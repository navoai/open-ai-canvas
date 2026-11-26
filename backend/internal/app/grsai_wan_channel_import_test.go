package app

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"yingce/backend/internal/model"
)

func TestGrsaiWanChannelExampleImportsDisabledAndUnpriced(t *testing.T) {
	// Isolate the official catalog so this exercises uploading the distributed
	// package to an existing deployment rather than a preloaded test adapter.
	t.Setenv("CANVAS_OFFICIAL_PLUGIN_DIR", t.TempDir())
	service, _ := newChannelModelTestService(t)
	const provider = "grsai-wan-video"
	packages := filepath.Join("..", "..", "..", "plugin-packages")
	archive, err := os.ReadFile(filepath.Join(packages, provider+".yingce-plugin"))
	if err != nil {
		t.Fatal(err)
	}
	service.pluginRuntime, err = newPluginRuntime(service.dataDir)
	if err != nil {
		t.Fatal(err)
	}
	installed, err := service.InstallPlugin(archive, provider+".yingce-plugin")
	if err != nil || installed.SHA256 == "" {
		t.Fatalf("install package: digest=%q, err=%v", installed.SHA256, err)
	}

	raw, err := os.ReadFile(filepath.Join(packages, provider, "config", provider+"-channel.json"))
	if err != nil {
		t.Fatal(err)
	}
	var example SystemChannelExport
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields() // a credential field silently ignored by import is still unsafe to ship
	if err := decoder.Decode(&example); err != nil {
		t.Fatal(err)
	}
	if len(example.Channels) != 1 || example.Channels[0].Enabled || len(example.Channels[0].Models) != 1 {
		t.Fatal("expected a single disabled channel containing one model")
	}
	input := example.Channels[0].Models[0]
	preset := officialVideoPreset(provider, "wan-3.0")
	if preset == nil || input.CapabilityConfig == nil || !reflect.DeepEqual(input.CapabilityConfig.Video, preset) {
		t.Fatal("importable example must declare the same video capability as the shared model preset")
	}

	admin := &model.User{ID: "admin", Role: model.UserRoleAdmin}
	result, err := service.ImportSystemChannels(admin, example)
	if err != nil || len(result.Errors) != 0 || len(result.CreatedChannels) != 1 || result.SavedModels != 1 {
		t.Fatalf("import: result=%+v, err=%v", result, err)
	}
	channels, err := service.repo.SystemChannels(true)
	if err != nil || len(channels) != 1 {
		t.Fatalf("got %d channels, err=%v", len(channels), err)
	}
	channel := channels[0]
	if channel.Enabled || channel.APIKey != "" || channel.SecretKey != "" || channel.BaseURL != "https://grsaiapi.com" {
		t.Fatal("new channel must use the global node and be disabled without credentials")
	}
	models, err := service.repo.ChannelModels(channel.ID, true)
	if err != nil || len(models) != 1 {
		t.Fatalf("got %d models, err=%v", len(models), err)
	}
	item := models[0]
	if item.Enabled || item.PriceConfigured || !item.ReferenceVideoBilling || item.ModelKey != "wan-3.0" ||
		item.ProviderModelKey != "wan-3.0" || string(item.Protocol) != provider || item.ProtocolRelease != installed.SHA256 {
		t.Fatal("model must be disabled, unpriced, support reference-video billing and resolve the installed release")
	}
	saved, err := DecodeModelCapabilityConfig(item.CapabilityConfigJSON)
	if err != nil || saved == nil || saved.Video == nil || !reflect.DeepEqual(saved.Video.Operations, preset.Operations) ||
		!reflect.DeepEqual(saved.Video.References, preset.References) || len(saved.Video.ModeLimits) != len(preset.Operations) {
		t.Fatalf("imported video capability or installed protocol mode limits missing: %v", err)
	}
	if len(item.PriceTiers) != 3 {
		t.Fatalf("got %d price tiers, want 3", len(item.PriceTiers))
	}
	missing := map[string]bool{"480p": true, "720p": true, "1080p": true}
	for _, tier := range item.PriceTiers {
		if !missing[tier.Resolution] || tier.Selector["vquality"] != tier.Resolution || tier.VideoSeconds != 0 ||
			tier.BillingMode != "per_second" || tier.ProviderModelKey != "wan-3.0" || tier.Enabled || tier.PriceConfigured ||
			tier.UnitPriceMicrocredits != 0 || tier.ReferenceVideoUnitPriceMicrocredits != 0 || tier.CostPricing.Configured {
			t.Fatalf("unexpected imported price tier: %#v", tier)
		}
		delete(missing, tier.Resolution)
	}
}
