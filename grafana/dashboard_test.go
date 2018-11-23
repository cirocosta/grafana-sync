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
		})

		Context("not having panels", func() {
			BeforeEach(func() {
				dashboard = map[string]interface{}{"foo": "bar"}
				expected = map[string]interface{}{"foo": "bar"}
			})

			It("works", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(dashboard).To(Equal(expected))
			})
		})

		Context("being a panel", func() {
			BeforeEach(func() {
				dashboard = map[string]interface{}{"foo": "bar", "datasource": "ds1"}
				expected = map[string]interface{}{"foo": "bar", "datasource": "ds2"}
				datasource = "ds2"
			})

			It("works", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(dashboard).To(Equal(expected))
			})
		})

		Context("having nil panels list", func() {
			BeforeEach(func() {
				dashboard = map[string]interface{}{
					"foo":    "bar",
					"panels": nil,
				}
				expected = map[string]interface{}{
					"foo":    "bar",
					"panels": nil,
				}
				datasource = "ds2"
			})

			It("does nothing", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(dashboard).To(Equal(expected))
			})
		})

		Context("being a row with nested panels", func() {
			BeforeEach(func() {
				dashboard = map[string]interface{}{
					"foo": "bar",
					"panels": []map[string]interface{}{
						{
							"datasource": "ds1",
						},
					},
				}
				expected = map[string]interface{}{
					"foo": "bar",
					"panels": []map[string]interface{}{
						{
							"datasource": "ds2",
						},
					},
				}
				datasource = "ds2"
			})

			It("works", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(dashboard).To(Equal(expected))
			})
		})
	})
})
