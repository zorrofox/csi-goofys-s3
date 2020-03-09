package s3_test

import (
	"log"
	"os"

    "github.com/zorrofox/csi-goofys-s3/pkg/s3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/kubernetes-csi/csi-test/pkg/sanity"
)

var _ = Describe("S3Driver", func() {

	Context("goofys", func() {
		socket := "/tmp/csi-goofys.sock"
		csiEndpoint := "unix://" + socket
		if err := os.Remove(socket); err != nil && !os.IsNotExist(err) {
			Expect(err).NotTo(HaveOccurred())
		}
		driver, err := s3.NewS3("test-node", csiEndpoint)
		if err != nil {
			log.Fatal(err)
		}
		go driver.Run()

		Describe("CSI sanity", func() {
			sanityCfg := &sanity.Config{
				TargetPath:  os.TempDir() + "/goofys-target",
				StagingPath: os.TempDir() + "/goofys-staging",
				Address:     csiEndpoint,
				//SecretsFile: "../../test/secret.yaml",
				TestVolumeParameters: map[string]string{
					"mounter": "goofys",
				},
				IDGen: &sanity.DefaultIDGenerator{}, // Bug for ID Generator Go Panic
			}
			sanity.GinkgoTest(sanityCfg)
		})
	})

})
