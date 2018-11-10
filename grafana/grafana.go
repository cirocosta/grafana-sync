package grafana

import (
	"context"
	"net/http"
	"encoding/json"

	"github.com/pkg/errors"
)

const (
	DashboardSearch = "/api/search"
)

type client struct {
	address string
	client  *http.Client
}

func NewClient(address string) (c *client) {
	c = &client{
		address: address,
		client:  &http.Client{},
	}
	return
}

type Dashboard struct {
	Title  string `json:"title"`
	Folder string `json:"folderTitle"`
}

func (c *client) AllDashboards(ctx context.Context) (dashboards []*Dashboard, err error) {
	req, err := http.NewRequest("GET", c.address+DashboardSearch, nil)
	if err != nil {
		err = errors.Wrapf(err,
			"couldn't prepare request")
		return
	}

	req = req.WithContext(ctx)
	resp, err := c.client.Do(req)
	if err != nil {
		err = errors.Wrapf(err,
			"failed while performing request")
		return
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(dashboards)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to decode dashboards")
		return
	}

	return
}
