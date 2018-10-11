package storage

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
)

type Filesystem struct {
	root string
}

func NewFilesystem(root string) (Store, error) {
	store := &Filesystem{
		root: root,
	}

	return store, nil
}

func (store *Filesystem) New(ctx context.Context, name string) Blob {
	log.Println(name)
	f := newFile(path.Join(store.root, name))

	f.Stat()

	return f
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
			path:   prefix,
			name:   f.Name(),
			isDir:  f.IsDir(),
			loaded: false,
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
	f := store.New(ctx, name)
	return f.Stat()
}

func (store *Filesystem) Put(ctx context.Context, name string, blob Blob) error {
	fd, err := os.Open(name)
	if err != nil {
		return err
	}

	if _, err = io.Copy(fd, blob); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (store *Filesystem) Get(ctx context.Context, name string, blob Blob) error {
	fd, err := os.Open(name)
	if err != nil {
		return err
	}

	if _, err = io.Copy(blob, fd); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (store *Filesystem) Delete(ctx context.Context, name string) error {
	fd, err := os.Open(name)
	if err != nil {
		return err
	}

	stat, err := fd.Stat()
	if err != nil {
		return err
	}

	if stat.IsDir() {
		err = os.RemoveAll(name)
	} else {
		err = os.Remove(name)
	}

	if err != nil {
		return err
	}

	return nil
}

func (store *Filesystem) Move(ctx context.Context, from string, to string) error {
	err := os.Rename(from, to)
	if err != nil {
		return err
	}
	return nil
}
