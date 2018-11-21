package grafana

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
)

type Dashboard map[string]interface{}

func SetPanelDatasources(panel map[string]interface{}, datasource string) (err error) {
	for k, v := range panel {
		switch k {
		case "datasource":
			panel[k] = datasource

		case "panels":
			panels, ok := v.([]map[string]interface{})
			if !ok {
				err = errors.Errorf("can't convert panel list")
				return
			}

			for _, panel := range panels {
				err = SetPanelDatasources(panel, datasource)
				if err != nil {
					err = errors.Wrapf(err,
						"failed to set inner datasource")
					return
				}
			}

			continue
		}
	}

	return
}

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
