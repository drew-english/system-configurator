package store_test

import (
	"errors"
	"os"
	"path"
	"testing"

	"github.com/drew-english/system-configurator/internal/model"
	"github.com/drew-english/system-configurator/internal/store"
	"github.com/google/go-cmp/cmp/cmpopts"

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
			FileName: "config.yaml",
		}
	})

	JustBeforeEach(func() {
		if cfgFixture != "" {
			os.MkdirAll("./tmp/system-configurator", 0755)
			f, _ := os.Create("./tmp/system-configurator/config.yaml")
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
				_, err := os.Open("./tmp/system-configurator/config.yaml")
				Expect(errors.Is(err, os.ErrNotExist)).To(BeTrue())

				subject()

				file, err := os.Open("./tmp/system-configurator/config.yaml")
				Expect(err).ToNot(HaveOccurred())
				defer file.Close()

				buf := make([]byte, 2)
				file.Read(buf)
				Expect(buf).To(Equal([]byte("{}")))
			})
		})

		Context("when no configuration is given", func() {
			var original func() string
			BeforeEach(func() {
				cfgFixture = ""
				cfg = nil
				original = store.LocalDefaultLocation
				store.LocalDefaultLocation = func() string {
					return "./default-tmp/system-configurator"
				}
			})

			AfterEach(func() {
				os.RemoveAll("./default-tmp/system-configurator")
				store.LocalDefaultLocation = original
			})

			It("loads the config from the default location", func() {
				subject()
				_, err := os.Open(path.Join("./default-tmp/system-configurator", store.LocalDefaultFileName))
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
			}, cmpopts.IgnoreUnexported(model.Package{})))
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

	Describe("LocalDefaultLocation", func() {
		subject := func() string {
			return store.LocalDefaultLocation()
		}

		It("returns the default location for the local config file", func() {
			Expect(subject()).To(Equal(path.Join(os.Getenv("HOME"), ".config/system-configurator")))
		})

		Context("when the $HOME environment variable is not set", func() {
			BeforeEach(func() {
				os.Unsetenv("HOME")
			})

			It("panics", func() {
				Expect(func() { subject() }).To(Panic())
			})
		})
	})
})
