package datagears_test

import (
	"github.com/jsam/dg/datagears"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
)

var _ = Describe("Manifest", func() {

	manifestExample := `version: 1
remotes:
  myredis:
    host: localhost
    port: 6379
    database: 0
gears:
  mygear1:
    entrypoint: ./gear.py
    values:
      databaseHost: localhost
      databasePort: 6543
    secrets:
      - DATABASE_PASSWORD
    requirements:
      - numpy==0.19.2
`

	Context("check for manifest", func() {
		It("check expected manifest construction", func() {
			err := ioutil.WriteFile("./datagears.yml", []byte(manifestExample), 0644)
			Expect(err).ShouldNot(HaveOccurred())

			manifest := datagears.NewDGManifest("")
			Expect(manifest).NotTo(BeNil())

			Expect(manifest.Version).To(Equal(1))

			Expect(len(manifest.Remotes)).To(Equal(1))
			Expect(len(manifest.Gears)).To(Equal(1))

			remote, ok := manifest.Remotes["myredis"]
			Expect(ok).To(BeTrue())
			Expect(remote.Host).To(Equal("localhost"))
			Expect(remote.Port).To(Equal(6379))
			Expect(remote.Database).To(Equal(0))

			gear, ok := manifest.Gears["mygear1"]
			Expect(ok).To(BeTrue())
			Expect(gear.Entrypoint).To(Equal("./gear.py"))

			expectedValues := make(map[string]string)
			expectedValues["databaseHost"] = "localhost"
			expectedValues["databasePort"] = "6543"
			Expect(gear.Values).To(Equal(expectedValues))

			Expect(gear.Secrets).To(Equal([]string{
				"DATABASE_PASSWORD",
			}))

			Expect(gear.Requirements).To(Equal([]string{
				"numpy==0.19.2",
			}))

			_ = os.Remove("./datagears.yml")
		})
	})
})
