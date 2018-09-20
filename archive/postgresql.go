package archive

import (
	"fmt"
	"log"

	"github.com/samjohnduke/flextool/shared"
)

var CMD = "pg_dump"

type PostgreSQL struct {
	Host     string
	Port     string
	Username string
	Password string

	Options []string
	Out     string

	DB string
}

func (p *PostgreSQL) Archive() (*Archive, error) {
	ar := &Archive{}

	logs := make(chan string, 100)
	go func() {
		for s := range logs {
			log.Println(s)
		}
	}()

	opts := p.opts()
	status, err := shared.Exec(logs, CMD, opts...)
	if err != nil {
		return nil, err
	}

	ar.Path = p.Out
	ar.MIME = "application/x-tar"

	log.Printf(`Backup of database %s completed in: %f`, p.DB, status.Runtime)

	return ar, err
}

func (p *PostgreSQL) Restore(ar *Archive) error {
	panic("Unimplemented - go bug the maintainer")
}

func (p *PostgreSQL) opts() []string {
	options := p.Options
	options = append(options, fmt.Sprintf(`-d%v`, p.DB))
	options = append(options, fmt.Sprintf(`-h%v`, p.Host))
	options = append(options, fmt.Sprintf(`-p%v`, p.Port))
	options = append(options, fmt.Sprintf(`-U%v`, p.Username))
	return options
}
