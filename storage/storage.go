package storage

type Store interface {
	List()
	Upload()
	Download()
	Sync()
	SyncDir()
	Delete()
	Move()
}
