package token

import (
	"github.com/vite-cloud/vite/core/domain/datadir"
	"time"
)

const Store = datadir.Store("tokens")

type Token struct {
	Label      string
	Value      string
	CreatedAt  time.Time
	LastUsedAt time.Time
}

func (t Token) Time() time.Time {
	return t.CreatedAt
}

func (t Token) ID() string {
	return t.Label
}
