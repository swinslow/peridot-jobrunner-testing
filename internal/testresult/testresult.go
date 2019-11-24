// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package testresult

// TestResult contains data on the test, identifying it
// and whether it succeeded or failed.
type TestResult struct {
	// Suite is the overall type of test, e.g. "endpoint"
	Suite string

	// Element is the sub-type of test, e.g. "projects"
	Element string

	// ID is a unique identifier, within the element, for
	// a particular test, e.g. "GET-success"
	ID string

	// Success indicates whether the test succeeded.
	Success bool

	// FailStep indicates which step failed, if any.
	FailStep string

	// FailError provides the error of the failing step,
	// if any.
	FailError error

	// Wanted holds the latest JSON string that was desired.
	Wanted string

	// Got holds the latest JSON byte slice that was received.
	Got []byte
}

// TestFunc defines a function that takes a string with the
// API root URL for a test, and returns a TestResult.
type TestFunc func(string) *TestResult
