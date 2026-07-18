package aria2

import (
	"bytes"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"os/exec"
	"reflect"
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

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestTellDownloadsIncludesPausedWaitingItems(t *testing.T) {
	results := map[string]any{
		"aria2.tellActive": []map[string]any{{
			"gid": "active", "status": "active", "files": []any{},
		}},
		"aria2.tellWaiting": []map[string]any{
			{"gid": "active", "status": "waiting", "files": []any{}},
			{"gid": "paused", "status": "paused", "files": []any{}},
		},
	}
	var requests []rpcRequest
	c := NewClient(6800, "")
	c.http.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		var rpcReq rpcRequest
		if err := json.NewDecoder(req.Body).Decode(&rpcReq); err != nil {
			t.Fatal(err)
		}
		requests = append(requests, rpcReq)
		body, err := json.Marshal(map[string]any{
			"jsonrpc": "2.0",
			"id":      rpcReq.ID,
			"result":  results[rpcReq.Method],
		})
		if err != nil {
			t.Fatal(err)
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})

	got, err := c.TellDownloads()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0].GID != "active" || got[1].GID != "paused" {
		t.Fatalf("downloads = %#v, want active and paused without duplicates", got)
	}

	methods := []string{requests[0].Method, requests[1].Method}
	if !reflect.DeepEqual(methods, []string{"aria2.tellActive", "aria2.tellWaiting"}) {
		t.Fatalf("RPC methods = %v", methods)
	}
	if len(requests[0].Params) != 1 || len(requests[1].Params) != 3 {
		t.Fatalf("unexpected RPC params: %#v", requests)
	}
}

func TestRPCGetVersionIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	binary, err := exec.LookPath("aria2c")
	if err != nil {
		t.Skip("aria2c is not installed")
	}
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	if err := listener.Close(); err != nil {
		t.Fatal(err)
	}

	d, err := Start(binary, port, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer d.Stop()

	_, err = d.Client().call("aria2.getVersion")
	if err != nil {
		t.Fatal(err)
	}
}
