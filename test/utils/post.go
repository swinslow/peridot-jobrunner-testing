// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/swinslow/peridot-jobrunner-testing/internal/testresult"
)

// Post makes an HTTP POST call to the indicated URL with the
// specified body text.
// It checks whether the expected HTTP status code is returned;
// a different code is treated as a failure.
// On success, it reads the response body into a got byte slice
// and handles closing the body. On failure, it fills in the
// failure code in the TestResult and returns an error.
func Post(res *testresult.TestResult, step string, url string, bodystr string, code int, ghUsername string) error {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(bodystr))
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

// PostNoRes acts similarly to Post, but does not take a testresult
// or step value. It is primarily useful for fixture setup. It does
// not check the response body (but does ensure it is closed).
func PostNoRes(url string, bodystr string, code int, ghUsername string) error {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(bodystr))
	if err != nil {
		return err
	}
	AddAuthHeader(nil, "0", req, ghUsername)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// ignoring content body

	// check expected status code
	if resp.StatusCode != code {
		fmt.Printf("===> error in PostNoRes for %s\n", url)
		fmt.Printf("===> resp: %#v\n", resp)
		err = fmt.Errorf("expected HTTP status code %d, got %d", code, resp.StatusCode)
		return err
	}

	return nil
}
