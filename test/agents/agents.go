// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package agents

import (
	"github.com/swinslow/peridot-jobrunner-testing/internal/testresult"
)

// GetTests returns all of the endpoints test suites.
func GetTests() []testresult.TestFunc {
	allTests := []testresult.TestFunc{}

	allTests = append(allTests, getNopTests()...)

	return allTests
}
