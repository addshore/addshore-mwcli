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
	"gerrit.wikimedia.org/r/mediawiki/tools/cli/internal/exec"
	"gerrit.wikimedia.org/r/mediawiki/tools/cli/internal/mwdd"
	"fmt"
	"github.com/spf13/cobra"
)

var mwddCmd = &cobra.Command{
	Use:   "mwdd",
	Short: "The MediaWiki-Docker-Dev like development environment",
	RunE:  nil,
}

var mwddWhereCmd = &cobra.Command{
	Use:   "where",
	Short: "States the working directory for the environment",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(mwdd.DefaultForUser().Directory())
	},
}

var mwddCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create the Default containers",
	Run: func(cmd *cobra.Command, args []string) {
		mwdd.DefaultForUser().EnsureReady()
		options := exec.HandlerOptions{
			Verbosity:   Verbosity,
		}
		// TODO mediawiki should come from some default definition set?
		mwdd.DefaultForUser().UpDetached( []string{"mediawiki"}, options )
		// TODO add functionality for writing to the hosts file...
		//mwdd.DefaultForUser().EnsureHostsFile()
	},
}

var mwddDestroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy the Default containers",
	Run: func(cmd *cobra.Command, args []string) {
		mwdd.DefaultForUser().EnsureReady()
		options := exec.HandlerOptions{
			Verbosity:   Verbosity,
		}
		mwdd.DefaultForUser().DownWithVolumesAndOrphans( options )
		// TODO add functionality for writing to the hosts file...
		//mwdd.DefaultForUser().EnsureHostsFile()
	},
}

var mwddSuspendCmd = &cobra.Command{
	Use:   "suspend",
	Short: "Suspend the Default containers",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Not yet implemented!");
	},
}

var mwddResumeCmd = &cobra.Command{
	Use:   "resume",
	Short: "Resume the Default containers",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Not yet implemented!");
	},
}

func init() {
	mwddCmd.PersistentFlags().IntVarP(&Verbosity, "verbosity", "v", 1, "verbosity level (1-2)")

	mwddCmd.AddCommand(mwddWhereCmd)
	mwddCmd.AddCommand(mwddCreateCmd)
	mwddCmd.AddCommand(mwddDestroyCmd)
	mwddCmd.AddCommand(mwddSuspendCmd)
	mwddCmd.AddCommand(mwddResumeCmd)

	rootCmd.AddCommand(mwddCmd)
}
