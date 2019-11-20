package extractor_test

import (
	"io/ioutil"
	"path"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"

	"github.com/op/go-logging"

	. "github.com/compozed/deployadactyl/artifetcher/extractor"
	"github.com/compozed/deployadactyl/interfaces"
	"github.com/compozed/deployadactyl/randomizer"
)

const deployadactylManifest = `---
applications:
- name: deployadactyl
  memory: 256M
  disk_quota: 256M
`

var _ = Describe("Extracting", func() {
	var (
		af          *afero.Afero
		file        string
		tarFile     string
		destination string
		extractor   Extractor
	)

	BeforeEach(func() {
		file = "/artifact.jar"
		tarFile = "/artifact.tar"
		destination = "../fixtures/deployadactyl-fixture"
		af = &afero.Afero{Fs: afero.NewMemMapFs()}
		extractor = Extractor{interfaces.DeploymentLogger{Log: interfaces.DefaultLogger(GinkgoWriter, logging.DEBUG, "extractor_test")}, af}

		fileBytes, err := ioutil.ReadFile("../fixtures/deployadactyl-fixture.jar")
		tarFileBytes, err := ioutil.ReadFile("../fixtures/deployadactyl-fixture.tar")
		Expect(err).ToNot(HaveOccurred())

		Expect(af.WriteFile(file, fileBytes, 0644)).To(Succeed())
		Expect(af.WriteFile(tarFile, tarFileBytes, 0644)).To(Succeed())
	})

	AfterEach(func() {
		Expect(af.RemoveAll(destination)).To(Succeed())
	})

	Context("when zip and jar formats are used", func() {
		It("unzips the artifact", func() {
			Expect(extractor.Unzip(file, destination, "")).To(Succeed())

			extractedFile, err := af.ReadFile(path.Join(destination, "index.html"))
			Expect(err).ToNot(HaveOccurred())

			Expect(extractedFile).To(ContainSubstring("public/assets/images/pterodactyl.png"))
		})

		Context("when manifest is an empty string", func() {
			It("leaves the manifest alone", func() {
				Expect(extractor.Unzip(file, destination, "")).To(Succeed())

				extractedManifest, err := af.ReadFile(path.Join(destination, "manifest.yml"))
				Expect(err).ToNot(HaveOccurred())

				Expect(extractedManifest).To(BeEquivalentTo(deployadactylManifest))
			})
		})

		Context("when manifest is not an empty string", func() {
			It("unzips the artifact and overwrites the manifest", func() {
				manifestContents := "manifestContents-" + randomizer.StringRunes(10)
				Expect(extractor.Unzip(file, destination, manifestContents)).To(Succeed())

				extractedManifest, err := af.ReadFile(path.Join(destination, "manifest.yml"))
				Expect(err).ToNot(HaveOccurred())

				Expect(extractedManifest).To(BeEquivalentTo(manifestContents))
			})
		})

		It("can not unzip an invalid file", func() {
			file := "../fixtures/bad-deployadactyl-fixture.tgz"
			destination = "../fixtures/bad-deployadactyl-fixture"
			af = &afero.Afero{Fs: afero.NewMemMapFs()}

			extractor := Extractor{interfaces.DeploymentLogger{Log: interfaces.DefaultLogger(GinkgoWriter, logging.DEBUG, "extractor_test")}, af}

			Expect(extractor.Unzip(file, destination, "")).ToNot(Succeed())
		})
	})

	Context("when tar or gz formats are used", func() {
		It("should unarchive the provided tar file", func() {
			err := extractor.Untar(tarFile, destination, "")
			Expect(err).ToNot(HaveOccurred())

			extractedFile, err := af.ReadFile(path.Join(destination, "index.html"))
			Expect(extractedFile).To(ContainSubstring("public/assets/images/pterodactyl.png"))
		})

		Context("when manifest is an empty string", func() {
			It("leaves manifest alone", func() {
				err := extractor.Untar(tarFile, destination, "")
				Expect(err).ToNot(HaveOccurred())

				extractedManifest, err := af.ReadFile(path.Join(destination, "manifest.yml"))
				Expect(err).ToNot(HaveOccurred())

				Expect(extractedManifest).To(BeEquivalentTo(deployadactylManifest))
			})
		})

		Context("when manifest is not an empty string", func() {
			It("untars the files and copies manifest string to directory", func() {
				manifestContents := "manifestContents-" + randomizer.StringRunes(10)

				err := extractor.Untar(tarFile, destination, manifestContents)
				Expect(err).ToNot(HaveOccurred())

				extractedManifest, err := af.ReadFile(path.Join(destination, "manifest.yml"))
				Expect(err).ToNot(HaveOccurred())

				Expect(extractedManifest).To(BeEquivalentTo(manifestContents))
			})
		})

		Context("when it is unable to untar the file", func() {
			It("returns an error", func() {
				err := extractor.Untar(file, destination, "")
				Expect(err).To(HaveOccurred())

				Expect(err.Error()).To(ContainSubstring("Failed to untar file"))
			})
		})
	})

})
