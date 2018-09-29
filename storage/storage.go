package storage

import (
	"context"
	"io"
	"time"
)

type Store interface {
	List(ctx context.Context, path string, opts ListOpts) (Blobs, error)
	Stat(ctx context.Context, name string) (*Stat, error)
	Put(ctx context.Context, name string, data io.ReadCloser) error
	Get(ctx context.Context, name string) (io.Reader, error)
	Sync(ctx context.Context, from Blob, to string, opts SyncOpts) error
	SyncList(ctx context.Context, list Blobs, to string, opts SyncListOpts) error
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
	Stat() (*Stat, error)
	Reader() (io.Reader, error)
}

type Blobs map[string]Blob

type Stat struct {
	LastModified time.Time
	ETag         string
	ContentType  string
	Size         int64
}
