package nzbget

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type Client struct {
	url      string
	username string
	password string
	http     *http.Client

	versionOnce  sync.Once
	cachedVersion string
}

func NewClient(host string, port int, username, password string) *Client {
	return &Client{
		url:      fmt.Sprintf("http://%s:%d/jsonrpc", host, port),
		username: username,
		password: password,
		http:     &http.Client{Timeout: 5 * time.Second},
	}
}

type rpcRequest struct {
	Method string `json:"method"`
	Params []any  `json:"params"`
	ID     int    `json:"id"`
}

type rpcResponse struct {
	Result json.RawMessage `json:"result"`
	Error  any             `json:"error"`
}

func (c *Client) call(method string, out any) error {
	body, err := json.Marshal(rpcRequest{Method: method, Params: []any{}, ID: 1})
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, c.url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.username, c.password)

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("nzbget API returned HTTP %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var rpc rpcResponse
	if err := json.Unmarshal(data, &rpc); err != nil {
		return err
	}
	if rpc.Error != nil {
		return fmt.Errorf("nzbget RPC error: %v", rpc.Error)
	}
	return json.Unmarshal(rpc.Result, out)
}

func (c *Client) Status() (*StatusResult, error) {
	var s StatusResult
	if err := c.call("status", &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func (c *Client) Version() (string, error) {
	var err error
	c.versionOnce.Do(func() {
		var v string
		if callErr := c.call("version", &v); callErr != nil {
			err = callErr
			c.versionOnce = sync.Once{} // reset so next call retries
			return
		}
		c.cachedVersion = v
	})
	if err != nil {
		return "", err
	}
	return c.cachedVersion, nil
}

// ListGroups returns the number of queued NZBs.
func (c *Client) QueuedCount() (int, error) {
	var groups []json.RawMessage
	if err := c.call("listgroups", &groups); err != nil {
		return 0, err
	}
	return len(groups), nil
}
