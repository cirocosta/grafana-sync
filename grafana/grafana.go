package grafana

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"path"

	"github.com/pkg/errors"
)

const (
	DashboardSearch         = "/api/search"
	DashboardCreateOrUpdate = "/api/dashboards/db"
	DashboardByUid          = "/api/dashboards/uid"
	FoldersApi              = "/api/folders"
)

type ClientConfig struct {
	Verbose     bool
	Address     string
	AccessToken string
	Username    string
	Password    string
}

type Client struct {
	client *http.Client
	cfg    ClientConfig
}

func NewClient(opts ClientConfig) (c *Client) {
	c = &Client{
		client: &http.Client{},
		cfg:    opts,
	}
	return
}

type DashboardCreateOrUpdateEntry struct {
	Overwrite bool                   `json:"overwrite"`
	FolderId  int                    `json:"folderId"`
	Dashboard map[string]interface{} `json:"dashboard"`
}

type DashboardSearchEntry struct {
	Uid    string `json:"uid"`
	Title  string `json:"title"`
	Folder string `json:"folderTitle"`
}

func (c *Client) doRequest(req *http.Request) (resp *http.Response, err error) {
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

type Folder struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Uid   string `json:"uid"`
}

func (c *Client) GetDashboard(ctx context.Context, uid string) (dashboard map[string]interface{}, err error) {
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

	var (
		decoder  = json.NewDecoder(resp.Body)
		contents = map[string]interface{}{}
	)

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

	dashboard, _ = innerDashboard.(map[string]interface{})
	dashboard["id"] = nil

	return
}

func (c *Client) CreateFolder(ctx context.Context, title string) (folder Folder, err error) {
	var content []byte

	content, err = json.Marshal(&Folder{Title: title})
	if err != nil {
		err = errors.Wrapf(err,
			"failed to encode folder as json")
		return
	}

	req, err := http.NewRequest("POST",
		c.cfg.Address+FoldersApi, bytes.NewBuffer(content))
	if err != nil {
		err = errors.Wrapf(err,
			"couldn't prepare request")
		return
	}

	req.Header.Set("content-type", "application/json")

	req = req.WithContext(ctx)

	resp, err := c.doRequest(req)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = errors.Errorf("non successful response - %s", resp.Status)
		return
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(folder)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to decode folder creation response")
		return
	}

	return
}

func (c *Client) ListFolders(ctx context.Context) (folders []Folder, err error) {
	req, err := http.NewRequest("GET", c.cfg.Address+FoldersApi, nil)
	if err != nil {
		err = errors.Wrapf(err,
			"couldn't prepare request")
		return
	}

	q := req.URL.Query()
	q.Add("limit", "100")

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
	err = decoder.Decode(&folders)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to decode dashboards ref json")
		return
	}

	return
}

func (c *Client) PushDashboard(ctx context.Context, dashboard DashboardCreateOrUpdateEntry) (err error) {
	var content []byte

	content, err = json.Marshal(dashboard)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to encode dashboard as json")
		return
	}

	req, err := http.NewRequest("POST",
		c.cfg.Address+DashboardCreateOrUpdate, bytes.NewBuffer(content))
	if err != nil {
		err = errors.Wrapf(err,
			"couldn't prepare request")
		return
	}

	req.Header.Set("content-type", "application/json")

	req = req.WithContext(ctx)

	resp, err := c.doRequest(req)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = errors.Errorf("non successful response - %s", resp.Status)
		return
	}

	return
}

func (c *Client) ListDashboardSearchEntries(ctx context.Context) (dashboards []*DashboardSearchEntry, err error) {
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
