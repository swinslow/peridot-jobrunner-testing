// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/swinslow/peridot-jobrunner-testing/internal/testresult"
)

// Put makes an HTTP PUT call to the indicated URL with the
// specified body text.
// It checks whether the expected HTTP status code is returned;
// a different code is treated as a failure.
// On success, it reads the response body into a got byte slice
// and handles closing the body. On failure, it fills in the
// failure code in the TestResult and returns an error.
func Put(res *testresult.TestResult, step string, url string, bodystr string, code int, ghUsername string) error {
	client := &http.Client{}
	req, err := http.NewRequest("PUT", url, strings.NewReader(bodystr))
	if err != nil {
		FailTest(res, step, err)
		return err
	}
	AddAuthHeader(res, step, req, ghUsername)
	resp, err := client.Do(req)
	if err != nil {
		FailTest(res, step, err)
		return err
	}
	defer resp.Body.Close()

	// parse content body
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		FailTest(res, step, err)
		return err
	}

	// record in testresult
	res.Got = b

	// check expected status code
	if resp.StatusCode != code {
		err = fmt.Errorf("expected HTTP status code %d, got %d", code, resp.StatusCode)
		FailTest(res, step, err)
		return err
	}

	return nil
}
