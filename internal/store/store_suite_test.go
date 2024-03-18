package store_test

import (
	"errors"
	"os"
	"path"
	"testing"

	"github.com/drew-english/system-configurator/internal/model"
	"github.com/drew-english/system-configurator/internal/store"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestStore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Store Suite")
}

var _ = Describe("Local", func() {
	var (
		cfg        *store.LocalCfg
		cfgFixture string
	)

	BeforeEach(func() {
		cfgFixture = `{
			"packages": [
				{
					"name": "some-package",
					"version": "1.2.3",
					"alternates": {
						"apt": {
							"name": "somepackage",
							"version": "1.2.3"
						}
					}
				}
			]
		}`

		cfg = &store.LocalCfg{
			Location: "./tmp/system-configurator",
			FileName: "config.json",
		}
	})

	JustBeforeEach(func() {
		if cfgFixture != "" {
			os.MkdirAll("./tmp/system-configurator", 0755)
			f, _ := os.Create("./tmp/system-configurator/config.json")
			f.WriteString(cfgFixture)
			f.Close()
		}
	})

	AfterEach(func() {
		os.RemoveAll("./tmp")
	})

	Describe("NewLocal", func() {
		subject := func() (store.Store, error) {
			return store.NewLocal(cfg)
		}

		It("returns a new local store", func() {
			s, _ := subject()
			Expect(s).ToNot(BeNil())
		})

		Context("when the local config file does not exist", func() {
			BeforeEach(func() {
				cfgFixture = ""
			})

			It("creates a new one", func() {
				_, err := os.Open("./tmp/system-configurator/config.json")
				Expect(errors.Is(err, os.ErrNotExist)).To(BeTrue())

				subject()
				_, err = os.Open("./tmp/system-configurator/config.json")
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when no configuration is given", func() {
			BeforeEach(func() {
				cfgFixture = ""
				cfg = nil
			})

			AfterEach(func() {
				os.RemoveAll("~") // safe in tests as it references the cwd
			})

			It("loads the config from the default location", func() {
				subject()
				_, err := os.Open(path.Join(store.LocalDefaultLocation, store.LocalDefaultFileName))
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe(".LoadConfiguration", func() {
		var localStore store.Store

		JustBeforeEach(func() {
			localStore, _ = store.NewLocal(cfg)
		})

		subject := func() (*store.Configuration, error) {
			return localStore.LoadConfiguration()
		}

		It("returns the configuration specified in the file", func() {
			config, err := subject()
			Expect(err).ToNot(HaveOccurred())
			Expect(config).To(BeComparableTo(&store.Configuration{
				Packages: []*model.Package{
					{
						Name:    "some-package",
						Version: "1.2.3",
						Alternates: map[string]*model.Package{
							"apt": {
								Name:    "somepackage",
								Version: "1.2.3",
							},
						},
					},
				},
			}))
		})

		Context("when the config file was not loaded correctly", func() {
			BeforeEach(func() {
				f, err := os.Create("./a-file.tmp")
				Expect(err).ToNot(HaveOccurred())
				f.Close()

				cfg = &store.LocalCfg{
					Location: "./a-file.tmp",
				}
			})

			AfterEach(func() {
				os.Remove("./a-file.tmp")
			})

			It("returns an error", func() {
				_, err := subject()
				Expect(err).To(MatchError("error referencing local configuration file"))
			})
		})
	})
})
