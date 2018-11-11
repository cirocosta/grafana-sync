package grafana

import (
	"os"
	"encoding/json"

	"github.com/pkg/errors"
)

type Dashboard map[string]interface{}

func (d *Dashboard) SaveToDisk(filepath string) (err error) {
	var dashboardFile *os.File

	dashboardFile, err = os.Create(filepath)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to create file for dashboard %s", filepath)
		return

	}

	defer dashboardFile.Close()

	err = json.NewEncoder(dashboardFile).Encode(d)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to encode dashboard json to file %s", filepath)
		return
	}

	return
}

