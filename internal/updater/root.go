/*Package updater is used to update the cli

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
package updater

import (
	"log"
	"os"
	"strings"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

/*CanUpdate ...*/
func CanUpdate(currentVersion string, gitSummary string, verboseOutput bool) (bool, *selfupdate.Release) {
	if(verboseOutput){
		selfupdate.EnableLog()
	}

	v := semver.MustParse(strings.Trim(gitSummary,"v"))

	// TODO when builds are on wm.o then allow for a "dev" or "stable" update option and checks

	rel, ok, err := selfupdate.DetectLatest("addshore/mwcli")
	if err != nil {
		if(verboseOutput){
			log.Println("Some unknown error occurred")
		}
		return false, rel
	}
	if !ok {
		if(verboseOutput){
			log.Println("No release detected. Current version is considered up-to-date")
		}
		return false, rel
	}
	if v.Equals(rel.Version) {
		if(verboseOutput){
			log.Println("Current version", v, "is the latest. Update is not needed")
		}
		return false, rel
	}
	if(verboseOutput){
		log.Println("Update available", rel.Version)
	}
	return true, rel
}

/*UpdateTo ...*/
func UpdateTo(release selfupdate.Release, verboseOutput bool) (success bool, message string) {
	if(verboseOutput){
		selfupdate.EnableLog()
	}

	cmdPath, err := os.Executable()
	if err != nil {
		return false, "Failed to grab local executable location"
	}

	err = selfupdate.UpdateTo(release.AssetURL, cmdPath)
	if err != nil {
		if(verboseOutput){
			log.Println("Binary update failed:", err)
		}
		return false, "Binary update failed"
	}

	return true, "Successfully updated to version" + release.Version.String() + "\nRelease note:\n" + release.ReleaseNotes
}

/*ShouldAllowUpdates ...*/
func ShouldAllowUpdates(currentVersion string, gitSummary string, verboseOutput bool) bool {
	if !strings.HasPrefix(gitSummary, currentVersion) || strings.HasSuffix(gitSummary,"dirty") {
		if(verboseOutput){
			log.Println("Can only update tag built releases")
		}
		return false
	}
	return true
}
