package archive

import (
	"fmt"
	"log"
	"path"
	"time"

	"github.com/pkg/errors"
	"github.com/samjohnduke/flextool/shared"
)

var cmdAll = "pg_dumpall"
var cmdOne = "pg_dump"

type PostgreSQL struct {
	Username string
	Options  []string
	Path     string
	All      bool
	DB       string
}

func (p *PostgreSQL) Archive() (*Archive, error) {
	ar := &Archive{}

	logs := make(chan string)
	go func() {
		for s := range logs {
			log.Println(s)
		}
	}()

	opts := p.opts()
	name := fmt.Sprintf("pgbackup-all-%v.sql.tar.gz", time.Now().Unix())
	path := path.Join(p.Path, name)

	opts = append(opts, fmt.Sprintf(`-f%v`, path))

	cmd := cmdAll
	if !p.All {
		cmd = cmdOne
	}

	status, err := shared.Exec(logs, cmd, opts...)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if len(status.Stderr) > 0 {
		return nil, errors.New("Error with script")
	}

	ar.Path = path
	ar.MIME = "application/x-tar"

	log.Printf(`Backup of database %s completed in: %f`, "all", status.Runtime)

	return ar, err
}

func (p *PostgreSQL) Restore(ar *Archive) error {
	panic("Unimplemented - go bug the maintainer")
}

func (p *PostgreSQL) opts() []string {
	options := p.Options

	if !p.All {
		options = append(options, "-Fc")
	}

	options = append(options, fmt.Sprintf(`-U%v`, p.Username))
	return options
}
