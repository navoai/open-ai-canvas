package protocol

import (
	"net/url"
	"reflect"
	"strings"
	"testing"
)

const grsaiWanProvider = "grsai-wan-video"

func grsaiWanRequest() GenerationRequest {
	return GenerationRequest{
		Capability: CapabilityVideo, Model: "wan-3.0", Prompt: "船驶过湖面",
		AspectRatio: "16:9", Resolution: "720p", Duration: 5,
	}
}

func TestGrsaiWanVideoCreateAndPollRequests(t *testing.T) {
	adapter := officialPackageAdapter(t, grsaiWanProvider+".yingce-plugin", grsaiWanProvider)
	for _, base := range []string{"https://grsaiapi.com", "https://grsai.dakka.com.cn"} {
		t.Run(base, func(t *testing.T) {
			req := grsaiWanRequest()
			create, err := adapter.BuildCreate(t.Context(), RequestContext{BaseURL: base, Request: req})
			if err != nil {
				t.Fatal(err)
			}
			body := manifestTestBody(t, create)
			want := map[string]any{"model": "wan-3.0", "prompt": req.Prompt, "aspectRatio": "16:9", "resolution": "720p", "duration": float64(5), "replyType": "async"}
			if create.Method != "POST" || create.Path != "/v1/api/generate" || create.ContentType != "application/json" || !reflect.DeepEqual(body, want) {
				t.Fatalf("create: %s %s (%s), body %#v", create.Method, create.Path, create.ContentType, body)
			}
			// The adapter must supply one raw query value. Encoding belongs to the host;
			// pre-escaping the ID or interpolating it in the path breaks opaque task IDs.
			const id = "opaque/id+A&B=?%"
			poll, err := adapter.BuildPoll(t.Context(), PollContext{BaseURL: base, TaskID: id})
			if err != nil {
				t.Fatal(err)
			}
			if poll.Method != "GET" || poll.Path != "/v1/api/result" || poll.Body != nil ||
				!reflect.DeepEqual(poll.Query, map[string][]string{"id": {id}}) ||
				url.Values(poll.Query).Encode() != "id=opaque%2Fid%2BA%26B%3D%3F%25" {
				t.Fatalf("poll = %#v", poll)
			}
			for _, spec := range []RequestSpec{create, poll} {
				if spec.Auth.Type != "bearer" || spec.Auth.Field != "apiKey" || spec.Headers["Authorization"] != "" {
					t.Fatalf("expected host-injected bearer auth, got %#v with headers %#v", spec.Auth, spec.Headers)
				}
			}
		})
	}
	if adapter.Metadata().IdempotentCreate || adapter.Metadata().Cancel != "" {
		t.Fatal("provider claims undocumented create idempotency or upstream cancellation")
	}
}

func TestGrsaiWanVideoOutputContract(t *testing.T) {
	adapter := officialPackageAdapter(t, grsaiWanProvider+".yingce-plugin", grsaiWanProvider)
	for i, ratio := range []string{"16:9", "9:16", "21:9", "4:3", "1:1", "3:4"} {
		req := grsaiWanRequest()
		req.AspectRatio, req.Resolution, req.Duration = ratio, []string{"480p", "720p", "1080p"}[i%3], []int{2, 30}[i%2]
		spec, err := adapter.BuildCreate(t.Context(), RequestContext{Request: req})
		if err != nil {
			t.Fatalf("%s/%s/%ds: %v", ratio, req.Resolution, req.Duration, err)
		}
		body := manifestTestBody(t, spec)
		if body["aspectRatio"] != ratio || body["resolution"] != req.Resolution || body["duration"] != float64(req.Duration) {
			t.Fatalf("wrong JSON names, types, or values: %#v", body)
		}
	}
	for name, mutate := range map[string]func(*GenerationRequest){
		"other model":                          func(r *GenerationRequest) { r.Model = "wan3.0-video" },
		"missing prompt":                       func(r *GenerationRequest) { r.Prompt = "" },
		"missing ratio":                        func(r *GenerationRequest) { r.AspectRatio = "" },
		"unknown ratio":                        func(r *GenerationRequest) { r.AspectRatio = "auto" },
		"missing resolution":                   func(r *GenerationRequest) { r.Resolution = "" },
		"unknown resolution":                   func(r *GenerationRequest) { r.Resolution = "4K" },
		"missing duration":                     func(r *GenerationRequest) { r.Duration = 0 },
		"short duration":                       func(r *GenerationRequest) { r.Duration = 1 },
		"long duration":                        func(r *GenerationRequest) { r.Duration = 31 },
		"pricing selector is not an operation": func(r *GenerationRequest) { r.Operation = "video_to_video" },
		"fractional seed": func(r *GenerationRequest) {
			r.ProviderOptions = map[string]map[string]any{grsaiWanProvider: {"seed": 1.5}}
		},
	} {
		t.Run(name, func(t *testing.T) {
			req := grsaiWanRequest()
			mutate(&req)
			if _, err := adapter.BuildCreate(t.Context(), RequestContext{Request: req}); err == nil {
				t.Fatal("invalid request accepted")
			}
		})
	}
}

