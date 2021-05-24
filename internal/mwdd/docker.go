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
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"time"

	"gerrit.wikimedia.org/r/mediawiki/tools/cli/internal/exec"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"golang.org/x/crypto/ssh/terminal"
)

// DockerExecCommand to be run with Docker, which directly uses the docker SDK
type DockerExecCommand struct {
	DockerComposeService      string
	Command      []string
	WorkingDir      string
	User      string
	HandlerOptions exec.HandlerOptions
}

// DockerRunCommand to be run with Docker, which directly uses the docker SDK
type DockerRunCommand struct {
	CustomContainerSuffix      string
	Image      string
	Command      []string
	WorkingDir      string
	User      string
	MountInPlace      []string
	HandlerOptions exec.HandlerOptions
}

/*UserAndGroupForDockerExecution gets a user and group id combination for the current user that can be used for execution*/
func UserAndGroupForDockerExecution() string {
	if(runtime.GOOS == "windows") {
		// TODO confirm that just using 2000 will always work on Windows?
		// This user won't exist, but that fact doesn't really matter on pure Windows
		return "2000:2000"
	}
	return fmt.Sprint(os.Getuid(), ":", os.Getgid())
}

/*DockerRun runs a docker run command using the dockerSDK*/
func (m MWDD) DockerRun( command DockerRunCommand ) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println("Unable to create docker client")
		panic(err)
	}

	ctx := context.Background()

	ContainerConfig := &container.Config{
		Image: command.Image,
		WorkingDir: command.WorkingDir,
		User: command.User,
		Entrypoint: command.Command,
		AttachStderr:true,
		AttachStdin: true,
		AttachStdout:true,
		Tty:		 true,
		OpenStdin:   true,
	}

	var emptyMountsSliceEntry []mount.Mount
	HostConfig := &container.HostConfig{
		Mounts: emptyMountsSliceEntry,
		PortBindings: nat.PortMap{},
		AutoRemove: true,
		// XXX: Must be kept in sync with docker-compose
		DNS: []string{"10.0.0.10"},
	}

	NetworkingConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			// XXX: Must be kept in sync with docker-compose
			"mwcli-mwdd-default_dps": &network.EndpointSettings{},
		},
	}

	for _, toMount := range command.MountInPlace {
		HostConfig.Mounts = append(
			HostConfig.Mounts,
			mount.Mount{
				Type:   mount.TypeBind,
				Source: toMount,
				Target: toMount,
			},
		)
	}

	// TODO if something breaks between creating and starting, then the container will hang around :( FIXME
	containerCreated, err := cli.ContainerCreate(
		ctx,
		ContainerConfig,
		HostConfig,
		NetworkingConfig,
		nil,
		"mwcli-mwdd-default-custom_" + command.CustomContainerSuffix,
		);
	if err != nil {
		fmt.Println("Error Creating container: " + containerCreated.ID)
		panic(err)
	}

	waiter, err := cli.ContainerAttach(ctx, containerCreated.ID, types.ContainerAttachOptions{
		Stderr:	   true,
		Stdout:	   true,
		Stdin:		true,
		Stream:	   true,
	})
	if err != nil {
		fmt.Println("Error Attaching to container: " + containerCreated.ID)
		panic(err)
	}

	err = cli.ContainerStart(ctx, containerCreated.ID, types.ContainerStartOptions{})
	if err != nil {
		fmt.Println("Error Starting container: " + containerCreated.ID)
		panic(err)
	}

	// When TTY is ON, just copy stdout https://phabricator.wikimedia.org/T282340
	// See: https://github.com/docker/cli/blob/70a00157f161b109be77cd4f30ce0662bfe8cc32/cli/command/container/hijack.go#L121-L130
	go io.Copy(os.Stdout, waiter.Reader)
	go io.Copy(waiter.Conn, os.Stdin)

	fd := int(os.Stdin.Fd())
	var oldState *terminal.State
	if terminal.IsTerminal(fd) {
		oldState, err = terminal.MakeRaw(fd)
		if err != nil {
			// print error
		}
		defer terminal.Restore(fd, oldState)
	}

	for {
		resp, err := cli.ContainerInspect(ctx, containerCreated.ID)
		time.Sleep(50 * time.Millisecond)
		if err != nil {
			break
		}

		if !resp.State.Running {
			break
		}
	}
}

/*DockerExec runs a docker exec command using the docker SDK for a service withing docker-compose*/
func (m MWDD) DockerExec( command DockerExecCommand ) {
	containerID := m.DockerComposeProjectName() + "_" + command.DockerComposeService + "_1"

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println("Unable to create docker client")
		panic(err)
	}

	config :=  types.ExecConfig{
		AttachStderr: true,
		AttachStdout: true,
		AttachStdin: true,
		Tty: true,
		WorkingDir: command.WorkingDir,
		User: command.User,
		Cmd: command.Command,
	}

	ctx := context.Background()
	response, err := cli.ContainerExecCreate(ctx, containerID, config)
	if err != nil {
		return
	}

	execID := response.ID
	if execID == "" {
		fmt.Println("exec ID empty")
		return
	}

	execStartCheck := types.ExecStartCheck{
		Tty: true,
	}

	waiter, err := cli.ContainerExecAttach(ctx, execID, execStartCheck)
	if err != nil {
		fmt.Println(err)
		return
	}

	// When TTY is ON, just copy stdout https://phabricator.wikimedia.org/T282340
	// See: https://github.com/docker/cli/blob/70a00157f161b109be77cd4f30ce0662bfe8cc32/cli/command/container/hijack.go#L121-L130
	go io.Copy(os.Stdout, waiter.Reader)
	go io.Copy(waiter.Conn, os.Stdin)

	fd := int(os.Stdin.Fd())
	var oldState *terminal.State
	if terminal.IsTerminal(fd) {
		oldState, err = terminal.MakeRaw(fd)
		if err != nil {
			// print error
		}
		defer terminal.Restore(fd, oldState)
	}

	for {
		resp, err := cli.ContainerExecInspect(ctx, execID)
		time.Sleep(50 * time.Millisecond)
		if err != nil {
			break
		}

		if !resp.Running {
			break
		}
	}
}
