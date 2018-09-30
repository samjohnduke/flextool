package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/minio/minio-go"
)

type DOSpace struct {
	accessKey string
	secretKey string
	endpoint  string
	space     string
	region    string
	secure    bool

	client *minio.Client
}

func NewDOSpace(accessKey, secretKey, region, space string, secure bool) (Store, error) {
	store := &DOSpace{
		accessKey: accessKey,
		secretKey: secretKey,
		endpoint:  fmt.Sprintf("%s.digitaloceanspaces.com", region),
		region:    region,
		space:     space,
		secure:    secure,
	}

	minioClient, err := minio.New(store.endpoint, store.accessKey, store.secretKey, store.secure)
	if err != nil {
		log.Fatalln(err)
	}

	store.client = minioClient
	return store, nil
}

func (dos *DOSpace) New(ctx context.Context, name string) Blob {
	f := NewFile(name)
	f.Stat()

	return f
}

// TODO - check for leading slash
// TODO - check for if prefix is a directory with missing trailing slash
func (dos *DOSpace) List(ctx context.Context, prefix string, opts ListOpts) (Blobs, error) {
	blobs := make(Blobs)
	doneCh := make(chan struct{})

	// Indicate to our routine to exit cleanly upon return.
	defer close(doneCh)

	objectCh := dos.client.ListObjectsV2(dos.space, prefix, opts.Recursive, doneCh)
	for object := range objectCh {
		if object.Err != nil {
			fmt.Println(object.Err)
			return nil, object.Err
		}

		file := &File{
			name:  object.Key,
			isDir: object.Size == 0 && strings.HasSuffix(object.Key, "/"),
			path:  "",
		}

		blobs[object.Key] = file
	}

	return blobs, nil
}

func (dos *DOSpace) Stat(ctx context.Context, name string) (*Stat, error) {
	info, err := dos.client.StatObject(dos.space, name, minio.StatObjectOptions{})
	if err != nil {
		return nil, err
	}

	stat := &Stat{
		LastModified: info.LastModified,
		ContentType:  info.ContentType,
		Size:         info.Size,
		ETag:         info.ETag,
	}

	return stat, nil
}

func (dos *DOSpace) Put(ctx context.Context, name string, blob Blob) error {
	stat, err := blob.Stat()
	if err != nil {
		return err
	}

	n, err := dos.client.PutObjectWithContext(ctx, dos.space, name, blob, stat.Size, minio.PutObjectOptions{
		ContentType: stat.ContentType,
	})
	if err != nil {
		return err
	}

	fmt.Println("Successfully uploaded bytes: ", n)
	return nil
}

func (dos *DOSpace) Get(ctx context.Context, name string, blob Blob) error {
	object, err := dos.client.GetObjectWithContext(ctx, dos.space, name, minio.GetObjectOptions{})
	if err != nil {
		fmt.Println(err)
		return err
	}

	if _, err = io.Copy(blob, object); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (dos *DOSpace) Delete(ctx context.Context, name string) error {
	return nil
}

func (dos *DOSpace) Move(ctx context.Context, name string, to string) error {
	return nil
}