func TestGrsaiWanVideoMedia(t *testing.T) {
	adapter := officialPackageAdapter(t, grsaiWanProvider+".yingce-plugin", grsaiWanProvider)
	image := MediaReference{Kind: "image", Role: "reference_image", URL: "https://media.example/reference.png", Order: 1}
	inline := MediaReference{Kind: "image", Role: "reference_image", DataURL: "data:image/png;base64,aW1n", Order: 2}
	video := MediaReference{Kind: "video", Role: "reference_video", URL: "https://media.example/v.mp4", MIMEType: "video/mp4", Metadata: map[string]any{"durationMs": 7500}}
	audio := MediaReference{Kind: "audio", Role: "reference_audio", URL: "https://media.example/s.mp3", MIMEType: "audio/mpeg", Metadata: map[string]any{"durationMs": 15000}}
	build := func(inputs ...MediaReference) (RequestSpec, error) {
		req := grsaiWanRequest()
		req.Inputs = inputs
		req.ProviderOptions = map[string]map[string]any{grsaiWanProvider: {"seed": 0}}
		return adapter.BuildCreate(t.Context(), RequestContext{Request: req})
	}
	spec, err := build(inline, audio, image, video)
	if err != nil {
		t.Fatal(err)
	}
	body := manifestTestBody(t, spec)
	if !reflect.DeepEqual(body["images"], []any{image.URL, inline.DataURL}) ||
		!reflect.DeepEqual(body["audios"], []any{audio.URL}) ||
		!reflect.DeepEqual(body["videos"], []any{map[string]any{"url": video.URL}}) ||
		body["seed"] != float64(0) || body["firstFrame"] != nil || body["lastFrame"] != nil {
		t.Fatalf("reference payload = %#v", body)
	}
	first, last := image, image
	first.Role, first.URL, last.Role, last.URL = "first_frame", "https://media.example/first.png", "last_frame", "https://media.example/last.png"
	frameSpec, err := build(last, first)
	if err != nil {
		t.Fatal(err)
	}
	frames := manifestTestBody(t, frameSpec)
	if frames["firstFrame"] != first.URL || frames["lastFrame"] != last.URL || frames["images"] != nil || frames["audios"] != nil {
		t.Fatalf("keyframe payload = %#v", frames)
	}
	notMP4, longAudio, shortClip, shortAudio := video, audio, video, audio
	notMP4.URL, notMP4.MIMEType = "https://media.example/v.webm", "video/webm"
	longAudio.Metadata = map[string]any{"durationMs": 16000}
	shortClip.Metadata = map[string]any{"durationMs": 2500}
	shortAudio.Metadata = map[string]any{"durationMs": 2000}
	repeat := func(ref MediaReference, n int) []MediaReference {
		out := make([]MediaReference, n)
		for i := range out {
			out[i] = ref
		}
		return out
	}
	for _, tc := range []struct {
		field  string
		inputs []MediaReference
	}{
		{"videos", repeat(video, 2)}, // exactly 15s in total
		{"videos", repeat(shortClip, 5)},
		{"audios", repeat(shortAudio, 5)},
		{"images", repeat(image, 10)},
	} {
		spec, err := build(tc.inputs...)
		if err != nil {
			t.Fatalf("%d %s: %v", len(tc.inputs), tc.field, err)
		}
		items, _ := manifestTestBody(t, spec)[tc.field].([]any)
		if len(items) != len(tc.inputs) {
			t.Fatalf("%s count = %d, want %d", tc.field, len(items), len(tc.inputs))
		}
	}
	for name, inputs := range map[string][]MediaReference{
		"frame and reference images": {first, image}, "frames and audio": {first, last, audio},
		"tail only is not offered by this release": {last}, "keyframes with video not offered": {first, video},
		"not MP4": {notMP4}, "videos over 15s in total": {video, video, video},
		"oversize audio duration": {longAudio}, "11 images": repeat(image, 11),
		"6 audios": repeat(shortAudio, 6), "6 videos within 15s": repeat(shortClip, 6),
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := build(inputs...); err == nil {
				t.Fatal("invalid media inputs accepted")
			}
		})
	}
}

