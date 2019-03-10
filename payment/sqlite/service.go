package sqlite

import (
	"context"
	"path/filepath"
	"strings"
	"time"

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
//   resolving the absolute path of the fname is a path to a file.
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

	var c, pc, err = svc.openConn(context.Background())
	if err != nil {
		if pc == sqlite3.ERROR && isURI {
			// Invalid URI format returns this error code
			return nil, errors.Wrap(err, ErrInvalidArgDBFname, payment.ErrMDArg("fname", fname))
		}

		return nil, err
	}

	if err := c.Close(); err != nil {
		return nil, errors.Wrap(err, payment.ErrUnexpectedStoreError)
	}

	return &svc, nil
}

type service struct {
	fname     string
	openFlags int
}

// TODO: methods should consider ctx.Cancel() using SQLite3 conn.Interrupt()

// Create stores p in the database.
//
// The function will return all the errors that payment.Service documents plus
// the following ones:
//
// * ErrDBCantOpen
//
// * ErrDBLimit
//
// * ErrDBSchemaChanged
//
// * ErrInvalidPayment
func (s *service) Create(ctx context.Context, p payment.PymtUpsert) (uuid.UUID, error) {
	if err := p.Validate(); err != nil {
		return uuid.Nil, err
	}

	var id, err = uuid.NewV4()
	if err != nil {
		return uuid.Nil, errors.Wrap(err, payment.ErrUnexpectedStoreError)
	}

	var pd []byte
	{
		d := &pymtData{}
		d.Init(p)
		pd, err = d.Serialize()
		if err != nil {
			return uuid.Nil, err
		}
	}

	conn, _, err := s.openConn(ctx)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, payment.ErrUnexpectedStoreError)
	}

	defer func() {
		_ = conn.Close()
	}()

	err = conn.Exec(
		"INSERT INTO payments(id, organisation_id, data) VALUES (?, ?, ?)",
		id.String(), p.OrgID.String(), pd,
	)
	if err != nil {
		if cerr := handleCommonSQLiteErr(err); cerr != nil {
			return uuid.Nil, cerr
		}

		var pc, _, serr = isSQLiteErr(err)
		if serr != nil {
			if pc == sqlite3.CONSTRAINT {
				return uuid.Nil, errors.Wrap(serr, ErrInvalidPayment)
			}
		}

		return uuid.Nil, errors.Wrap(err, payment.ErrUnexpectedStoreError)
	}

	return id, nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New(payment.ErrInvalidPaymentID, payment.ErrMDArg("id", id))
	}

	conn, _, err := s.openConn(ctx)
	if err != nil {
		return errors.Wrap(err, payment.ErrUnexpectedStoreError)
	}

	defer func() {
		_ = conn.Close()
	}()

	err = conn.Exec("DELETE FROM payments WHERE id = ?", id.String())
	if err != nil {
		if cerr := handleCommonSQLiteErr(err); cerr != nil {
			return cerr
		}

		return errors.Wrap(err, payment.ErrUnexpectedStoreError)
	}

	if conn.TotalChanges() == 0 {
		return errors.New(payment.ErrNotFound, payment.ErrMDVar("id", id))
	}

	return nil
}

func (s *service) Find(
	_ context.Context, _ payment.Filter, _ payment.Selection, _ payment.Sort, _ payment.Chunk,
) ([]payment.Pymt, error) {
	// TODO: WIP
	// Implement it
	return nil, nil
}

func (s *service) Get(ctx context.Context, id uuid.UUID, _ payment.Selection) (payment.Pymt, error) {
	if id == uuid.Nil {
		return payment.Pymt{}, errors.New(payment.ErrInvalidPaymentID, payment.ErrMDArg("id", id))
	}
	// TODO: WIP
	// Implement it
	return payment.Pymt{}, nil
}

func (s *service) Update(_ context.Context, id uuid.UUID, _ uint32, p payment.PymtUpsert) error {
	if id == uuid.Nil {
		return errors.New(payment.ErrInvalidPaymentID, payment.ErrMDArg("id", id))
	}

	if err := p.Validate(); err != nil {
		return err
	}

	// TODO: WIP
	// Implement it
	return nil
}

