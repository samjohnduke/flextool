package storage

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"path"
)

type File struct {
	name        string
	path        string
	isDir       bool
	exists      bool
	loaded      bool
	inReadMode  bool
	inWriteMode bool
	fd          *os.File
}

func newFile(name string) *File {
	dir, n := path.Split(name)

	return &File{
		path:   dir,
		name:   n,
		isDir:  false,
		loaded: false,
		exists: false,
	}
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

func (f *File) Exists() bool {
	return f.exists
}

func (f *File) Stat() (*Stat, error) {
	fd, err := os.Open(path.Join(f.path, f.name))
	if err != nil {
		if os.IsNotExist(err) {
			f.exists = false
		}
		return nil, err
	}

	stat := &Stat{}

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

func (f *File) Read(p []byte) (n int, err error) {
	f.ensureReadable()
	return f.fd.Read(p)
}

func (f *File) Write(p []byte) (n int, err error) {
	f.ensureWriteable()
	return f.fd.Write(p)
}

func (f *File) Close() error {
	if f.fd == nil {
		return nil
	}

	return f.fd.Close()
}

func (f *File) Seek(offset int64, whence int) (int64, error) {
	return f.fd.Seek(offset, whence)
}

func (f *File) ensureReadable() error {
	if f.fd != nil {
		return nil
	}

	fd, err := os.Open(path.Join(f.path, f.name))
	if err != nil {
		log.Println(err)
		return err
	}

	f.fd = fd
	f.inReadMode = true
	f.loaded = true

	return nil
}

func (f *File) ensureWriteable() error {
	if f.fd != nil && f.inWriteMode {
		return nil
	}

	fd, err := os.Create(path.Join(f.path, f.name))
	if err != nil {
		log.Println(err)
		return err
	}

	f.fd = fd
	f.inWriteMode = true
	f.loaded = true
	return nil
}
