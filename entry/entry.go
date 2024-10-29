package entry

import "time"

type Vector []float32

type Entry struct {
	Time      time.Time
	Content   string
	Embedding Vector
}
