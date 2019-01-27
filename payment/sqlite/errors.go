package sqlite

type code uint8

// The list of specific error codes that the sqlite service implementation can
// return.
const (
	ErrDBCantOpen code = iota + 1
	ErrInvalidArgDBFname
)

func (c code) String() string {
	switch c {
	case ErrInvalidArgDBFname:
		return "InvalidArgDBFname"
	case ErrDBCantOpen:
		return "DBCantOpen"
	}

	return ""
}

func (c code) Message() string {
	switch c {
	case ErrInvalidArgDBFname:
		return "The SQLite filename isn't of a valid format"
	case ErrDBCantOpen:
		return "The SQLite DB cannot be open. Causes: DB is corrupted, file " +
			"isn't a valid DB file or some DB files cannot be open"
	}

	return ""
}
