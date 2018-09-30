package storage

import (
	"context"
	"io"
	"time"
)

type Store interface {
	List(ctx context.Context, path string, opts ListOpts) (Blobs, error)
	Stat(ctx context.Context, name string) (*Stat, error)
	New(ctx context.Context, name string) Blob
	Put(ctx context.Context, name string, out Blob) error
	Get(ctx context.Context, name string, in Blob) error
	Delete(ctx context.Context, name string) error
	Move(ctx context.Context, name string, to string) error
}

type SyncOpts struct {
}

type SyncListOpts struct {
}

type ListOpts struct {
	Recursive bool
}

type Blob interface {
	Name() string
	Path() string
	IsDir() bool
	Exists() bool
	Stat() (*Stat, error)

	io.ReadWriteCloser
}

type Blobs map[string]Blob

type Stat struct {
	LastModified time.Time
	ETag         string
	ContentType  string
	Size         int64
}
