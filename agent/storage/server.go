package store

import (
	"path"

	"github.com/samjohnduke/flextool/storage"
	context "golang.org/x/net/context"
)

type StoreServer struct {
	system storage.Filesystem
}

func (s *StoreServer) List(ctx context.Context, list *ListRequest) (*Blobs, error) {
	opts := storage.ListOpts{
		Recursive: list.Opts.Rescursive,
	}

	blobs, err := s.system.List(ctx, list.Path, opts)
	if err != nil {
		return nil, err
	}

	bs := &Blobs{}
	for _, b := range blobs {
		bs.Blob = append(bs.Blob, &Blob{
			Name:   path.Join(b.Path(), b.Name()),
			IsDir:  b.IsDir(),
			Exists: true,
		})
	}
	return bs, nil
}

func (s *StoreServer) Stat(context.Context, *StatRequest) (*StatResponse, error) {
	return nil, nil
}

func (s *StoreServer) Put(context.Context, *PutPartRequest) (*PutPartResponse, error) {
	return nil, nil
}

func (s *StoreServer) Get(context.Context, *GetPartRequest) (*GetPartResponse, error) {
	return nil, nil
}

func (s *StoreServer) Delete(context.Context, *DeleteRequest) (*Response, error) {
	return nil, nil
}

func (s *StoreServer) Move(context.Context, *MoveRequest) (*Response, error) {
	return nil, nil
}