// openConn create a new sqlite3 connection.
// It returns error the connection creation fails or the WAL journal model cannot
// be set. When an error is returned, the sqlite3 error primary code is also
// returned (see https://www.sqlite.org/rescode.html).
//
// If ctx has a deadline, it's used to set the BusyTimeout to the connection,
// see https://godoc.org/github.com/bvinc/go-sqlite-lite/sqlite3#Conn.BusyFunc
//
// The following error codes can be returned:
//
// * ErrDBCantOpen
//
// * payment.ErrUnexpectedStoreError
func (s *service) openConn(ctx context.Context) (*sqlite3.Conn, uint8, error) {
	var (
		err  error
		conn *sqlite3.Conn
	)
	if s.openFlags == 0 {
		conn, err = sqlite3.Open(s.fname)
	} else {
		conn, err = sqlite3.Open(s.fname, s.openFlags)
	}

	if err != nil {
		var pc, _, serr = isSQLiteErr(err)
		if serr == nil {
			return nil, 0, errors.Wrap(err, payment.ErrUnexpectedStoreError,
				payment.ErrMDVar("sqlite_open_filename", s.fname),
				payment.ErrMDVar("sqlite_open_flags", s.openFlags),
			)
		}

		switch pc {
		case sqlite3.CANTOPEN, sqlite3.NOTADB:
			return nil, pc, errors.Wrap(serr, ErrDBCantOpen, payment.ErrMDArg("fname", s.fname))
		default:
			return nil, pc, errors.Wrap(serr, payment.ErrUnexpectedStoreError,
				payment.ErrMDVar("sqlite_open_filename", s.fname),
				payment.ErrMDVar("sqlite_open_flags", s.openFlags),
			)
		}
	}

	if err = conn.Exec("PRAGMA journal_mode=wal"); err != nil {
		pc, _, _ := isSQLiteErr(err)

		return nil, pc, errors.Wrap(err, payment.ErrUnexpectedStoreError,
			payment.ErrMDFnCall("sqlite3.Conn.Exec", "PRAGMA journal_mode=wal"),
		)
	}

	if t, ok := ctx.Deadline(); ok {
		conn.BusyTimeout(time.Until(t))
	}

	return conn, 0, nil
}

// isSQLiteErr returns the primary error code, the extended error code and the
// specific sqlite3 error when it's of such type or 0, 0 and nil when not.
// See https://www.sqlite.org/rescode.html
func isSQLiteErr(err error) (uint8, int, *sqlite3.Error) {
	if serr, ok := err.(*sqlite3.Error); ok {
		c := serr.Code()
		return uint8(c & 0xFF), c, serr
	}

	return 0, 0, nil
}

// handleCommonSQLiteErr is a convenient function to map SQLite error codes
// which are common to the most of the SQL statements executed by the service
// methods. It returns an error is the err is successful mapped otherwise nil,
// which means that the caller should deal with it because err isn't a common
// SQLite error.
func handleCommonSQLiteErr(err error) error {
	var pc, _, serr = isSQLiteErr(err)
	if serr == nil {
		return errors.Wrap(err, payment.ErrUnexpectedStoreError)
	}

	switch pc {
	case sqlite3.BUSY, sqlite3.INTERRUPT, sqlite3.LOCKED:
		return errors.Wrap(serr, payment.ErrAbortedOperation)
	case sqlite3.SCHEMA:
		return errors.Wrap(serr, ErrDBSchemaChanged)
	case sqlite3.CORRUPT, sqlite3.IOERR, sqlite3.FULL, sqlite3.INTERNAL, sqlite3.PROTOCOL:
		return errors.Wrap(serr, payment.ErrUnexpectedStoreError)
	case sqlite3.NOMEM:
		return errors.Wrap(serr, payment.ErrUnexpectedSysError)
	case sqlite3.TOOBIG:
		return errors.Wrap(err, ErrDBLimit)
	}

	return nil
}
