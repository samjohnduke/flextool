package storage

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

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
	stat := &Stat{}
	fd, err := os.Open(path.Join(f.path, f.name))
	if err != nil {
		return nil, err
	}

	fstat, err := fd.Stat()
	if err != nil {
		return nil, err
	}

	f.isDir = fstat.IsDir()
	stat.Size = fstat.Size()
	stat.LastModified = fstat.ModTime()
	stat.ContentType = "application/octet-stream"

	h := md5.New()
	if _, err := io.Copy(h, fd); err != nil {
		return nil, err
	}
	stat.ETag = fmt.Sprintf("%x", h.Sum(nil))

	return stat, nil
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

	if opts.Recursive {
		return store.listTree(ctx, fullPath, opts)
	}

	return store.listFlat(ctx, fullPath, opts)
}

func (store *Filesystem) listFlat(ctx context.Context, prefix string, opts ListOpts) (Blobs, error) {
	fList, err := ioutil.ReadDir(prefix)
	if err != nil {
		return nil, errors.Wrap(err, "cannot read dir")
	}

	var blobs = make(Blobs)
	for _, f := range fList {
		file := &File{
			path:  prefix,
			name:  f.Name(),
			isDir: f.IsDir(),
		}

		blobs[f.Name()] = file
	}

	return blobs, nil
}

func (store *Filesystem) listTree(ctx context.Context, prefix string, opts ListOpts) (Blobs, error) {
	var blobs = make(Blobs)

	err := filepath.Walk(prefix, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		file := &File{
			path:  path.Join(prefix, p),
			name:  info.Name(),
			isDir: info.IsDir(),
		}

		blobs[info.Name()] = file

		return nil
	})
	if err != nil {
		return nil, err
	}

	return blobs, nil
}

func (store *Filesystem) Stat(ctx context.Context, name string) (*Stat, error) {
	dir, n := path.Split(name)

	f := File{
		path: dir,
		name: n,
	}

	return f.Stat()
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
