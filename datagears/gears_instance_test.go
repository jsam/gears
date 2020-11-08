package datagears_test

import (
	"github.com/jsam/dg/datagears"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GearsInstance", func() {

	Context("gears instance builder", func() {
		It("check expected object via builder", func() {
			gearsInstance, err := datagears.FromRedisDSN("redis://password@redis:6379/1")
			Expect(err).ShouldNot(HaveOccurred())

			Expect(gearsInstance.Password).To(Equal("password"))
			Expect(gearsInstance.DB).To(Equal(1))
			Expect(gearsInstance.Addr).To(Equal("redis:6379"))
		})
	})
})
