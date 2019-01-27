package sqlite

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/bvinc/go-sqlite-lite/sqlite3"
	"github.com/gofrs/uuid"
	"github.com/ifraixedes/go-payments-api-example/payment"
	"go.fraixed.es/errors"
)

// New creates an instance of the SQLite implementation of the payment Service.
//
// fname is any of the values that SQLite filename can take (see
// https://www.sqlite.org/c3ref/open.html, for more information), which are:
//
// * Path to a file, which is created if it doesn't exists
//
// * An URI as described in https://www.sqlite.org/uri.html
//
// * The string ":memory:", which creates an in-memory database
//
// When not using the URI fname, the SQLite database is always opened for
// read/write operations and with shared cache mode enabled (see
// https://www.sqlite.org/sharedcache.html for more information), otherwise URI
// parameters must specify the connection and cache mode through the query
// parameters.
// All the connections are opened with WAL method (see
// https://www.sqlite.org/wal.html for more information).
//
// The following error codes can be returned:
//
// * ErrInvalidArgDBFname
//
// * ErrDBCantOpen
//
// * payment.ErrUnexpectedStoreError
//
// * payment.ErrUnexpectedOSError - this error happens if there is an error when
// resolving the absolute path of the fname is a path to a file.
func New(fname string) (payment.Service, error) {
	if fname == "" {
		return nil, errors.New(ErrInvalidArgDBFname, payment.ErrMDArg("fname", fname))
	}

	var (
		svc = service{
			fname: fname,
		}
		isURI bool
	)

	switch {
	case fname == ":memory:":
		// When in-memory, the URI format must be used for being able to enabled the
		// shared cache
		svc.fname = "file::memory:?cache=shared"
	case !strings.HasPrefix(fname, "file:"):
		var err error
		svc.fname, err = filepath.Abs(fname)
		if err != nil {
			return nil, errors.Wrap(err, payment.ErrUnexpectedOSError, payment.ErrMDArg("fname", fname))
		}
		svc.openFlags = sqlite3.OPEN_READWRITE | sqlite3.OPEN_CREATE | sqlite3.OPEN_SHAREDCACHE
	default:
		isURI = true
	}

	var _, c, err = openConn(svc.fname, svc.openFlags)
	if err != nil {
		if c == sqlite3.ERROR && isURI {
			// Invalid URI format returns this error code
			return nil, errors.Wrap(err, ErrInvalidArgDBFname, payment.ErrMDArg("fname", fname))
		}

		return nil, err
	}

	return &svc, nil
}

type service struct {
	fname     string
	openFlags int
}

func (s *service) Create(_ context.Context, _ payment.PymtUpsert) (uuid.UUID, error) {
	// TODO: WIP
	// Implement it
	return uuid.Nil, nil
}

func (s *service) Delete(_ context.Context, _ uuid.UUID) error {
	// TODO: WIP
	// Implement it
	return nil
}

func (s *service) Find(
	_ context.Context, _ payment.Filter, _ payment.Selection, _ payment.Sort, _ payment.Chunk,
) ([]payment.Pymt, error) {
	// TODO: WIP
	// Implement it
	return nil, nil
}

func (s *service) Get(_ context.Context, _ uuid.UUID, _ payment.Selection) (payment.Pymt, error) {
	// TODO: WIP
	// Implement it
	return payment.Pymt{}, nil
}

func (s *service) Update(_ context.Context, _ uuid.UUID, _ uint32, _ payment.PymtUpsert) error {
	// TODO: WIP
	// Implement it
	return nil
}

// openConn create a new sqlite3 connection.
// It returns error the connection creation fails or the WAL journal model cannot
// be set. When an error is returned, the sqlite3 error code is also returned.
//
// The following error codes can be returned:
//
// * ErrDBCantOpen
//
// * payment.ErrUnexpectedStoreError
func openConn(fname string, flags int) (*sqlite3.Conn, int, error) {
	var (
		err  error
		conn *sqlite3.Conn
	)
	if flags == 0 {
		conn, err = sqlite3.Open(fname)
	} else {
		conn, err = sqlite3.Open(fname, flags)
	}

	if err != nil {
		var serr, ok = err.(*sqlite3.Error)
		if !ok {
			return nil, 0, errors.Wrap(err, payment.ErrUnexpectedStoreError,
				payment.ErrMDVar("sqlite_open_filename", fname),
				payment.ErrMDVar("sqlite_open_flags", flags),
			)
		}

		switch c := serr.Code(); c {
		case sqlite3.CANTOPEN:
			return nil, c, errors.Wrap(serr, ErrDBCantOpen, payment.ErrMDArg("fname", fname))
		default:
			return nil, c, errors.Wrap(serr, payment.ErrUnexpectedStoreError,
				payment.ErrMDVar("sqlite_open_filename", fname),
				payment.ErrMDVar("sqlite_open_flags", flags),
			)
		}
	}

	if err = conn.Exec("PRAGMA journal_mode=wal"); err != nil {
		var c int
		if serr, ok := err.(*sqlite3.Error); !ok {
			c = serr.Code()
		}

		return nil, c, errors.Wrap(err, payment.ErrUnexpectedStoreError,
			payment.ErrMDFnCall("sqlite3.Conn.Exec", "PRAGMA journal_mode=wal"),
		)
	}

	return conn, 0, nil
}
