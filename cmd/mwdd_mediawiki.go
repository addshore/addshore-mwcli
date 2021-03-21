/*Package cmd is used for command line.

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
package cmd

import (
	"fmt"

	"gerrit.wikimedia.org/r/mediawiki/tools/cli/internal/exec"
	"gerrit.wikimedia.org/r/mediawiki/tools/cli/internal/mwdd"
	"github.com/spf13/cobra"
)

var mwddMediawikiCmd = &cobra.Command{
	Use:   "mediawiki",
	Short: "MediaWiki service",
	RunE:  nil,
}

var mwddMediawikiInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs a new MediaWiki site using install.php",
	Run: func(cmd *cobra.Command, args []string) {
		dbname := "default"
		mwdd.DefaultForUser().Exec("mediawiki",[]string{
			"php",
			"/var/www/html/w/maintenance/install.php",
			"--dbuser", "root",
			"--dbpass", "toor",
			"--dbname", dbname,
			"--dbserver", "db-master",
			"--lang", "en",
			"--pass", "mwddpassword",
			"docker-" + dbname,
			"admin",
			}, exec.HandlerOptions{})
	},
}

var mwddMediawikiComposerCmd = &cobra.Command{
	Use:   "composer",
	Short: "Runs composer in a container in the context of MediaWiki",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO run a composer container with the mediawiki main directory checked out (and other volumes?)
		// as the correct user and on the correct network?
		// READ https://github.com/addshore/mediawiki-docker-dev/blob/v1/control/src/Command/MediaWiki/Composer.php
		//READ https://github.com/addshore/mediawiki-docker-dev/blob/v1/docker-compose/mw-composer.yml
		mwdd.DefaultForUser().EnsureReady()
		options := exec.HandlerOptions{
			Verbosity:   Verbosity,
		}
		// TODO only run up if not already up?
		mwdd.DefaultForUser().UpDetached( []string{"composer"}, options )
		mwdd.DefaultForUser().Run(
			"composer",
			// TODO pass in other arguments
			[]string{"info"},
			options,
			)
	},
}

var mwddMediawikiCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create the Mediawiki containers",
	Run: func(cmd *cobra.Command, args []string) {
		mwdd.DefaultForUser().EnsureReady()
		options := exec.HandlerOptions{
			Verbosity:   Verbosity,
		}
		// TODO check mediawiki is here..
		mwdd.DefaultForUser().UpDetached( []string{"mediawiki"}, options )
	},
}

var mwddMediawikiDestroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy the Mediawiki containers",
	Run: func(cmd *cobra.Command, args []string) {
		mwdd.DefaultForUser().EnsureReady()
		options := exec.HandlerOptions{
				Verbosity:   Verbosity,
		}
		mwdd.DefaultForUser().DownWithVolumesAndOrphans( options )
},
}

var mwddMediawikiSuspendCmd = &cobra.Command{
	Use:   "suspend",
	Short: "Suspend the Mediawiki containers",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Not yet implemented!");
	},
}

var mwddMediawikiResumeCmd = &cobra.Command{
	Use:   "resume",
	Short: "Resume the Mediawiki containers",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Not yet implemented!");
	},
}

func init() {
	mwddCmd.AddCommand(mwddMediawikiCmd)
	mwddMediawikiCmd.AddCommand(mwddMediawikiCreateCmd)
	mwddMediawikiCmd.AddCommand(mwddMediawikiDestroyCmd)
	mwddMediawikiCmd.AddCommand(mwddMediawikiSuspendCmd)
	mwddMediawikiCmd.AddCommand(mwddMediawikiResumeCmd)
	mwddMediawikiCmd.AddCommand(mwddMediawikiInstallCmd)
	mwddMediawikiCmd.AddCommand(mwddMediawikiComposerCmd)
}
