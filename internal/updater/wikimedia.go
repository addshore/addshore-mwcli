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
	"strings"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

/*CanUpdate ...*/
func CanUpdateFromWikimedia(currentVersion string, gitSummary string, verboseOutput bool) (bool) {
	if(verboseOutput){
		selfupdate.EnableLog()
	}

	v, err := semver.Parse(strings.Trim(gitSummary,"v"))
	if err != nil {
		if(verboseOutput){
			log.Println("Could not parse git summary version, maybe you are not using a real release?")
		}
		return false
	}

	log.Println(v)

	// TODO check the latest.txt on rel.wm.o vs our local version

	return true;

}

/*UpdateFromWikimedia ...*/
func UpdateFromWikimedia(verboseOutput bool) (success bool, message string) {

	// TODO actually update from rel.wm.o

	// if(verboseOutput){
	// 	selfupdate.EnableLog()
	// }

	// cmdPath, err := os.Executable()
	// if err != nil {
	// 	return false, "Failed to grab local executable location"
	// }

	// err = selfupdate.UpdateTo(release.AssetURL, cmdPath)
	// if err != nil {
	// 	return false, "Binary update failed" + err.Error()
	// }

	// return true, "Successfully updated to version" + release.Version.String() + "\nRelease note:\n" + release.ReleaseNotes
}