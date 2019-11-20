// Package extractor unzips artifacts.
package extractor

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path"

	"archive/tar"
	"compress/gzip"
	I "github.com/compozed/deployadactyl/interfaces"
	"github.com/spf13/afero"
)

type ExtractorConstructor func(log I.DeploymentLogger, fs *afero.Afero) I.Extractor

func NewExtractor(log I.DeploymentLogger, fs *afero.Afero) I.Extractor {
	return &Extractor{
		Log:        log,
		FileSystem: fs,
	}
}

// Extractor has a file system from which files are extracted from.
type Extractor struct {
	Log        I.DeploymentLogger
	FileSystem *afero.Afero
}

// Unzip unzips from source into destination.
// If there is no manifest provided to this function, it will attempt to read a manifest file within the zip file.
func (e *Extractor) Unzip(source, destination, manifest string) error {
	e.Log.Info("extracting application")
	e.Log.Debugf(`parameters for extractor:
	source: %+v
	destination: %+v`, source, destination)

	err := e.FileSystem.MkdirAll(destination, 0755)
	if err != nil {
		return CreateDirectoryError{err}
	}

	file, err := e.FileSystem.Open(source)
	if err != nil {
		return err
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		return err
	}

	//reader, err := zip.NewReader(file, fileStat.Size())
	reader, err := zip.NewReader(file, fileStat.Size())
	if err != nil {
		return OpenZipError{source, err}
	}

	for _, file := range reader.File {
		err := e.unzipFile(destination, file)
		if err != nil {
			return ExtractFileError{file.Name, err}
		}
	}

	if manifest != "" {
		manifestFile, err := e.FileSystem.OpenFile(path.Join(destination, "manifest.yml"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return OpenManifestError{err}
		}
		defer manifestFile.Close()

		_, err = fmt.Fprint(manifestFile, manifest)
		if err != nil {
			return PrintToManifestError{err}
		}
	}

	e.Log.Info("extract was successful")
	return nil
}

func (e *Extractor) unzipFile(destination string, file *zip.File) error {
	contents, err := file.Open()
	if err != nil {
		return ExtractFileError{file.Name, err}
	}
	defer contents.Close()

	if file.FileInfo().IsDir() {
		return nil
	}

	savedLocation := path.Join(destination, file.Name)
	directory := path.Dir(savedLocation)
	err = e.FileSystem.MkdirAll(directory, 0755)
	if err != nil {
		return MakeDirectoryError{directory, err}
	}

	mode := file.Mode()
	newFile, err := e.FileSystem.OpenFile(savedLocation, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return OpenFileError{savedLocation, err}
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, contents)
	if err != nil {
		return WriteFileError{savedLocation, err}
	}

	return nil
}

func (e *Extractor) Untar(source, destination, manifest string) error {
	e.Log.Info("extracting application")
	e.Log.Debugf(`parameters for extractor:
	source: %+v
	destination: %+v`, source, destination)

	f, err := e.FileSystem.Open(source)
	if err != nil {
		return err
	}
	defer f.Close()

	uncompressedStream, err := gzip.NewReader(f)
	if err != nil {
		return ReadTarError{err}
	}

	tarReader := tar.NewReader(uncompressedStream)

	for true {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Printf("ExtractTarGz: Next() failed: %s", err.Error())
		}

		switch header.Typeflag {
		case tar.TypeReg:
			savedLocation := path.Join(destination, header.Name)
			directory := path.Dir(savedLocation)
			err = e.FileSystem.MkdirAll(directory, 0755)

			newFile, err := e.FileSystem.OpenFile(savedLocation, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)

			if err != nil {
				fmt.Printf("ExtractTarGz: Create() failed: %s", err.Error())
			}
			defer newFile.Close()
			if _, err := io.Copy(newFile, tarReader); err != nil {
				fmt.Printf("ExtractTarGz: Copy() failed: %s", err.Error())
			}
		}
	}

	if manifest != "" {
		manifestFile, err := e.FileSystem.OpenFile(path.Join(destination, "manifest.yml"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return OpenManifestError{err}
		}
		defer manifestFile.Close()

		_, err = fmt.Fprint(manifestFile, manifest)
		if err != nil {
			return PrintToManifestError{err}
		}
	}

	e.Log.Info("extract was successful")

	return nil
}
