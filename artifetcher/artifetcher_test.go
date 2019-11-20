package artifetcher_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"

	"github.com/op/go-logging"

	. "github.com/compozed/deployadactyl/artifetcher"
	E "github.com/compozed/deployadactyl/artifetcher/extractor"
	"github.com/compozed/deployadactyl/interfaces"
	"github.com/compozed/deployadactyl/mocks"
	"github.com/compozed/deployadactyl/randomizer"
)

var _ = Describe("Artifetcher", func() {
	var (
		artifetcher *Artifetcher
		af          *afero.Afero
		extractor   *mocks.Extractor
		testserver  *httptest.Server
		manifest    string
		log         interfaces.DeploymentLogger
	)

	BeforeEach(func() {
		log = interfaces.DeploymentLogger{Log: interfaces.DefaultLogger(GinkgoWriter, logging.DEBUG, "artifetcher_test")}
		af = &afero.Afero{Fs: afero.NewMemMapFs()}
		extractor = &mocks.Extractor{}
		artifetcher = &Artifetcher{af, extractor, log}
		manifest = "manifest-" + randomizer.StringRunes(10)

		testserver = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.RequestURI == "/timeoutTest" {
				w.WriteHeader(504)
			} else if r.RequestURI == "/tarFile" {
				http.ServeFile(w, r, "./fixtures/deployadactyl-fixture.tar")
			} else if r.RequestURI == "/rarFile" {
				http.ServeFile(w, r, "./fixtures/deployadactyl-fixture.rar")
			} else {
				http.ServeFile(w, r, "./fixtures/deployadactyl-fixture.jar")
			}

		}))
	})

	AfterEach(func() {
		testserver.Close()
	})

	Context("when a zip file is provided", func() {
		Describe("fetching a zip file", func() {
			It("can fetch a jar file", func() {
				extractor.UnzipCall.Returns.Error = nil

				unzippedPath, err := artifetcher.Fetch(testserver.URL, "")
				Expect(err).ToNot(HaveOccurred())

				Expect(af.IsDir(unzippedPath)).To(BeTrue())

				Expect(extractor.UnzipCall.Received.Source).To(ContainSubstring("deployadactyl-artifact"))
				Expect(extractor.UnzipCall.Received.Destination).To(Equal(unzippedPath))
				Expect(extractor.UnzipCall.Received.Manifest).To(BeEmpty())
			})

			It("returns an error when an invalid url is given", func() {
				_, err := artifetcher.Fetch("example://example.example", manifest)
				Expect(err).To(HaveOccurred())
			})

			It("returns an error when the URL returns a 404 not found", func() {
				testserver = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					http.Error(w, "not found", 404)
				}))

				_, err := artifetcher.Fetch(testserver.URL, manifest)
				Expect(err).To(HaveOccurred())
			})

			Context("when extractor fails", func() {
				It("returns an error", func() {
					extractor.UnzipCall.Returns.Error = errors.New("unzip call failed")

					_, err := artifetcher.Fetch(testserver.URL, "")

					Expect(err).To(MatchError(NonProcessError{errors.New("unzip call failed")}))
				})
			})

			Context("when request to retrieve artifact times out", func() {
				It("should return an error", func() {
					_, err := artifetcher.Fetch(testserver.URL+"/timeoutTest", "")
					Expect(err).To(HaveOccurred())

					Expect(err.Error()).To(ContainSubstring("Download of application artifact timed out"))
				})
			})
		})

		Describe("fetching a zip file from a request", func() {
			It("returns the path to the unzipped directory and manifest", func() {
				artifetcher = &Artifetcher{af, E.NewExtractor(log, af), log}

				expectManifest := `---
applications:
- name: artifact-with-manifest
  memory: 512M`

				body, err := os.Open("./fixtures/artifact-with-manifest.jar")
				Expect(err).ToNot(HaveOccurred())

				path, manifest, err := artifetcher.FetchArtifactFromRequest(body, "application/zip")
				Expect(err).ToNot(HaveOccurred())

				Expect(path).To(ContainSubstring("deployadactyl-"))
				Expect(manifest).To(ContainSubstring(expectManifest))
			})

			Context("when extractor fails", func() {
				It("returns an error", func() {
					errorMessage := "test extract fail"
					extractor.UnzipCall.Returns.Error = errors.New(errorMessage)

					body, err := os.Open("./fixtures/artifact-with-manifest.jar")
					Expect(err).ToNot(HaveOccurred())

					path, _, err := artifetcher.FetchArtifactFromRequest(body, "application/zip")
					Expect(err).To(MatchError(NonProcessError{errors.New(errorMessage)}))

					Expect(path).To(BeEmpty())
				})
			})
		})
	})

	Context("when a tar file is provided", func() {
		Describe("fetching a tar file", func() {
			It("can fetch a tar file", func() {
				extractor.UntarCall.Returns.Error = nil

				untarPath, err := artifetcher.Fetch(testserver.URL+"/tarFile", "")
				Expect(err).ToNot(HaveOccurred())

				Expect(af.IsDir(untarPath)).To(BeTrue())

				Expect(extractor.UntarCall.Received.Source).To(ContainSubstring("deployadactyl-artifact"))
				Expect(extractor.UntarCall.Received.Destination).To(Equal(untarPath))
				Expect(extractor.UntarCall.Received.Manifest).To(BeEmpty())
			})
		})

		Describe("fetching tar file from a request", func() {
			It("should return the path to the unarchived directory and manifest", func() {
				artifetcher = &Artifetcher{af, E.NewExtractor(log, af), log}

				expectManifest := `---
applications:
- name: deployadactyl
  memory: 256M
  disk_quota: 256M`

				body, err := os.Open("./fixtures/deployadactyl-fixture.tar")
				Expect(err).ToNot(HaveOccurred())

				path, manifest, err := artifetcher.FetchArtifactFromRequest(body, "application/x-tar")
				Expect(err).ToNot(HaveOccurred())

				Expect(path).To(ContainSubstring("deployadactyl-"))
				Expect(manifest).To(ContainSubstring(expectManifest))
			})
		})
	})

	Context("when neither tar nor zip format is provided", func() {
		Context("when a request with a json body is sent", func() {
			It("should return invalid format error", func() {
				_, err := artifetcher.Fetch(testserver.URL+"/rarFile", "")
				Expect(err).To(HaveOccurred())

				Expect(err.Error()).To(ContainSubstring("File format not supported"))
				Expect(err).To(MatchError(UnsupportedFormatError{}))
			})
		})

		Context("when a request with the artifact attached is sent", func() {
			It("should return invalid format error", func() {
				body, err := os.Open("./fixtures/deployadactyl-fixture.rar")
				Expect(err).ToNot(HaveOccurred())

				_, _, err = artifetcher.FetchArtifactFromRequest(body, "application/not-a-tar-or-a-zip")
				Expect(err).To(HaveOccurred())

				Expect(err).To(MatchError(UnsupportedFormatError{}))
			})
		})
	})
})
