// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package agents

import (
	"github.com/swinslow/peridot-jobrunner-testing/internal/testresult"
	"github.com/swinslow/peridot-jobrunner-testing/test/utils"
)

func getNopTests() []testresult.TestFunc {
	return []testresult.TestFunc{
		jobsSubGetOperator,
		jobsSubPostOperator,
		jobsGetOneViewer,
		jobsPutOneOperator,
		jobsPutOneViewer,
		jobsDeleteOneAdmin,
		jobsDeleteOneOperator,
	}
}

// ===== GET /repopulls/id/jobs

func jobsSubGetOperator(root string) *testresult.TestResult {
	res := &testresult.TestResult{
		Suite:   "endpoints",
		Element: "repopulls/{id}/jobs",
		ID:      "GET (viewer)",
	}

	url := root + "/repopulls/4/jobs"

	res.Wanted = `{"jobs":[
		{"id":2, "repopull_id":4, "agent_id":1, "started_at":"0001-01-01T00:00:00Z", "finished_at":"0001-01-01T00:00:00Z", "status":"startup", "health":"ok", "is_ready":true, "config":{}},
		{"id":3, "repopull_id":4, "agent_id":2, "priorjob_ids": [2], "started_at":"0001-01-01T00:00:00Z", "finished_at":"0001-01-01T00:00:00Z", "status":"startup", "health":"ok", "is_ready":true, "config":{"codereader": {"primary": {"path": "/somewhere"}}}},
		{"id":4, "repopull_id":4, "agent_id":4, "priorjob_ids": [2,3], "started_at":"0001-01-01T00:00:00Z", "finished_at":"0001-01-01T00:00:00Z", "status":"startup", "health":"ok", "is_ready":false, "config":{"kv": {"hello":"world"}, "codereader": {"godeps": {"priorjob_id": 3}}, "spdxreader": {"primary": {"path": "/path/wherever"}, "godeps": {"priorjob_id": 3}}}}
	]}`
	err := utils.GetContent(res, "1", url, 200, "viewer")
	if err != nil {
		return res
	}

	if !utils.IsMatch(res) {
		utils.FailMatch(res, "2")
		return res
	}

	utils.Pass(res)
	return res
}

// ===== POST /repopulls/id/jobs

func jobsSubPostOperator(root string) *testresult.TestResult {
	res := &testresult.TestResult{
		Suite:   "endpoints",
		Element: "repopulls/{id}/jobs",
		ID:      "POST (operator)",
	}

	url := root + "/repopulls/3/jobs"

	// first, send POST to add a new job
	body := `{"agent_id":1, "is_ready":false, "priorjob_ids":[],
		"config":{"kv": {"hi": "there", "hello": "world"}}
	}`
	res.Wanted = `{"id": 5}`
	err := utils.Post(res, "1", url, body, 201, "operator")
	if err != nil {
		return res
	}

	if !utils.IsMatch(res) {
		utils.FailMatch(res, "2")
		return res
	}

	// now, confirm that a new job was actually added
	// this should be the only one for repopull 3 so we can reuse the same url
	// priorjob_ids and some config vals should be absent
	res.Wanted = `{"jobs":[
		{"id":5, "repopull_id":3, "agent_id":1, "started_at":"0001-01-01T00:00:00Z", "finished_at":"0001-01-01T00:00:00Z", "status":"startup", "health":"ok", "is_ready":false, "config":{"kv": {"hi": "there", "hello": "world"}}}
	]}`
	err = utils.GetContent(res, "3", url, 200, "operator")
	if err != nil {
		return res
	}

	if !utils.IsMatch(res) {
		utils.FailMatch(res, "4")
		return res
	}

	utils.Pass(res)
	return res
}

// ===== GET /jobs/id

func jobsGetOneViewer(root string) *testresult.TestResult {
	res := &testresult.TestResult{
		Suite:   "endpoints",
		Element: "jobs/{id}",
		ID:      "GET (viewer)",
	}

	url := root + "/jobs/4"

	res.Wanted = `{"job":{"id":4, "repopull_id":4, "agent_id":4, "priorjob_ids": [2,3], "started_at":"0001-01-01T00:00:00Z", "finished_at":"0001-01-01T00:00:00Z", "status":"startup", "health":"ok", "is_ready":false, "config":{"kv": {"hello":"world"}, "codereader": {"godeps": {"priorjob_id": 3}}, "spdxreader": {"primary": {"path": "/path/wherever"}, "godeps": {"priorjob_id": 3}}}}}`
	err := utils.GetContent(res, "1", url, 200, "viewer")
	if err != nil {
		return res
	}

	if !utils.IsMatch(res) {
		utils.FailMatch(res, "2")
		return res
	}

	utils.Pass(res)
	return res
}

// ===== PUT /jobs/id

func jobsPutOneOperator(root string) *testresult.TestResult {
	res := &testresult.TestResult{
		Suite:   "endpoints",
		Element: "jobs/{id}",
		ID:      "PUT (operator)",
	}

	url := root + "/jobs/4"

	// first, send PUT to update an existing job
	// only is_ready can currently be updated
	body := `{"is_ready": true}`
	res.Wanted = ``
	err := utils.Put(res, "1", url, body, 204, "operator")
	if err != nil {
		return res
	}

	if !utils.IsEmpty(res) {
		utils.FailMatch(res, "2")
		return res
	}

	// now, confirm that the job was actually updated
	// is_ready should now be true
	res.Wanted = `{"job":{"id":4, "repopull_id":4, "agent_id":4, "priorjob_ids": [2,3], "started_at":"0001-01-01T00:00:00Z", "finished_at":"0001-01-01T00:00:00Z", "status":"startup", "health":"ok", "is_ready":true, "config":{"kv": {"hello":"world"}, "codereader": {"godeps": {"priorjob_id": 3}}, "spdxreader": {"primary": {"path": "/path/wherever"}, "godeps": {"priorjob_id": 3}}}}}`
	err = utils.GetContent(res, "3", url, 200, "operator")
	if err != nil {
		return res
	}

	if !utils.IsMatch(res) {
		utils.FailMatch(res, "4")
		return res
	}

	utils.Pass(res)
	return res
}

