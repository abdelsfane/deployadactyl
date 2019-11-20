package interfaces

import "github.com/compozed/deployadactyl/structs"

type EnvResolver interface {
	Resolve(env string) (structs.Environment, error)
}
