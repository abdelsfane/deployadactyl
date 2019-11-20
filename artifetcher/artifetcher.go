// Package artifetcher downloads the artifact given a URL.
package artifetcher

import (
	"io"
	"net"
	"net/http"
	"time"

	I "github.com/compozed/deployadactyl/interfaces"
	"github.com/spf13/afero"
)

type ArtifetcherConstructor func(fs *afero.Afero, ex I.Extractor, log I.DeploymentLogger) I.Fetcher

func NewArtifetcher(fs *afero.Afero, ex I.Extractor, log I.DeploymentLogger) I.Fetcher {
	return &Artifetcher{
		FileSystem: fs,
		Extractor:  ex,
		Log:        log,
	}
}

// Artifetcher fetches artifacts within a file system with an Extractor.
type Artifetcher struct {
	FileSystem *afero.Afero
	Extractor  I.Extractor
	Log        I.DeploymentLogger
}

// Fetch downloads an artifact located at URL.
// It then passes it to the extractor with the manifest for unzipping.
//
// Returns a string to the unzipped artifacts path and an error.
func (a *Artifetcher) Fetch(url, manifest string) (string, error) {
	a.Log.Info("fetching artifact")
	a.Log.Debugf("artifact URL: %s", url)

	artifactFile, err := a.FileSystem.TempFile("", "deployadactyl-artifact-")
	if err != nil {
		return "", CreateTempFileError{err}
	}
	defer artifactFile.Close()
	defer a.FileSystem.Remove(artifactFile.Name())

	var client = &http.Client{
		Timeout: 15 * time.Minute,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   60 * time.Second,
				KeepAlive: 60 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   15 * time.Second,
			ResponseHeaderTimeout: 15 * time.Second,
			ExpectContinueTimeout: 2 * time.Second,
		},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", FetcherRequestError{err}
	}

	response, err := client.Do(req)
	if err != nil {
		return "", GetUrlError{url, err}
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == 504 {
			return "", ArtifactoryTimeoutError{url, response.Status}
		} else {
			return "", GetStatusError{url, response.Status}
		}
	}

	_, err = io.Copy(artifactFile, response.Body)
	if err != nil {
		return "", WriteResponseError{err}
	}

	unarchivedPath, err := a.FileSystem.TempDir("", "deployadactyl-unarchived-")
	if err != nil {
		return "", CreateTempDirectoryError{err}
	}

	if response.Header.Get("Content-Type") == "application/java-archive" || response.Header.Get("Content-Type") == "application/zip" {
		err = a.Extractor.Unzip(artifactFile.Name(), unarchivedPath, manifest)
		if err != nil {
			a.FileSystem.RemoveAll(unarchivedPath)
			return "", NonProcessError{err}

		}
	} else if response.Header.Get("Content-Type") == "application/x-tar" || response.Header.Get("Content-Type") == "application/x-gzip" {
		err = a.Extractor.Untar(artifactFile.Name(), unarchivedPath, manifest)
		if err != nil {
			a.FileSystem.RemoveAll(unarchivedPath)
			return "", NonProcessError{err}

		}
	} else {
		return "", UnsupportedFormatError{}
	}

	a.Log.Debugf("fetched and unarchived to tempdir: %s", unarchivedPath)
	return unarchivedPath, nil
}

// FetchZipFromRequest fetches files from a compressed zip file in the request body.
//
// Returns a string to the unzipped application path and an error.
func (a *Artifetcher) FetchArtifactFromRequest(body io.Reader, contentType string) (string, string, error) {

	file, err := a.FileSystem.TempFile("", "deployadactyl-")
	if err != nil {
		return "", "", CreateTempFileError{err}
	}
	defer file.Close()
	defer a.FileSystem.Remove(file.Name())

	a.Log.Infof("fetching file %s", file.Name())
	_, err = io.Copy(file, body)
	if err != nil {
		return "", "", WriteResponseError{err}
	}

	unarchivedPath, err := a.FileSystem.TempDir("", "deployadactyl-")
	if err != nil {
		return "", "", CreateTempDirectoryError{err}
	}

	if contentType == "application/zip" || contentType == "application/java-archive" {
		err = a.Extractor.Unzip(file.Name(), unarchivedPath, "")
		if err != nil {
			a.FileSystem.RemoveAll(unarchivedPath)
			return "", "", NonProcessError{err}
		}
	} else if contentType == "application/x-tar" || contentType == "application/x-gzip" {
		err = a.Extractor.Untar(file.Name(), unarchivedPath, "")
	} else {
		return "", "", UnsupportedFormatError{}
	}

	manifest, err := a.FileSystem.ReadFile(unarchivedPath + "/manifest.yml")
	if err != nil {
		return "", "", err
	}

	a.Log.Debugf("fetched and unarchived to tempdir %s", unarchivedPath)
	return unarchivedPath, string(manifest), nil
}
