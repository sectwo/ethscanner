package env

import (
	"github.com/naoina/toml"
	"os"
)

type Env struct {
	DB struct {
		Uri string
		DB  string

		Block string
		Tx    string
	}

	Node struct {
		Dial       string
		StartBlock uint64
		EndBlock   uint64
	}
	//Log struct {
	//}
}

func NewEnv(path string) *Env {
	env := new(Env)

	if file, err := os.Open(path); err != nil {
		panic(err)
	} else if err = toml.NewDecoder(file).Decode(env); err != nil {
		panic(err)
	} else {
		return env
	}
}
