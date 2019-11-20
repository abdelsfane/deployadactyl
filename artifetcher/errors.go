package artifetcher

import "fmt"

type CreateTempFileError struct {
	Err error
}

func (e CreateTempFileError) Error() string {
	return fmt.Sprintf("cannot create temp file: %s", e.Err)
}

type FetcherRequestError struct {
	Err error
}

func (e FetcherRequestError) Error() string {
	return fmt.Sprintf("cannot create artifact fetch request: %s", e.Err)
}

type GetUrlError struct {
	Url string
	Err error
}

func (e GetUrlError) Error() string {
	return fmt.Sprintf("cannot GET url: %s: %s", e.Url, e.Err)
}

type GetStatusError struct {
	Url    string
	Status string
}

func (e GetStatusError) Error() string {
	return fmt.Sprintf("cannot GET url: %s: %s", e.Url, e.Status)
}

type ArtifactoryTimeoutError struct {
	Url    string
	Status string
}

func (e ArtifactoryTimeoutError) Error() string {
	return fmt.Sprintf(`*******************

The following error was found in the above logs:

Error: Download of application artifact timed out

Potential Solution: Reduce the size of the artifact

*******************`)
}

type WriteResponseError struct {
	Err error
}

func (e WriteResponseError) Error() string {
	return fmt.Sprintf("cannot write response to file: %s", e.Err)
}

type CreateTempDirectoryError struct {
	Err error
}

func (e CreateTempDirectoryError) Error() string {
	return fmt.Sprintf("cannot create temp directory: %s", e.Err)
}

type NonProcessError struct {
	Err error
}

func (e NonProcessError) Error() string {
	return fmt.Sprintf("cannot process artifact: %s", e.Err)
}

type UnsupportedFormatError struct{}

func (e UnsupportedFormatError) Error() string {
	return fmt.Sprintf("File format not supported")
}
