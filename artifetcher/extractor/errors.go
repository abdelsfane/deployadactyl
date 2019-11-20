package extractor

import "fmt"

type CreateDirectoryError struct {
	Err error
}

func (e CreateDirectoryError) Error() string {
	return fmt.Sprintf("cannot create directory: %s", e.Err)
}

type OpenZipError struct {
	Source string
	Err    error
}

func (e OpenZipError) Error() string {
	niceFixYourZipMessage := `Please double check your zip compression method and that the correct files are zipped.
You can try confirming that it's valid on your computer by opening or performing some other action on it.
Once you've confirmed that it's valid, please try again.`

	return fmt.Sprintf("cannot open zip file: %s: %s\n%s", e.Source, e.Err, niceFixYourZipMessage)
}

type ExtractFileError struct {
	FileName string
	Err      error
}

func (e ExtractFileError) Error() string {
	return fmt.Sprintf("cannot extract file from archive: %s: %s", e.FileName, e.Err)
}

type OpenManifestError struct {
	Err error
}

func (e OpenManifestError) Error() string {
	return fmt.Sprintf("cannot open manifest file: %s", e.Err)
}

type PrintToManifestError struct {
	Err error
}

func (e PrintToManifestError) Error() string {
	return fmt.Sprintf("cannot print to open manifest file: %s", e.Err)
}

type MakeDirectoryError struct {
	Directory string
	Err       error
}

func (e MakeDirectoryError) Error() string {
	return fmt.Sprintf("cannot make directory: %s: %s", e.Directory, e.Err)
}

type OpenFileError struct {
	SavedLocation string
	Err           error
}

func (e OpenFileError) Error() string {
	return fmt.Sprintf("cannot open file for writing: %s: %s", e.SavedLocation, e.Err)
}

type WriteFileError struct {
	SavedLocation string
	Err           error
}

func (e WriteFileError) Error() string {
	return fmt.Sprintf("cannot write to file: %s: %s", e.SavedLocation, e.Err)
}

type ReadTarError struct {
	Err error
}

func (e ReadTarError) Error() string {
	return fmt.Sprintf("Failed to untar file: %s", e.Err)
}
