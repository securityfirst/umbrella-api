package main

import (
	"github.com/go-gorp/gorp"
)

type Umbrella struct {
	Db *gorp.DbMap
}
