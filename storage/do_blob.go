package storage

import "path"

type DOBlob struct {
	isDir  bool
	name   string
	path   string
	loaded bool

	space *DOSpace
}

func newDOBlob(name string) *DOBlob {
	dir, n := path.Split(name)

	return &DOBlob{
		path: dir,
		name: n,
	}
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
	return false
}

func (f *DOBlob) Stat() (*Stat, error) {
	return nil, nil
}

func (f *DOBlob) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (f *DOBlob) Write(p []byte) (n int, err error) {
	return 0, nil
}

func (f *DOBlob) Close() error {
	return nil
}

func (f *DOBlob) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}
