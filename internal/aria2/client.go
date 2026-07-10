package aria2

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Client struct {
	baseURL string
	secret  string
	http    *http.Client
}

func NewClient(port int, secret string) *Client {
	return &Client{
		baseURL: fmt.Sprintf("http://127.0.0.1:%d/jsonrpc", port),
		secret:  secret,
		http:    &http.Client{Timeout: 10 * time.Second},
	}
}

type rpcRequest struct {
	JSONRPC string `json:"jsonrpc"`
	ID      string `json:"id"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
}

type rpcResponse struct {
	ID     string          `json:"id"`
	Result json.RawMessage `json:"result"`
	Error  *rpcError       `json:"error"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (c *Client) call(method string, params ...any) (json.RawMessage, error) {
	rpcParams := make([]any, 0, len(params)+1)
	if c.secret != "" {
		rpcParams = append(rpcParams, "token:"+c.secret)
	}
	rpcParams = append(rpcParams, params...)

	reqBody, err := json.Marshal(rpcRequest{
		JSONRPC: "2.0",
		ID:      "adl",
		Method:  method,
		Params:  rpcParams,
	})
	if err != nil {
		return nil, err
	}
	resp, err := c.http.Post(c.baseURL, "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var rpcResp rpcResponse
	if err := json.Unmarshal(body, &rpcResp); err != nil {
		return nil, err
	}
	if rpcResp.Error != nil {
		return nil, fmt.Errorf("aria2: %s", rpcResp.Error.Message)
	}
	return rpcResp.Result, nil
}

func (c *Client) AddURI(uri string, dir string) (string, error) {
	opts := map[string]string{"dir": dir}
	result, err := c.call("aria2.addUri", []string{uri}, opts)
	if err != nil {
		return "", err
	}
	var gid string
	if err := json.Unmarshal(result, &gid); err != nil {
		return "", err
	}
	return gid, nil
}

func (c *Client) AddTorrent(path string, dir string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	encoded := base64.StdEncoding.EncodeToString(data)
	opts := map[string]string{"dir": dir}
	result, err := c.call("aria2.addTorrent", encoded, []any{}, opts)
	if err != nil {
		return "", err
	}
	var gid string
	if err := json.Unmarshal(result, &gid); err != nil {
		return "", err
	}
	return gid, nil
}

type Status struct {
	GID       string
	Name      string
	Status    string
	Total     int64
	Completed int64
	Speed     int64
	ETA       int64
	Dir       string
}

func (c *Client) TellActive() ([]Status, error) {
	result, err := c.call("aria2.tellActive")
	if err != nil {
		return nil, err
	}
	var raw []map[string]any
	if err := json.Unmarshal(result, &raw); err != nil {
		return nil, err
	}
	out := make([]Status, 0, len(raw))
	for _, item := range raw {
		st, err := parseStatus(item)
		if err != nil {
			continue
		}
		out = append(out, st)
	}
	return out, nil
}

func (c *Client) TellStatus(gid string) (Status, error) {
	result, err := c.call("aria2.tellStatus", gid)
	if err != nil {
		return Status{}, err
	}
	var raw map[string]any
	if err := json.Unmarshal(result, &raw); err != nil {
		return Status{}, err
	}
	return parseStatus(raw)
}

func (c *Client) Pause(gid string) error {
	_, err := c.call("aria2.pause", gid)
	return err
}

func (c *Client) Unpause(gid string) error {
	_, err := c.call("aria2.unpause", gid)
	return err
}

func (c *Client) Remove(gid string) error {
	_, err := c.call("aria2.remove", gid)
	return err
}

func parseStatus(raw map[string]any) (Status, error) {
	st := Status{
		GID:    str(raw, "gid"),
		Status: str(raw, "status"),
		Dir:    str(raw, "dir"),
	}
	if files, ok := raw["files"].([]any); ok && len(files) > 0 {
		if f, ok := files[0].(map[string]any); ok {
			st.Name = str(f, "path")
		}
	}
	st.Total = atoi(str(raw, "totalLength"))
	st.Completed = atoi(str(raw, "completedLength"))
	st.Speed = atoi(str(raw, "downloadSpeed"))
	if st.Speed > 0 && st.Total > st.Completed {
		st.ETA = (st.Total - st.Completed) / st.Speed
	}
	return st, nil
}

func str(m map[string]any, key string) string {
	v, ok := m[key]
	if !ok {
		return ""
	}
	switch t := v.(type) {
	case string:
		return t
	default:
		return fmt.Sprint(t)
	}
}

func atoi(s string) int64 {
	n, _ := strconv.ParseInt(s, 10, 64)
	return n
}
