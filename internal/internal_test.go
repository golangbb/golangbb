package internal

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Constants", func() {
	Context("PORT", func() {
		When("constant is accessed without environment variable being set", func() {
			It("should return the default/fallback value", func() {
				Expect(PORT).Should(BeIdenticalTo(defaultPORT))
			})
		})
	})
	Context("DATABASENAME", func() {
		When("constant is accessed without environment variable being set", func() {
			It("should return the default/fallback value", func() {
				Expect(DATABASENAME).Should(BeIdenticalTo(defaultDATABASENAME))
			})
		})
	})
})
