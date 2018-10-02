package storage

import (
	"context"
	"path"

	"github.com/minio/minio-go"
)

type DOBlob struct {
	isDir  bool
	name   string
	path   string
	loaded bool
	exists bool

	object *minio.Object
	space  *DOSpace
}

func newDOBlob(name string) *DOBlob {
	dir, n := path.Split(name)

	b := &DOBlob{
		path:   dir,
		name:   n,
		exists: false,
		loaded: false,
		isDir:  false,
	}

	b.Stat()

	return b
}

func (f *DOBlob) Name() string {
	return f.name
}

func (f *DOBlob) Path() string {
	return f.path
}

func (f *DOBlob) IsDir() bool {
	return f.isDir
}

func (f *DOBlob) Exists() bool {
	return f.exists
}

func (f *DOBlob) Stat() (*Stat, error) {
	stat, err := f.space.Stat(context.Background(), path.Join(f.path, f.name))
	if err != nil {
		f.exists = true
	}

	if stat.Size <= 0 {
		f.isDir = true
	}

	return stat, err
}

func (f *DOBlob) Read(p []byte) (n int, err error) {
	err = f.ensureObject()
	if err != nil {
		return 0, err
	}

	return f.object.Read(p)
}

func (f *DOBlob) Write(p []byte) (n int, err error) {
	return 0, nil
}

func (f *DOBlob) Close() (err error) {
	err = f.ensureObject()
	if err != nil {
		return err
	}

	err = f.object.Close()
	if err != nil {
		return err
	}

	f.object = nil
	f.loaded = false

	return nil
}

func (f *DOBlob) Seek(offset int64, whence int) (n int64, err error) {
	err = f.ensureObject()
	if err != nil {
		return
	}

	return f.object.Seek(offset, whence)
}

func (f *DOBlob) ensureObject() error {
	if f.object != nil {
		return nil
	}

	object, err := f.space.getObject(path.Join(f.path, f.name))
	if err != nil {
		return err
	}

	f.object = object
	f.loaded = true
	return nil
}
