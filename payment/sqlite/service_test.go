package sqlite_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/ifraixedes/go-payments-api-example/payment"
	"github.com/ifraixedes/go-payments-api-example/payment/internal/testutil"
	"github.com/ifraixedes/go-payments-api-example/payment/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.fraixed.es/errors"
)

func TestNew(t *testing.T) {
	type params struct {
		fname string
	}
	type tcase struct {
		desc   string
		args   params
		assert func(*testing.T, tcase, payment.Service, error)
		after  func(*testing.T, tcase)
	}

	var tcases = []tcase{
		{
			desc: "successful: file path",
			args: params{
				fname: func() string {
					f, err := ioutil.TempFile(os.TempDir(), "pymt-api-ex-sqlite-*.db")
					require.NoError(t, err)
					return f.Name()
				}(),
			},
			assert: func(t *testing.T, _ tcase, s payment.Service, err error) {
				assert.NotNil(t, s)
				assert.NoError(t, err)
			},
			after: func(t *testing.T, tc tcase) {
				require.NoError(t, os.Remove(tc.args.fname))
			},
		},
		{
			desc: "successful: in-memory",
			args: params{
				fname: ":memory:",
			},
			assert: func(t *testing.T, _ tcase, s payment.Service, err error) {
				assert.NotNil(t, s)
				assert.NoError(t, err)
			},
		},
		{
			desc: "successful: URI",
			args: params{
				fname: func() string {
					f, err := ioutil.TempFile(os.TempDir(), "pymt-api-ex-sqlite-*.db")
					require.NoError(t, err)
					return fmt.Sprintf("file://%s?cache=private", f.Name())
				}(),
			},
			assert: func(t *testing.T, _ tcase, s payment.Service, err error) {
				assert.NotNil(t, s)
				assert.NoError(t, err)
			},
			after: func(t *testing.T, tc tcase) {
				var m = regexp.MustCompile(`file://(.[^?]+)\?`).FindStringSubmatch(tc.args.fname)
				require.NoError(t, os.Remove(m[1]))
			},
		},
		{
			desc: "error: empty string",
			args: params{
				fname: "",
			},
			assert: func(t *testing.T, tc tcase, _ payment.Service, err error) {
				testutil.AssertError(t, err, sqlite.ErrInvalidArgDBFname, payment.ErrMDArg("fname", tc.args.fname))
			},
		},
		{
			desc: "error: invalid file path",
			args: params{
				fname: "/do-not/exist/path/some.db",
			},
			assert: func(t *testing.T, tc tcase, _ payment.Service, err error) {
				testutil.AssertError(t, err, sqlite.ErrDBCantOpen, payment.ErrMDArg("fname", tc.args.fname))
			},
		},
		{
			desc: "error: invalid URI",
			args: params{
				fname: "file://example.com/tmp/some.db",
			},
			assert: func(t *testing.T, tc tcase, _ payment.Service, err error) {
				testutil.AssertError(t, err, sqlite.ErrInvalidArgDBFname, payment.ErrMDArg("fname", tc.args.fname))
			},
		},
		{
			desc: "error: invalid DB (dir)",
			args: params{
				fname: "/",
			},
			assert: func(t *testing.T, tc tcase, _ payment.Service, err error) {
				testutil.AssertError(t, err, sqlite.ErrDBCantOpen, payment.ErrMDArg("fname", tc.args.fname))
			},
		},
	}

	for i := range tcases {
		var tc = tcases[i]
		t.Run(tc.desc, func(t *testing.T) {
			var s, err = sqlite.New(tc.args.fname)
			tc.assert(t, tc, s, err)

			if tc.after != nil {
				tc.after(t, tc)
			}
		})
	}
}

func TestService_Create(t *testing.T) {
	t.Run("error invalid paymet", func(t *testing.T) {
		var s, err = sqlite.New(testingDB)
		require.NoError(t, err)

		var p = payment.PymtUpsert{
			Type:  "Payment",
			OrgID: testutil.NewUUID(t),
		}
		var ec, ok = errors.GetCode(p.Validate())
		require.True(t, ok)

		_, err = s.Create(context.Background(), p)
		testutil.AssertError(t, err, ec)
	})
}

func TestService_Delete(t *testing.T) {
	t.Run("error nil UUID", func(t *testing.T) {
		var s, err = sqlite.New(testingDB)
		require.NoError(t, err)

		err = s.Delete(context.Background(), uuid.Nil)
		testutil.AssertError(t, err, payment.ErrInvalidPaymentID, payment.ErrMDArg("id", uuid.Nil))
	})
}

func TestService_Get(t *testing.T) {
	t.Run("error nil UUID", func(t *testing.T) {
		var s, err = sqlite.New(testingDB)
		require.NoError(t, err)

		_, err = s.Get(context.Background(), uuid.Nil, payment.Selection{})
		testutil.AssertError(t, err, payment.ErrInvalidPaymentID, payment.ErrMDArg("id", uuid.Nil))
	})
}

func TestService_Update(t *testing.T) {
	t.Run("error nil UUID", func(t *testing.T) {
		var s, err = sqlite.New(testingDB)
		require.NoError(t, err)

		err = s.Update(context.Background(), uuid.Nil, 0, payment.PymtUpsert{})
		testutil.AssertError(t, err, payment.ErrInvalidPaymentID, payment.ErrMDArg("id", uuid.Nil))
	})

	t.Run("error invalid payment upsert", func(t *testing.T) {
		var s, err = sqlite.New(testingDB)
		require.NoError(t, err)

		var p = payment.PymtUpsert{}
		var ec, ok = errors.GetCode(p.Validate())
		require.True(t, ok)

		err = s.Update(context.Background(), testutil.NewUUID(t), 0, p)
		testutil.AssertError(t, err, ec)
	})
}
