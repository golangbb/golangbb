package helpers

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("GetEnv", func() {
	var envVarKey = "envVarKey"
	var envVarValue = "envVarValue"
	var fallbackValue = "fallbackValue"

	BeforeEach(func() {
		os.Unsetenv(envVarKey)
	})
	When("Environment Variable is set", func() {
		It("should get the Environment Variable value from the OS", func() {
			os.Setenv(envVarKey, envVarValue)
			grabbedEnvVarValue := GetEnv(envVarKey, fallbackValue)
			Expect(grabbedEnvVarValue).Should(BeIdenticalTo(envVarValue))
		})
	})
	When("Environment Variable is not set", func() {
		It("should get use the fallback value", func() {
			grabbedEnvVarValue := GetEnv(envVarKey, fallbackValue)
			Expect(grabbedEnvVarValue).Should(BeIdenticalTo(fallbackValue))
		})
	})
})
