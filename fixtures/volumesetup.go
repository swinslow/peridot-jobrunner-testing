// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package fixtures

import (
	"io/ioutil"
	"os"
	"path"
)

// ResetVolume clears the /code and /spdx directories.
func ResetVolume() error {
	// delete contents of /code directory but not /code itself
	contents, err := ioutil.ReadDir("/code")
	if err != nil {
		return err
	}
	for _, c := range contents {
		os.RemoveAll(path.Join([]string{"code", c.Name()}...))
	}

	// delete contents of /spdx directory but not /spdx itself
	contents, err = ioutil.ReadDir("/spdx")
	if err != nil {
		return err
	}
	for _, c := range contents {
		os.RemoveAll(path.Join([]string{"spdx", c.Name()}...))
	}

	return nil
}
