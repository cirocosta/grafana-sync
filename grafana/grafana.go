package grafana

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"path"

	"github.com/pkg/errors"
)

const (
	DashboardSearch = "/api/search"
	DashboardByUid  = "/api/dashboards/uid"
)

type client struct {
	address     string
	verbose     bool
	accessToken string
	client      *http.Client
}

func NewClient(address, accessToken string, verbose bool) (c *client) {
	c = &client{
		accessToken: accessToken,
		address:     address,
		client:      &http.Client{},
		verbose:     verbose,
	}
	return
}

type DashboardRef struct {
	Uid    string `json:"uid"`
	Title  string `json:"title"`
	Folder string `json:"folderTitle"`
}

func (c *client) doRequest(req *http.Request) (resp *http.Response, err error) {
	var verboseBytes []byte

	req.Header.Add("Authorization", "Bearer "+c.accessToken)

	if c.verbose {
		verboseBytes, err = httputil.DumpRequestOut(req, true)
		if err != nil {
			err = errors.Wrapf(err,
				"failed to dump verbose version of request")
			return
		}

		log.Println(string(verboseBytes))
	}

	resp, err = c.client.Do(req)
	if err != nil {
		return
	}

	if c.verbose {
		verboseBytes, err = httputil.DumpResponse(resp, true)
		if err != nil {
			err = errors.Wrapf(err,
				"failed to dump verbose version of response")
			return
		}

		log.Println(string(verboseBytes))
	}

	return
}

func (c *client) GetDashboard(ctx context.Context, uid string) (dashboard Dashboard, err error) {
	req, err := http.NewRequest("GET", c.address+path.Join(DashboardByUid, uid), nil)
	if err != nil {
		err = errors.Wrapf(err,
			"couldn't prepare request")
		return
	}

	req = req.WithContext(ctx)
	resp, err := c.doRequest(req)
	if err != nil {
		err = errors.Wrapf(err,
			"failed while performing request")
		return
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&dashboard)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to decode dashboard json")
		return
	}

	dashboard["id"] = nil

	return
}

func (c *client) ListDashboardRefs(ctx context.Context) (dashboards []*DashboardRef, err error) {
	req, err := http.NewRequest("GET", c.address+DashboardSearch, nil)
	if err != nil {
		err = errors.Wrapf(err,
			"couldn't prepare request")
		return
	}

	q := req.URL.Query()
	q.Add("type", "dash-db")

	req.URL.RawQuery = q.Encode()

	req = req.WithContext(ctx)
	resp, err := c.doRequest(req)
	if err != nil {
		err = errors.Wrapf(err,
			"failed while performing request")
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = errors.Errorf("non successful response - %s", resp.Status)
		return
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&dashboards)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to decode dashboards ref json")
		return
	}

	return
}
