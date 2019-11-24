// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package main

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"text/tabwriter"

	"github.com/swinslow/peridot-jobrunner-testing/fixtures"
	"github.com/swinslow/peridot-jobrunner-testing/internal/testresult"
	"github.com/swinslow/peridot-jobrunner-testing/test/agents"
)

func main() {
	apiRoot := "http://api:3005"
	anyFailed := false

	allRs := []*testresult.TestResult{}
	var rs *testresult.TestResult

	// get all test suites
	allTests := agents.GetTests()

	// and run them, resetting DB and volume each time
	fmt.Printf("Testing (%d total): \n", len(allTests))
	for _, t := range allTests {
		fmt.Printf("  %s\n", runtime.FuncForPC(reflect.ValueOf(t).Pointer()).Name())
		err := fixtures.ResetVolume()
		if err != nil {
			fmt.Printf("Error resetting volume before test: %v\n", err)
			os.Exit(1)
		}
		err = fixtures.ResetDB(apiRoot)
		if err != nil {
			fmt.Printf("Error resetting DB before test: %v\n", err)
			os.Exit(1)
		}
		err = fixtures.SetupFixture(apiRoot)
		if err != nil {
			fmt.Printf("Error setting fixtures before test: %v\n", err)
			os.Exit(1)
		}

		rs = t(apiRoot)
		allRs = append(allRs, rs)
	}

	fmt.Printf("\n\n")

	// set up tabwriter for outputting test result table
	w := tabwriter.NewWriter(os.Stdout, 8, 4, 1, ' ', 0)

	// output results
	for _, r := range allRs {
		var result string
		if r.Success {
			result = "ok"
		} else {
			result = "FAIL"
			anyFailed = true
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", r.Suite, r.Element, r.ID, result)
	}
	w.Flush()

	if anyFailed {
		// print details of failing tests
		fmt.Printf("\n\n==========\n\n")
		for _, r := range allRs {
			if !r.Success {
				fmt.Printf("%s:%s:%s\n", r.Suite, r.Element, r.ID)
				fmt.Printf("    Status: FAIL\n")
				fmt.Printf("    Step:   %s\n", r.FailStep)
				fmt.Printf("    Errors: %v\n", r.FailError)
				fmt.Printf("    Wanted: %s\n", r.Wanted)
				fmt.Printf("    Got:    %s\n", r.Got)
				fmt.Printf("\n==========\n\n")
			}
		}

		// return failure status code
		os.Exit(1)
	}
}
