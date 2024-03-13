package repo

import (
	"scanner/env"
	"scanner/repo/db"
	"scanner/repo/node"
)

type Repo struct {
	env *env.Env

	DB   db.DBImpl
	Node node.NodeImpl
}

func NewRepo(env *env.Env) (*Repo, error) {
	r := &Repo{
		env: env,
	}

	var err error

	if r.DB, err = db.NewDB(env); err != nil {
		panic(err)
	} else if r.Node, err = node.NewNode(env); err != nil {
		panic(err)
	} else {
		return r, nil
	}
}
