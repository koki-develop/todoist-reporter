package todoist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type httpAPI interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	token   string
	httpAPI httpAPI
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

	if resp.StatusCode != http.StatusOK {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		return nil, errors.New(string(b))
	}

	var r Resources
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, errors.WithStack(err)
	}

	return &r, nil
}

func (cl *Client) GetCompletedItems(projID int, since time.Time) (Items, error) {
	p, err := json.Marshal(map[string]interface{}{
		"project_id": projID,
		"limit":      200,
		"since":      since,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ep := "https://api.todoist.com/sync/v8/completed/get_all"
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

	if resp.StatusCode != http.StatusOK {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		return nil, errors.New(string(b))
	}

	var r getCompletedItemsResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, errors.WithStack(err)
	}

	var items Items
	eg := new(errgroup.Group)
	for _, item := range r.Items {
		item := item // https://golang.org/doc/faq#closures_and_goroutines

		eg.Go(func() error {
			p, err := json.Marshal(map[string]interface{}{
				"item_id": item.TaskID,
			})
			if err != nil {
				return errors.WithStack(err)
			}
			ep := "https://api.todoist.com/sync/v8/items/get"
			req, err := http.NewRequest(http.MethodPost, ep, bytes.NewBuffer(p))
			if err != nil {
				return errors.WithStack(err)
			}
			req.Header.Set("authorization", fmt.Sprintf("Bearer %s", cl.token))
			req.Header.Set("content-type", "application/json")

			resp, err := cl.httpAPI.Do(req)
			if err != nil {
				return errors.WithStack(err)
			}
			defer resp.Body.Close()

			// 何故か 404 が返ってくることがある。どうしようもないので無視する。
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}

			if resp.StatusCode != http.StatusOK {
				b, err := io.ReadAll(resp.Body)
				if err != nil {
					return errors.WithStack(err)
				}
				return errors.New(string(b))
			}
			var info getItemInfoResponse
			if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
				return errors.WithStack(err)
			}
			info.Item.Completed = true
			items = append(items, info.Item)

			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, errors.WithStack(err)
	}

	return items, nil
}
