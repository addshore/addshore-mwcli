/*Package mwdd is used to interact a mwdd v2 setup

Copyright © 2020 Addshore

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package mwdd

import (
	"fmt"
	"os"
	"os/user"
	"runtime"

	"gerrit.wikimedia.org/r/mediawiki/tools/cli/internal/env"
	"gerrit.wikimedia.org/r/mediawiki/tools/cli/internal/exec"
	"gerrit.wikimedia.org/r/mediawiki/tools/cli/internal/mwdd/files"
)

/*MWDD representation of a mwdd v2 setup*/
type MWDD string


/*DefaultForUser returns the default mwdd working directory for the user*/
func DefaultForUser() (MWDD) {
	return MWDD(mwddUserDirectory() + string(os.PathSeparator) + "default")
}

func mwddUserDirectory() string {
	currentUser, _ := user.Current()
	projectDirectory := currentUser.HomeDir + string(os.PathSeparator) + ".mwcli/mwdd"
	return projectDirectory
}

/*Directory the directory containing the development environment*/
func (m MWDD) Directory() string {
	return string(m)
}

/*DockerComposeProjectName the name of the docker-compose project*/
func (m MWDD) DockerComposeProjectName() string {
	return "mwcli-mwdd-default"
}

/*Env ...*/
func (m MWDD) Env() env.DotFile {
	return env.DotFileForDirectory(m.Directory())
}

/*EnsureReady ...*/
func (m MWDD) EnsureReady() {
	files.EnsureReady(m.Directory())
	m.EnsureEnvDefaults()
}

/*EnsureEnvDefaults ...*/
func (m MWDD) EnsureEnvDefaults() {
	neededVarDefaults := map[string]string{
		"MEDIAWIKI_VOLUMES_CODE": "~/dev/git/gerrit/mediawiki",
		"PORT": "8080",
	}
	env := m.Env()

	for key, value := range neededVarDefaults {
		if( env.Get(key) == "" ) {
			env.Set(key, value)
		}
	}

	// Always set the UID and GID (assume people shouldn't be setting this)
	if(runtime.GOOS == "windows") {
		// This user won't exist, but that fact doesn't really matter on Windows
		env.Set("UID", "2000")
		env.Set("GID", "2000")
	} else {
		if(env.Get("UID") == "") {
			env.Set("UID", fmt.Sprintf("%d",os.Getuid()))
		}
		if(env.Get("GID") == "") {
			env.Set("GID", fmt.Sprintf("%d",os.Getgid()))
		}
	}
}

/*EnsureHostsFile Make sure that a bunch of hosts that we will use are in the hosts file*/
func (m MWDD) EnsureHostsFile() {
	//TODO this will differ per serviceset though...
}

// DockerComposeCommand results in something like: `docker-compose <automatic project stuff> <command> <commandArguments>`
type DockerComposeCommand struct {
	Command      string
	CommandArguments    []string
	HandlerOptions exec.HandlerOptions
}

/*DockerCompose runs any docker-compose command for the mwdd project with the correct project settings and all files loaded*/
func (m MWDD) DockerCompose( command DockerComposeCommand ) {
	context := exec.ComposeCommandContext{
		ProjectDirectory: m.Directory(),
		ProjectName: m.DockerComposeProjectName(),
		Files: files.ListRawDcYamlFilesInContextOfProjectDirectory(m.Directory()),
	}

	exec.RunCommand(
		command.HandlerOptions,
		exec.ComposeCommand(
			context,
			command.Command,
			command.CommandArguments...
		),
	)
}

/*Exec runs `docker-compose exec -T <service> <commandAndArgs>`*/
func (m MWDD) Exec( service string, commandAndArgs []string, options exec.HandlerOptions ) {
	m.DockerCompose(
		DockerComposeCommand{
			Command: "exec",
			CommandArguments: append( []string{"-T", service }, commandAndArgs... ),
			HandlerOptions: options,
		},
	)
}

/*UpDetached runs `docker-compose up -d <services>`*/
func (m MWDD) UpDetached( services []string, options exec.HandlerOptions ) {
	m.DockerCompose(
		DockerComposeCommand{
			Command: "up",
			CommandArguments: append( []string{"-d" }, services... ),
			HandlerOptions: options,
		},
	)
}

/*DownWithVolumesAndOrphans runs `docker-compose down --volumes --remove-orphans`*/
func (m MWDD) DownWithVolumesAndOrphans( options exec.HandlerOptions ) {
	m.DockerCompose(
		DockerComposeCommand{
			Command: "down",
			CommandArguments: []string{"--volumes","--remove-orphans"},
			HandlerOptions: options,
		},
	)
}

/*Stop runs `docker-compose stop <services>`*/
func (m MWDD) Stop( services []string, options exec.HandlerOptions ) {
	m.DockerCompose(
		DockerComposeCommand{
			Command: "stop",
			CommandArguments: services,
			HandlerOptions: options,
		},
	)
}

/*Start runs `docker-compose start <services>`*/
func (m MWDD) Start( services []string, options exec.HandlerOptions ) {
	m.DockerCompose(
		DockerComposeCommand{
			Command: "start",
			CommandArguments: services,
			HandlerOptions: options,
		},
	)
}

/*Rm runs `docker-compose rm --stop --force -v <services>`*/
func (m MWDD) Rm( services []string, options exec.HandlerOptions ) {
	m.DockerCompose(
		DockerComposeCommand{
			Command: "rm",
			CommandArguments: append( []string{"--stop", "--force", "-v" }, services... ),
			HandlerOptions: options,
		},
	)
}

/*RmVolumes runs `docker volume rm <volume names with docker-compose project prefixed>`*/
func (m MWDD) RmVolumes( dcVolumes []string, options exec.HandlerOptions ) {
	dockerVolumes := []string{}
	for _, dcVolume := range dcVolumes {
		dockerVolumes = append( dockerVolumes, m.DockerComposeProjectName() + "_" + dcVolume )
	}
	exec.RunCommand(
		options,
		exec.Command("docker", append( []string{"volume", "rm" }, dockerVolumes... )... ),
	)
}

// TODO more from https://github.com/addshore/mediawiki-docker-dev/blob/4d380cf638bc60b5b6c22853a199639a3eb70b0b/control/src/Shell/DockerCompose.php#L53
// TODO execIt?
// TODO run?
// TODO runDetatched?
// TODO logsTail?
// TODO raw?