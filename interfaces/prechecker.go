package interfaces

import S "github.com/compozed/deployadactyl/structs"

// Prechecker interface.
type Prechecker interface {
	AssertAllFoundationsUp(environment S.Environment) error
}
