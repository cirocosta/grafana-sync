package grafana_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGrafana(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Metric Suite")
}
