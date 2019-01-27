package sqlite

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/bvinc/go-sqlite-lite/sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	type params struct {
		fname string
	}
	type tcase struct {
		desc   string
		args   params
		assert func(*testing.T, tcase, *service)
		after  func(*testing.T, tcase)
	}

	var tcases = []tcase{
		{
			desc: "file path",
			args: params{
				fname: func() string {
					f, err := ioutil.TempFile(os.TempDir(), "pymt-api-ex-sqlite-*.db")
					require.NoError(t, err)
					return f.Name()
				}(),
			},
			assert: func(t *testing.T, tc tcase, s *service) {
				if !assert.NotNil(t, s) {
					return
				}

				var fpath, err = filepath.Abs(tc.args.fname)
				require.NoError(t, err)

				assert.Equal(t, fpath, s.fname)
				assert.Equal(t, sqlite3.OPEN_READWRITE|sqlite3.OPEN_CREATE|sqlite3.OPEN_SHAREDCACHE, s.openFlags)
			},
			after: func(t *testing.T, tc tcase) {
				require.NoError(t, os.Remove(tc.args.fname))
			},
		},
		{
			desc: "in-memory",
			args: params{
				fname: ":memory:",
			},
			assert: func(t *testing.T, tc tcase, s *service) {
				if !assert.NotNil(t, s) {
					return
				}

				assert.Equal(t, "file::memory:?cache=shared", s.fname)
				assert.Equal(t, 0, s.openFlags)
			},
		},
		{
			desc: "URI",
			args: params{
				fname: func() string {
					f, err := ioutil.TempFile(os.TempDir(), "pymt-api-ex-sqlite-*.db")
					require.NoError(t, err)
					return fmt.Sprintf("file://%s?cache=private", f.Name())
				}(),
			},
			assert: func(t *testing.T, tc tcase, s *service) {
				if !assert.NotNil(t, s) {
					return
				}

				assert.Equal(t, tc.args.fname, s.fname)
				assert.Equal(t, 0, s.openFlags)
			},
			after: func(t *testing.T, tc tcase) {
				var m = regexp.MustCompile(`file://(.[^?]+)\?`).FindStringSubmatch(tc.args.fname)
				require.NoError(t, os.Remove(m[1]))
			},
		},
	}

	for i := range tcases {
		var tc = tcases[i]
		t.Run(tc.desc, func(t *testing.T) {
			var s, err = New(tc.args.fname)
			require.NoError(t, err)

			tc.assert(t, tc, s.(*service))

			if tc.after != nil {
				tc.after(t, tc)
			}
		})
	}
}
