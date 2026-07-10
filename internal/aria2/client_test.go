package aria2

import (
	"encoding/json"
	"testing"
)

func TestRPCParamsNeverNull(t *testing.T) {
	c := NewClient(6800, "")
	body, err := json.Marshal(rpcRequest{
		JSONRPC: "2.0",
		ID:      "adl",
		Method:  "aria2.getVersion",
		Params:  make([]any, 0),
	})
	if err != nil {
		t.Fatal(err)
	}
	if !json.Valid(body) {
		t.Fatal("invalid json")
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		t.Fatal(err)
	}
	if string(raw["params"]) == "null" {
		t.Fatal("params must not be null")
	}
	_ = c
}

func TestRPCGetVersionIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	d, err := Start("aria2c", 16802, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer d.Stop()

	_, err = d.Client().call("aria2.getVersion")
	if err != nil {
		t.Fatal(err)
	}
}
