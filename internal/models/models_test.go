package models

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/tools/go/packages"
	"strings"
)

var _ = Describe("Models", func() {
	When("Models is executed", func() {
		It("should return a slice containing a pointer to each model", func() {
			models := Models()
			Expect(models).Should(Equal([]interface{}{
				&Email{}, &Group{}, &User{},
			}))
		})
		It("should return a slice with a length equal to the number of models defined", func() {
			cfg := &packages.Config{
				Mode: packages.NeedImports | packages.NeedTypes,
			}
			pkgs, err := packages.Load(cfg, "github.com/golangbb/golangbb/v2/internal/models")
			Expect(err).Should(BeNil())
			Expect(len(pkgs)).Should(Equal(1))

			pkg := pkgs[0]
			scope := pkg.Types.Scope()
			numberOfExportedModels := 0
			for _, name := range scope.Names() {
				obj := scope.Lookup(name)
				if !obj.Exported() {
					continue
				}

				if !strings.HasPrefix(obj.String(), "type") {
					continue
				}

				numberOfExportedModels += 1
			}

			numberOfModels := len(Models())
			Expect(numberOfModels).Should(Equal(numberOfExportedModels))
		})
	})
})
