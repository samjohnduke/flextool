package storage

import (
	"context"
	"io"
)

type DOSpace struct {
	accessKey string
	secretKey string
	endpoint  string
	space     string
	secure    bool
}

func NewDOSpace(accessKey, secretKey, endpoint, space string, secure bool) (Store, error) {
	store := &DOSpace{
		accessKey: accessKey,
		secretKey: secretKey,
		endpoint:  endpoint,
		space:     space,
		secure:    secure,
	}

	return store, nil
}

func (dos *DOSpace) List(ctx context.Context, path string, opts ListOpts) (Blobs, error) {
	return nil, nil
}

func (dos *DOSpace) Stat(ctx context.Context, name string) (*Stat, error) {
	return nil, nil
}

func (dos *DOSpace) Put(ctx context.Context, name string, data io.ReadCloser) error {
	return nil
}

func (dos *DOSpace) Get(ctx context.Context, name string) (io.Reader, error) {
	return nil, nil
}

func (dos *DOSpace) Sync(ctx context.Context, from Blob, to string, opts SyncOpts) error {
	return nil
}

func (dos *DOSpace) SyncList(ctx context.Context, list Blobs, to string, opts SyncListOpts) error {
	return nil
}

func (dos *DOSpace) Delete(ctx context.Context, name string) error {
	return nil
}

func (dos *DOSpace) Move(ctx context.Context, name string, to string) error {
	return nil
}
