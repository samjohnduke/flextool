package archive

// Archiver is implemented as a means of archiving a specific tool into a
// single file that can be shipped off to another location
type Archiver interface {
	Archive() (*Archive, error)
	Restore(*Archive) error
}

// Archive is a location on disk + a MIME type that describes the archive
type Archive struct {
	Path string
	MIME string
}