func jobsPutOneViewer(root string) *testresult.TestResult {
	res := &testresult.TestResult{
		Suite:   "endpoints",
		Element: "jobs/{id}",
		ID:      "PUT (viewer)",
	}

	url := root + "/jobs/4"

	body := `{"is_ready": true}`
	res.Wanted = `{"error": "Access denied"}`
	err := utils.Put(res, "1", url, body, 403, "viewer")
	if err != nil {
		return res
	}

	if !utils.IsMatch(res) {
		utils.FailMatch(res, "2")
		return res
	}

	// now, confirm that the job was NOT actually updated
	// is_ready should still be false
	res.Wanted = `{"job":{"id":4, "repopull_id":4, "agent_id":4, "priorjob_ids": [2,3], "started_at":"0001-01-01T00:00:00Z", "finished_at":"0001-01-01T00:00:00Z", "status":"startup", "health":"ok", "is_ready":false, "config":{"kv": {"hello":"world"}, "codereader": {"godeps": {"priorjob_id": 3}}, "spdxreader": {"primary": {"path": "/path/wherever"}, "godeps": {"priorjob_id": 3}}}}}`
	err = utils.GetContent(res, "3", url, 200, "operator")
	if err != nil {
		return res
	}

	if !utils.IsMatch(res) {
		utils.FailMatch(res, "4")
		return res
	}

	utils.Pass(res)
	return res
}

// ===== DELETE /jobs/id

func jobsDeleteOneAdmin(root string) *testresult.TestResult {
	res := &testresult.TestResult{
		Suite:   "endpoints",
		Element: "jobs/{id}",
		ID:      "DELETE (admin)",
	}

	url := root + "/jobs/3"

	// send a delete request
	res.Wanted = ``
	err := utils.Delete(res, "1", url, ``, 204, "admin")
	if err != nil {
		return res
	}

	if !utils.IsEmpty(res) {
		utils.FailMatch(res, "2")
		return res
	}

	// now, confirm that the job is gone
	// NOTE that job ID 3 is also removed from priorjob_ids and config for job 4.
	// FIXME the deleted job should not cascade in this way.
	allURL := root + "/repopulls/4/jobs"
	res.Wanted = `{"jobs":[
		{"id":2, "repopull_id":4, "agent_id":1, "started_at":"0001-01-01T00:00:00Z", "finished_at":"0001-01-01T00:00:00Z", "status":"startup", "health":"ok", "is_ready":true, "config":{}},
		{"id":4, "repopull_id":4, "agent_id":4, "priorjob_ids": [2], "started_at":"0001-01-01T00:00:00Z", "finished_at":"0001-01-01T00:00:00Z", "status":"startup", "health":"ok", "is_ready":false, "config":{"kv": {"hello":"world"}, "spdxreader": {"primary": {"path": "/path/wherever"}}}}
	]}`
	err = utils.GetContent(res, "3", allURL, 200, "viewer")
	if err != nil {
		return res
	}

	if !utils.IsMatch(res) {
		utils.FailMatch(res, "4")
		return res
	}

	utils.Pass(res)
	return res
}

func jobsDeleteOneOperator(root string) *testresult.TestResult {
	res := &testresult.TestResult{
		Suite:   "endpoints",
		Element: "jobs/{id}",
		ID:      "DELETE (operator)",
	}

	url := root + "/jobs/3"

	// try and fail to delete the job
	res.Wanted = `{"error": "Access denied"}`
	err := utils.Delete(res, "1", url, ``, 403, "operator")
	if err != nil {
		return res
	}

	if !utils.IsMatch(res) {
		utils.FailMatch(res, "2")
		return res
	}

	// now, confirm that the job has NOT been deleted
	allURL := root + "/repopulls/4/jobs"
	res.Wanted = `{"jobs":[
		{"id":2, "repopull_id":4, "agent_id":1, "started_at":"0001-01-01T00:00:00Z", "finished_at":"0001-01-01T00:00:00Z", "status":"startup", "health":"ok", "is_ready":true, "config":{}},
		{"id":3, "repopull_id":4, "agent_id":2, "priorjob_ids": [2], "started_at":"0001-01-01T00:00:00Z", "finished_at":"0001-01-01T00:00:00Z", "status":"startup", "health":"ok", "is_ready":true, "config":{"codereader": {"primary": {"path": "/somewhere"}}}},
		{"id":4, "repopull_id":4, "agent_id":4, "priorjob_ids": [2,3], "started_at":"0001-01-01T00:00:00Z", "finished_at":"0001-01-01T00:00:00Z", "status":"startup", "health":"ok", "is_ready":false, "config":{"kv": {"hello":"world"}, "codereader": {"godeps": {"priorjob_id": 3}}, "spdxreader": {"primary": {"path": "/path/wherever"}, "godeps": {"priorjob_id": 3}}}}
	]}`
	err = utils.GetContent(res, "3", allURL, 200, "viewer")
	if err != nil {
		return res
	}

	if !utils.IsMatch(res) {
		utils.FailMatch(res, "4")
		return res
	}

	utils.Pass(res)
	return res
}