func TestGrsaiWanVideoLifecycle(t *testing.T) {
	adapter := officialPackageAdapter(t, grsaiWanProvider+".yingce-plugin", grsaiWanProvider)
	created, err := adapter.ParseCreate(t.Context(), []byte(`{"id":"grs-task","status":"running"}`))
	if err != nil || created.TaskID != "grs-task" || created.Status != StatusProcessing || created.Result != nil {
		t.Fatalf("create result = %#v, err = %v", created, err)
	}
	poll := PollContext{TaskID: created.TaskID}
	running, err := adapter.ParsePoll(t.Context(), poll, []byte(`{"id":"grs-task","status":"running","progress":25}`))
	if err != nil || running.Status != StatusProcessing || running.Result != nil || !strings.Contains(running.Stage, "25%") {
		t.Fatalf("running = %#v, err = %v", running, err)
	}
	for _, tc := range []struct{ status, message string }{{"failed", "generate failed"}, {"violation", "content rejected"}} {
		// These are parsed response bodies (e.g. HTTP 200 on poll); a transport HTTP
		// 400 is rejected by the host before this adapter receives a response.
		raw := []byte(`{"id":"grs-task","status":"` + tc.status + `","error":"` + tc.message + `"}`)
		state, err := adapter.ParsePoll(t.Context(), poll, raw)
		if err != nil || state.Status != StatusFailed || state.TaskID != created.TaskID || !strings.Contains(state.Message, tc.message) || state.Result != nil {
			t.Fatalf("%s = %#v, err = %v", tc.status, state, err)
		}
	}
	unknown, err := adapter.ParsePoll(t.Context(), poll, []byte(`{"id":"grs-task","status":"unexpected"}`))
	if err == nil && unknown.Status != StatusFailed {
		t.Fatalf("unrecognized state must not poll indefinitely: %#v", unknown)
	}
	for _, bad := range []string{
		`{"id":"another-task","status":"succeeded","results":[{"url":"https://media.example/clip.mp4"}]}`,
		`{"id":"grs-task","status":"succeeded","results":[]}`,
		`{"id":"grs-task","status":"succeeded","results":{"url":"https://media.example/clip.mp4"}}`,
		`{"id":"grs-task","status":"succeeded","results":[{"url":"file:///clip.mp4"}]}`,
	} {
		state, err := adapter.ParsePoll(t.Context(), poll, []byte(bad))
		if err == nil && (state.Status != StatusFailed || state.Result != nil || state.Message == "") {
			t.Fatalf("malformed success must not expose a video or keep polling: %+v", state)
		}
	}
	const success = `{"id":"grs-task","status":"succeeded","progress":100,"results":[{"url":"https://media.example/video.mp4"}]`
	state, err := adapter.ParsePoll(t.Context(), poll, []byte(success+"}"))
	if err != nil || state.Status != StatusSucceeded || state.Result == nil || len(state.Result.Videos) != 1 ||
		state.Result.Videos[0].URL != "https://media.example/video.mp4" || len(state.Result.Usage) != 0 {
		t.Fatalf("success without usage = %#v, err = %v", state, err)
	}
	// Usage is shown in the synchronous example. Preserve it if actually present;
	// do not synthesize it on async responses or interpret provider credits as tokens.
	withUsage, err := adapter.ParseCreate(t.Context(), []byte(success+`,"usage":{"input_video_duration":6,"input_image_count":0,"duration":10,"resolution":"480p","credits":62400}}`))
	wantUsage := map[string]any{"input_video_duration": float64(6), "input_image_count": float64(0), "duration": float64(10), "resolution": "480p", "credits": float64(62400)}
	if err != nil || withUsage.Result == nil || withUsage.Status != StatusSucceeded || !reflect.DeepEqual(withUsage.Result.Usage, wantUsage) {
		t.Fatalf("success with usage = %#v, err = %v", withUsage, err)
	}
}
