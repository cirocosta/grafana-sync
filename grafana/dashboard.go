package grafana

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
)

func SetPanelDatasources(panel map[string]interface{}, datasource string) (err error) {
	for k, v := range panel {
		switch k {
		case "datasource":
			panel[k] = datasource

		case "panels":
			if v == nil {
				continue
			}

			var panels []map[string]interface{}

			switch v.(type) {
			case []interface{}:
				for _, iface := range v.([]interface{}) {
					panels = append(panels,
						iface.(map[string]interface{}))
				}
			case []map[string]interface{}:
				panels, _ = v.([]map[string]interface{})
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

func SaveToDisk(filepath string, d map[string]interface{}) (err error) {
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

func LoadFromDisk(filepath string) (dashboard map[string]interface{}, err error) {
	dashboardFile, err := os.Open(filepath)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to open dashboard file %s", filepath)
		return
	}

	defer dashboardFile.Close()

	err = json.NewDecoder(dashboardFile).Decode(&dashboard)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to read dashboard %s into memory", filepath)

	}

	return
}
