// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/swinslow/peridot-jobrunner-testing/internal/testresult"
)

// GetContent makes an HTTP GET call to the indicated URL.
// It checks whether the expected HTTP status code is returned;
// a different code is treated as a failure.
// On success, it reads the response body into a got byte slice
// and handles closing the body. On failure, it fills in the
// failure code in the TestResult and returns an error.
func GetContent(res *testresult.TestResult, step string, url string, code int, ghUsername string) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		FailTest(res, step, err)
		return err
	}
	AddAuthHeader(res, step, req, ghUsername)
	resp, err := client.Do(req)

	return helperGetContent(res, resp, step, code)
}

// GetContentNoFollow makes an HTTP GET call to the indicated
// URL, and will NOT follow redirects. It otherwise acts
// identically to GetContent.
func GetContentNoFollow(res *testresult.TestResult, step string, url string, code int, ghUsername string) error {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest("GET", url, nil)
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

	return helperGetContent(res, resp, step, code)
}

// helperGetContent does the rest of the GetContent or
// GetContentNoFollow activities, after the decision is
// made on whether to follow any redirects.
func helperGetContent(res *testresult.TestResult, resp *http.Response, step string, code int) error {
	// parse content body
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
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
