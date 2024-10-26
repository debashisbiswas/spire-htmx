package entry

import (
	"github.com/ryanskidmore/libsql-vector-go"
	"time"
)

type Entry struct {
	Time      time.Time
	Content   string
	Embedding libsqlvector.Vector
}
