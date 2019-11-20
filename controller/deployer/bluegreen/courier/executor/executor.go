// Package executor runs commands against the Cloud Foundry binary.
package executor

import (
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/afero"
)

// New returns a new Executor struct.
func New(fileSystem *afero.Afero) (Executor, error) {
	tempDir, err := fileSystem.TempDir("", "deployadactyl-executor-")
	if err != nil {
		return Executor{}, err
	}

	return Executor{
		fileSystem: fileSystem,
		tempDir:    tempDir,
	}, nil
}

// Executor has a file system that is used to execute the Cloud Foundry CLI.
type Executor struct {
	tempDir    string
	fileSystem *afero.Afero
}

// Execute takes a slice of string args and runs them together against the cf command on the Cloud Foundry binary.
//
// Returns the combined standard output and standard error.
func (e Executor) Execute(args ...string) ([]byte, error) {
	command := exec.Command("cf", args...)
	command.Env = setEnv(os.Environ(), "CF_HOME", e.tempDir)
	return command.CombinedOutput()
}

// ExecuteInDirectory does the same thing as Execute does, but does it in a specific directory.
//
// Returns the combined standard output and standard error.
func (e Executor) ExecuteInDirectory(directory string, args ...string) ([]byte, error) {
	command := exec.Command("cf", args...)
	command.Env = setEnv(os.Environ(), "CF_HOME", e.tempDir)
	command.Dir = directory
	return command.CombinedOutput()
}

// CleanUp removes the temporary directory of the Executor.
func (e Executor) CleanUp() error {
	return e.fileSystem.RemoveAll(e.tempDir)
}

func setEnv(env []string, key, value string) []string {
	keyValuePair := key + "=" + value

	for i, envVar := range env {
		if strings.HasPrefix(envVar, key+"=") {
			env[i] = keyValuePair
			return env
		}
	}

	return append(env, keyValuePair)
}
