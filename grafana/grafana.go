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

type ClientConfig struct {
	Verbose     bool
	Address     string
	AccessToken string
	Username    string
	Password    string
}

type client struct {
	client *http.Client
	cfg    ClientConfig
}

func NewClient(opts ClientConfig) (c *client) {
	c = &client{
		client: &http.Client{},
		cfg:    opts,
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

	if c.cfg.AccessToken != "" {
		req.Header.Add("Authorization", "Bearer "+c.cfg.AccessToken)
	} else if c.cfg.Username != "" {
		req.SetBasicAuth(c.cfg.Username, c.cfg.Password)
	}

	if c.cfg.Verbose {
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

	if c.cfg.Verbose {
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
	req, err := http.NewRequest("GET", c.cfg.Address+path.Join(DashboardByUid, uid), nil)
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
	contents := map[string]interface{}{}

	err = decoder.Decode(&contents)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to decode dashboard json")
		return
	}

	innerDashboard, found := contents["dashboard"]
	if !found {
		err = errors.Errorf("dashboard not found in json response")
		return
	}

	dashboard, _ = innerDashboard.(Dashboard)
	dashboard["id"] = nil

	return
}

func (c *client) ListDashboardRefs(ctx context.Context) (dashboards []*DashboardRef, err error) {
	req, err := http.NewRequest("GET", c.cfg.Address+DashboardSearch, nil)
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
