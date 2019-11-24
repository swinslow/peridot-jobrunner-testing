// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package utils

import (
	"fmt"
	"net/http"

	"github.com/swinslow/peridot-jobrunner-testing/internal/testresult"
	"github.com/yudai/gojsondiff"
)

// Pass fills in the success fields.
func Pass(res *testresult.TestResult) {
	res.Success = true
}

// FailTest fills in the failure fields for a test that failed
// for some reason other than because the JSON strings did not
// match.
func FailTest(res *testresult.TestResult, step string, msg error) {
	res.Success = false
	res.FailStep = step
	res.FailError = msg
}

// FailMatch fills in the failure fields for a test that failed
// because the desired JSON string did not match the JSON string
// that was received.
func FailMatch(res *testresult.TestResult, step string) {
	res.Success = false
	res.FailStep = step
}

// IsMatch compares a wanted string and a got byte slice containing
// JSON data, and returns a bool indicating whether they contained
// equivalent content. It will also return "false" if there is e.g.
// an error with the JSON unmarshalling, etc.
func IsMatch(res *testresult.TestResult) bool {
	differ := gojsondiff.New()
	d, err := differ.Compare([]byte(res.Wanted), res.Got)
	// fmt.Printf("*** WANTED:  %#v\n", wanted)
	// fmt.Printf("*** GOT:     %#v\n", got)
	// fmt.Printf("*** ERR:     %#v\n", err)
	// fmt.Printf("*** DIFF:    %#v\n", d)
	if err != nil {
		return false
	}

	return !d.Modified()
}

// IsEmpty checks for an empty wanted string and a zero-length got
// byte slice.
func IsEmpty(res *testresult.TestResult) bool {
	return res.Wanted == "" && len(res.Got) == 0
}

// AddAuthHeader adds the appropriate auth token header to the
// request object, before it is sent. Including "none" as the
// username means that no token will be sent. The token values
// included here are JWT values for the signing key "keyForTesting"
// and it is intentional that they are used here -- but of course
// they should not be used in production in any way!
func AddAuthHeader(res *testresult.TestResult, step string, req *http.Request, ghUsername string) {
	switch ghUsername {
	case "none":
		return
	case "admin":
		req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJnaXRodWIiOiJhZG1pbiJ9.3KnAxp1Tn7O8vHQXBReUy81qG7qfRPsxRXW8Wr68xfc")
	case "operator":
		req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJnaXRodWIiOiJvcGVyYXRvciJ9.v8xJrGfBKDj9OYF2G58NeV1sGfKNahr-OHzqCXetwUU")
	case "commenter":
		req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJnaXRodWIiOiJjb21tZW50ZXIifQ.PQDdHhSmjDs9sceGi54cT71ef2IVxiO_Yw0-_YDJ-i8")
	case "viewer":
		req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJnaXRodWIiOiJ2aWV3ZXIifQ.YQUkHNTbsfA3ldtfxhTkoFI8eHVhfbFLF5vkmOrFJZg")
	case "disabled":
		req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJnaXRodWIiOiJkaXNhYmxlZCJ9.mqdsZIPPEb1RmmdI1zO0elHFieHbzmleYdg06qRfVbQ")
	default:
		if res != nil {
			FailTest(res, step, fmt.Errorf("invalid username %s", ghUsername))
		}
	}
}
