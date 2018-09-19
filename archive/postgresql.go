package archive

type PostgreSQL struct {
	Out string
}

func (p *PostgreSQL) Archive() (*Archive, error) {
	return nil, nil
}

func (p *PostgreSQL) Restore(ar *Archive) error {
	return nil
}
