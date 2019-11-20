package interfaces

import (
	"io"

	S "github.com/compozed/deployadactyl/structs"
)

type BlueGreener interface {
	Execute(
		actionCreator ActionCreator,
		environment S.Environment,
		response io.ReadWriter,
	) error
}
