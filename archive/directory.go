package archive

type Directory struct {
	Out         string
	In          string
	Compression string
}

func (d *Directory) Archive() (*Archive, error) {
	return nil, nil
}

func (d *Directory) Restore(ar *Archive) error {
	return nil
}
