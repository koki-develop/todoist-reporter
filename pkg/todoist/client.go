package todoist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/koki-develop/todoist-reporter/pkg/util"
	"github.com/pkg/errors"
)

type HTTPAPI interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	token   string
	httpAPI HTTPAPI
}

func New(token string) *Client {
	return &Client{
		token:   token,
		httpAPI: new(http.Client),
	}
}

func (cl *Client) GetResources(types []string) (*Resources, error) {
	p, err := json.Marshal(map[string]interface{}{
		"sync_token":     "*",
		"resource_types": types,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ep := "https://api.todoist.com/sync/v8/sync"
	req, err := http.NewRequest(http.MethodPost, ep, bytes.NewBuffer(p))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	req.Header.Set("authorization", fmt.Sprintf("Bearer %s", cl.token))
	req.Header.Set("content-type", "application/json")

	resp, err := cl.httpAPI.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer resp.Body.Close()

	var r Resources
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, errors.WithStack(err)
	}
	if util.Contains(types, "items") {
		r.Items.Organize()
	}

	return &r, nil
}
