package grafana_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cirocosta/grafana-sync/grafana"
)

var _ = Describe("Dashboard", func() {
	Describe("SetPanelDatasources", func() {
		var (
			dashboard  map[string]interface{}
			expected   map[string]interface{}
			datasource string
			err        error
		)

		JustBeforeEach(func() {
			err = grafana.SetPanelDatasources(dashboard, datasource)
			Expect(err).NotTo(HaveOccurred())
			Expect(dashboard).To(Equal(expected))
		})

		Context("not having panels", func() {
			BeforeEach(func() {
				dashboard = map[string]interface{}{"foo": "bar"}
				expected = map[string]interface{}{"foo": "bar"}
			})
		})

		Context("being a panel", func() {
			BeforeEach(func() {
				dashboard = map[string]interface{}{"foo": "bar", "datasource": "ds1"}
				expected = map[string]interface{}{"foo": "bar", "datasource": "ds2"}
				datasource = "ds2"
			})
		})

		Context("being a panel", func() {
			BeforeEach(func() {
				dashboard = map[string]interface{}{"foo": "bar", "datasource": "ds1"}
				expected = map[string]interface{}{"foo": "bar", "datasource": "ds2"}
				datasource = "ds2"
			})
		})

		Context("being a row with nested panels", func() {
			BeforeEach(func() {
				expected = map[string]interface{}{
					"foo": "bar",
					"panels": map[string]interface{}{
						"datasource": "ds1",
					},
				}
				expected = map[string]interface{}{
					"foo": "bar",
					"panels": map[string]interface{}{
						"datasource": "ds2",
					},
				}
				datasource = "ds2"
			})
		})
	})
})
