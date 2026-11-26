package app

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"yingce/backend/internal/protocol"
)

// Use the shipped package and actual HTTP request builder. The fixture CDN URL
// is only parsed, never fetched; no suppliers or external media are contacted.
func TestGrsaiWanVideoDetachedHTTP(t *testing.T) {
	allowLoopbackProviderTest(t)
	const id = "fixture/id+1&next"
	const clip = "https://cdn.example.test/clip.mp4"
	adapter := loadOfficialProtocolAdapterForTest(t, "grsai-wan-video.yingce-plugin")

	for _, scenario := range []string{"success", "violation", "mismatched-id", "http400"} {
		t.Run(scenario, func(t *testing.T) {
			var creates, polls atomic.Int32
			upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if r.Header.Get("Authorization") != "Bearer fixture-only" {
					t.Error("missing host-injected test credential")
				}
				switch r.Method + " " + r.URL.Path {
				case "POST /v1/api/generate":
					creates.Add(1)
					var body map[string]any
					if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body["replyType"] != "async" || body["firstFrame"] != "https://cdn.example.test/frame.png" || body["images"] != nil {
						t.Error("single image in image-to-video mode was not sent as firstFrame in an async request")
					}
					if scenario == "http400" {
						w.WriteHeader(http.StatusBadRequest)
						_ = json.NewEncoder(w).Encode(map[string]string{"id": id, "status": "failed", "error": "invalid reference"})
						return
					}
					_ = json.NewEncoder(w).Encode(map[string]string{"id": id, "status": "running"})
				case "GET /v1/api/result":
					if r.URL.Query().Get("id") != id || len(r.URL.Query()) != 1 {
						t.Errorf("unexpected poll query: %v", r.URL.Query())
					}
					if polls.Add(1) == 1 {
						_ = json.NewEncoder(w).Encode(map[string]any{"id": id, "status": "running", "progress": 25})
					} else if scenario == "violation" {
						_ = json.NewEncoder(w).Encode(map[string]string{"id": id, "status": "violation", "error": "content rejected"})
					} else if scenario == "mismatched-id" {
						_ = json.NewEncoder(w).Encode(map[string]string{"id": "someone-elses-job", "status": "failed"})
					} else {
						_ = json.NewEncoder(w).Encode(map[string]any{"id": id, "status": "succeeded", "results": []map[string]string{{"url": clip}}})
					}
				default:
					t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
					http.NotFound(w, r)
				}
			}))
			defer upstream.Close()

			input := canvasGenerationInput{
				Mode: "video", Prompt: "fixture clip", Metadata: map[string]any{"videoOperation": "image_to_video"},
				ReferenceImages: []providerMedia{{ID: "frame", Type: "image/png", URL: "https://cdn.example.test/frame.png"}},
				Config:          providerConfig{InterfaceType: "grsai-wan-video", Model: "wan-3.0", BaseURL: upstream.URL + "/v1", APIKey: "fixture-only", Size: "16:9", VQuality: "720p", VideoSeconds: "5"},
			}
			run := func(ctx context.Context) error {
				_, err := runProtocolAdapterTaskWithPolicy(ctx, input, adapter, fastVideoPollPolicy())
				return err
			}
			err := run(t.Context())
			if scenario == "http400" {
				var upstreamErr providerHTTPError
				if !errors.As(err, &upstreamErr) || upstreamErr.StatusCode != http.StatusBadRequest || !strings.Contains(err.Error(), "invalid reference") || creates.Load() != 1 || polls.Load() != 0 {
					t.Fatalf("create rejection: %v, creates=%d polls=%d", err, creates.Load(), polls.Load())
				}
				return
			}
			var pending providerStatePendingError
			if !errors.As(err, &pending) || !pending.Detached || pending.TaskID != id {
				t.Fatalf("create receipt: %v", err)
			}
			resumed := context.WithValue(t.Context(), providerAnalyticsKey{}, providerAnalyticsContext{ProviderRequestID: pending.TaskID})
			if err = run(resumed); !errors.As(err, &pending) || pending.TaskID != id || pending.Stage != "25%" {
				t.Fatalf("running receipt: %v (%+v)", err, pending)
			}
			if scenario == "violation" {
				var terminal providerTaskTerminalError
				if err = run(resumed); !errors.As(err, &terminal) || terminal.Status != protocol.StatusFailed || terminal.Message != "content rejected" {
					t.Fatalf("violation did not end the task: %v", err)
				}
			} else if scenario == "mismatched-id" {
				var terminal providerTaskTerminalError
				if err = run(resumed); err == nil || errors.As(err, &terminal) || !strings.Contains(err.Error(), "不匹配的任务 ID") {
					t.Fatalf("another job must not determine this job's terminal status: %v", err)
				}
			} else {
				state, err := pollProtocolTask(resumed, input.Config, adapter, protocolRequestFromInput(input), pending.TaskID)
				if err != nil || state.Status != protocol.StatusSucceeded || state.Result == nil || len(state.Result.Videos) != 1 || state.Result.Videos[0].URL != clip {
					t.Fatalf("success receipt: %+v, error=%v", state, err)
				}
			}
			if creates.Load() != 1 || polls.Load() != 2 {
				t.Fatalf("recovery must only poll: creates=%d polls=%d", creates.Load(), polls.Load())
			}
		})
	}
}
