package storage

import (
	"context"
	"io"
	"io/ioutil"
	"path"

	"github.com/pkg/errors"
)

type Filesystem struct {
	root string
}

type File struct {
	name  string
	path  string
	isDir bool
}

func (f *File) Name() string {
	return f.name
}

func (f *File) Path() string {
	return f.path
}

func (f *File) IsDir() bool {
	return f.isDir
}

func (f *File) Stat() (*Stat, error) {
	return nil, nil
}

func (f *File) Reader() (io.Reader, error) {
	return nil, nil
}

func NewFilesystem(root string) (Store, error) {
	store := &Filesystem{
		root: root,
	}

	return store, nil
}

func (store *Filesystem) List(ctx context.Context, prefix string, opts ListOpts) (Blobs, error) {
	fullPath := path.Join(store.root, prefix)
	fList, err := ioutil.ReadDir(fullPath)
	if err != nil {
		return nil, errors.Wrap(err, "cannot read dir")
	}

	var blobs = make(Blobs)
	for _, f := range fList {
		file := &File{
			path:  fullPath,
			name:  f.Name(),
			isDir: f.IsDir(),
		}

		blobs[f.Name()] = file
	}

	return blobs, nil
}

func (store *Filesystem) Stat(ctx context.Context, name string) (*Stat, error) {
	return nil, nil
}

func (store *Filesystem) Put(ctx context.Context, name string, data io.ReadCloser) error {
	return nil
}

func (store *Filesystem) Get(ctx context.Context, name string) (io.Reader, error) {
	return nil, nil
}

func (store *Filesystem) Sync(ctx context.Context, from Blob, to string, opts SyncOpts) error {
	return nil
}

func (store *Filesystem) SyncList(ctx context.Context, list Blobs, to string, opts SyncListOpts) error {
	return nil
}

func (store *Filesystem) Delete(ctx context.Context, name string) error {
	return nil
}

func (store *Filesystem) Move(ctx context.Context, name string, to string) error {
	return nil
}
