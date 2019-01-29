package dropbox

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// Client implements a Dropbox client. You may use the Files and Users
// clients directly if preferred, however Client exposes them both.
type Client struct {
	*Config
	Users   *Users
	Files   *Files
	Sharing *Sharing
	Paper   *Paper
}

// New client.
func New(config *Config) *Client {
	c := &Client{Config: config}
	c.Users = &Users{c}
	c.Files = &Files{c}
	c.Sharing = &Sharing{c}
	c.Paper = &Paper{c}
	return c
}

// call rpc style endpoint.
func (c *Client) call(
	ctx context.Context,
	path string,
	in interface{},
) (io.ReadCloser, error) {
	url := "https://api.dropboxapi.com/2" + path

	body, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	req.Header.Set("Content-Type", "application/json")
	if c.Namespace != nil {
		namespaceHeader, err := json.Marshal(c.Namespace)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Dropbox-API-Path-Root", string(namespaceHeader))
	}

	r, _, err := c.do(req)
	return r, err
}

// download style endpoint.
func (c *Client) download(
	ctx context.Context,
	subdomain string,
	path string,
	in interface{},
	r io.Reader,
) (io.ReadCloser, int64, error) {
	url := fmt.Sprintf("https://%s.dropboxapi.com/2/%s", subdomain, path)

	body, err := json.Marshal(in)
	if err != nil {
		return nil, 0, err
	}

	req, err := http.NewRequest("POST", url, r)
	if err != nil {
		return nil, 0, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	req.Header.Set("Dropbox-API-Arg", string(body))
	if c.Namespace != nil {
		namespaceHeader, err := json.Marshal(c.Namespace)
		if err != nil {
			return nil, 0, err
		}
		req.Header.Set("Dropbox-API-Path-Root", string(namespaceHeader))
	}

	if r != nil {
		req.Header.Set("Content-Type", "application/octet-stream")
	}

	return c.do(req)
}

// perform the request.
func (c *Client) do(req *http.Request) (io.ReadCloser, int64, error) {
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, 0, err
	}

	if res.StatusCode < 400 {
		return res.Body, res.ContentLength, err
	}

	defer res.Body.Close()

	e := &Error{
		Status:     http.StatusText(res.StatusCode),
		StatusCode: res.StatusCode,
		Header:     res.Header,
	}

	kind := res.Header.Get("Content-Type")

	if !strings.Contains(kind, "json") {
		if b, err := ioutil.ReadAll(res.Body); err == nil {
			e.Summary = string(b)
			return nil, 0, e
		}
		return nil, 0, err
	}

	if err := json.NewDecoder(res.Body).Decode(e); err != nil {
		return nil, 0, err
	}

	return nil, 0, e
}
