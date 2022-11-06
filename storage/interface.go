package storage

import "io"

type CloudStorage interface {
	Put(path string, reader *io.ReadSeeker) error
	Get(path string, writer *io.WriteSeeker) error
}
