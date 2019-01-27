package payment_test

import (
	"math/rand"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/ifraixedes/go-payments-api-example/payment"
	"github.com/ifraixedes/go-payments-api-example/payment/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFilter(t *testing.T) {
	type params struct {
		op payment.FilterLogical
		l  payment.Filter
		r  payment.Filter
	}

	type tcase struct {
		desc   string
		args   params
		assert func(*testing.T, tcase, payment.Filter, error)
	}

	var fnuid, err = payment.NewFilterByID(payment.FilterCmpNotEqual, uuid.Nil)
	require.NoError(t, err)
	ftdd, err := payment.NewFilterByType(payment.FilterCmpEqual, "Payment")
	require.NoError(t, err)

	// Set a value which isn't coined by any of the package constants
	var fopinv payment.FilterLogical
	// nolint:gosec
	if rand.Int()%2 == 0 {
		fopinv = payment.FilterLogical(rand.Intn(245) + 10)
	}

	var tcases = []tcase{
		{
			desc: "successful",
			args: params{
				op: payment.FilterLogicalOr,
				l:  ftdd,
				r:  fnuid,
			},
			assert: func(t *testing.T, _ tcase, f payment.Filter, err error) {
				assert.NotEqual(t, payment.Filter{}, f)
				assert.NoError(t, err)
			},
		},
		{
			desc: "error: invalid operator",
			args: params{
				op: fopinv,
				l:  ftdd,
				r:  fnuid,
			},
			assert: func(t *testing.T, tc tcase, f payment.Filter, err error) {
				testutil.AssertError(t, err, payment.ErrInvalidArgFilterLogicalOpNotExists, payment.ErrMDArg("op", tc.args.op))
			},
		},
		{
			desc: "error: empty left filter",
			args: params{
				op: payment.FilterLogicalAnd,
				l:  payment.Filter{},
				r:  fnuid,
			},
			assert: func(t *testing.T, tc tcase, _ payment.Filter, err error) {
				testutil.AssertError(t, err, payment.ErrInvalidArgFilterNodeEmpty, payment.ErrMDArg("left", tc.args.l))
			},
		},
		{
			desc: "error: empty left filter",
			args: params{
				op: payment.FilterLogicalAnd,
				l:  fnuid,
				r:  payment.Filter{},
			},
			assert: func(t *testing.T, tc tcase, _ payment.Filter, err error) {
				testutil.AssertError(t, err, payment.ErrInvalidArgFilterNodeEmpty, payment.ErrMDArg("right", tc.args.r))
			},
		},
	}

	for i := range tcases {
		var tc = tcases[i]
		t.Run(tc.desc, func(t *testing.T) {
			var f, err = payment.NewFilter(tc.args.op, tc.args.l, tc.args.r)
			tc.assert(t, tc, f, err)
		})
	}
}

func TestFilter_NodeType(t *testing.T) {
	t.Run("leaf", func(t *testing.T) {
		var f, err = payment.NewFilterByAmount(payment.FilterCmpEqual, 10)
		require.NoError(t, err)

		assert.Equal(t, payment.FilterNodeTypeLeaf, f.NodeType())
	})

	t.Run("non-leaf", func(t *testing.T) {
		var f1, err = payment.NewFilterByAmount(payment.FilterCmpEqual, 10)
		require.NoError(t, err)
		f2, err := payment.NewFilterByType(payment.FilterCmpEqual, "Payment")
		require.NoError(t, err)
		f, err := payment.NewFilter(payment.FilterLogicalAnd, f1, f2)
		require.NoError(t, err)

		assert.Equal(t, payment.FilterNodeTypeNonLeaf, f.NodeType())
	})

	t.Run("empty", func(t *testing.T) {
		var f = payment.Filter{}
		assert.Equal(t, payment.FilterNodeTypeEmpty, f.NodeType())
	})
}

func TestFilter_Nodes(t *testing.T) {
	var fempty = payment.Filter{}

	t.Run("leaf", func(t *testing.T) {
		var f, err = payment.NewFilterByAmount(payment.FilterCmpEqual, 10)
		require.NoError(t, err)

		var _, l, r = f.Nodes()
		assert.Equal(t, fempty, l)
		assert.Equal(t, fempty, r)
	})

	t.Run("non-leaf", func(t *testing.T) {
		var f1, err = payment.NewFilterByAmount(payment.FilterCmpEqual, 10)
		require.NoError(t, err)
		f2, err := payment.NewFilterByType(payment.FilterCmpEqual, "Payment")
		require.NoError(t, err)
		f, err := payment.NewFilter(payment.FilterLogicalAnd, f1, f2)
		require.NoError(t, err)

		var op, l, r = f.Nodes()
		assert.Equal(t, payment.FilterLogicalAnd, op)
		assert.Equal(t, f1, l)
		assert.Equal(t, f2, r)
	})

	t.Run("empty", func(t *testing.T) {
		var _, l, r = fempty.Nodes()
		assert.Equal(t, fempty, l)
		assert.Equal(t, fempty, r)
	})
}

func TestFilter_Leaf(t *testing.T) {
	t.Run("leaf", func(t *testing.T) {
		var f, err = payment.NewFilterByAmount(payment.FilterCmpGreaterThan, 85)
		require.NoError(t, err)

		var l = f.Leaf()
		require.NotNil(t, l)
		assert.True(t, l.IsSet())

		var op, val = l.Filter()
		assert.Equal(t, payment.FilterCmpGreaterThan, op)
		assert.Equal(t, 85.0, val)
	})

	t.Run("non-leaf", func(t *testing.T) {
		var f1, err = payment.NewFilterByAmount(payment.FilterCmpEqual, -10.4)
		require.NoError(t, err)
		f2, err := payment.NewFilterByType(payment.FilterCmpEqual, "Payment")
		require.NoError(t, err)
		f, err := payment.NewFilter(payment.FilterLogicalAnd, f1, f2)
		require.NoError(t, err)

		var l = f.Leaf()
		assert.Nil(t, l)
	})

	t.Run("empty", func(t *testing.T) {
		var (
			f = payment.Filter{}
			l = f.Leaf()
		)

		assert.Nil(t, l)
	})
}

func TestNewFilterByAmount(t *testing.T) {
	type params struct {
		cmp payment.FilterCmp
		val float64
	}

	type tcase struct {
		desc   string
		args   params
		assert func(*testing.T, tcase, payment.Filter, error)
	}

	var tcases = []tcase{
		{
			desc: "successful",
			args: params{
				cmp: payment.FilterCmp(rand.Intn(5) + 1),
				val: rand.Float64(),
			},
			assert: func(t *testing.T, tc tcase, f payment.Filter, err error) {
				assert.NoError(t, err)
				if assert.Equal(t, f.NodeType(), payment.FilterNodeTypeLeaf) {
					var l = f.Leaf()
					assert.True(t, l.IsSet())

					var cmp, val = l.Filter()
					assert.Equal(t, tc.args.cmp, cmp)
					assert.Equal(t, tc.args.val, val)
				}
			},
		},
		{
			desc: "error: unsupported cmp",
			args: params{
				cmp: payment.FilterCmpMatch,
				val: rand.Float64(),
			},
			assert: func(t *testing.T, tc tcase, _ payment.Filter, err error) {
				testutil.AssertError(t, err, payment.ErrInvalidArgFilterCmpNotSupported, payment.ErrMDArg("cmp", tc.args.cmp))
			},
		},
		{
			desc: "error: cmp doesn't exist",
			args: params{
				cmp: payment.FilterCmp(rand.Intn(240) + 15),
				val: rand.Float64(),
			},
			assert: func(t *testing.T, tc tcase, _ payment.Filter, err error) {
				testutil.AssertError(t, err, payment.ErrInvalidArgFilterCmpNotExists, payment.ErrMDArg("cmp", tc.args.cmp))
			},
		},
	}

	for i := range tcases {
		var tc = tcases[i]
		t.Run(tc.desc, func(t *testing.T) {
			var f, err = payment.NewFilterByAmount(tc.args.cmp, tc.args.val)
			tc.assert(t, tc, f, err)
		})
	}
}

func TestNewFilterByID(t *testing.T) {
	type params struct {
		cmp payment.FilterCmp
		val uuid.UUID
	}

	type tcase struct {
		desc   string
		args   params
		assert func(*testing.T, tcase, payment.Filter, error)
	}

	var tcases = []tcase{
		{
			desc: "successful",
			args: params{
				cmp: func() payment.FilterCmp {
					// nolint:gosec
					if rand.Int()%2 == 0 {
						return payment.FilterCmpEqual
					}

					return payment.FilterCmpNotEqual
				}(),
				val: uuid.Must(uuid.NewV4()),
			},
			assert: func(t *testing.T, tc tcase, f payment.Filter, err error) {
				assert.NoError(t, err)
				if assert.Equal(t, f.NodeType(), payment.FilterNodeTypeLeaf) {
					var l = f.Leaf()
					assert.True(t, l.IsSet())

					var cmp, val = l.Filter()
					assert.Equal(t, tc.args.cmp, cmp)
					assert.Equal(t, tc.args.val, val)
				}
			},
		},
		{
			desc: "error: unsupported cmp",
			args: params{
				cmp: payment.FilterCmp(rand.Intn(4) + 3),
				val: uuid.Must(uuid.NewV4()),
			},
			assert: func(t *testing.T, tc tcase, _ payment.Filter, err error) {
				testutil.AssertError(t, err, payment.ErrInvalidArgFilterCmpNotSupported, payment.ErrMDArg("cmp", tc.args.cmp))
			},
		},
		{
			desc: "error: cmp doesn't exist",
			args: params{
				cmp: payment.FilterCmp(rand.Intn(240) + 15),
				val: uuid.Must(uuid.NewV4()),
			},
			assert: func(t *testing.T, tc tcase, _ payment.Filter, err error) {
				testutil.AssertError(t, err, payment.ErrInvalidArgFilterCmpNotExists, payment.ErrMDArg("cmp", tc.args.cmp))
			},
		},
	}

	for i := range tcases {
		var tc = tcases[i]
		t.Run(tc.desc, func(t *testing.T) {
			var f, err = payment.NewFilterByID(tc.args.cmp, tc.args.val)
			tc.assert(t, tc, f, err)
		})
	}
}

func TestNewFilterByType(t *testing.T) {
	type params struct {
		cmp payment.FilterCmp
		val string
	}

	type tcase struct {
		desc   string
		args   params
		assert func(*testing.T, tcase, payment.Filter, error)
	}

	var tcases = []tcase{
		{
			desc: "successful",
			args: params{
				cmp: func() payment.FilterCmp {
					// nolint:gosec
					if rand.Int()%2 == 0 {
						return payment.FilterCmpEqual
					}

					return payment.FilterCmpNotEqual
				}(),
				val: "Payment",
			},
			assert: func(t *testing.T, tc tcase, f payment.Filter, err error) {
				assert.NoError(t, err)
				if assert.Equal(t, f.NodeType(), payment.FilterNodeTypeLeaf) {
					var l = f.Leaf()
					assert.True(t, l.IsSet())

					var cmp, val = l.Filter()
					assert.Equal(t, tc.args.cmp, cmp)
					assert.Equal(t, tc.args.val, val)
				}
			},
		},
		{
			desc: "error: unsupported cmp",
			args: params{
				cmp: payment.FilterCmpGreaterOrEqualThan,
				val: "Payment",
			},
			assert: func(t *testing.T, tc tcase, _ payment.Filter, err error) {
				testutil.AssertError(t, err, payment.ErrInvalidArgFilterCmpNotSupported, payment.ErrMDArg("cmp", tc.args.cmp))
			},
		},
		{
			desc: "error: cmp doesn't exist",
			args: params{
				cmp: payment.FilterCmp(rand.Intn(240) + 15),
				val: "Transfer",
			},
			assert: func(t *testing.T, tc tcase, _ payment.Filter, err error) {
				testutil.AssertError(t, err, payment.ErrInvalidArgFilterCmpNotExists, payment.ErrMDArg("cmp", tc.args.cmp))
			},
		},
		{
			desc: "error: invalid value",
			args: params{
				cmp: payment.FilterCmpNotEqual,
				val: "Transfer",
			},
			assert: func(t *testing.T, tc tcase, _ payment.Filter, err error) {
				testutil.AssertError(t, err, payment.ErrInvalidArgFilterValue, payment.ErrMDArg("val", tc.args.val))
			},
		},
	}

	for i := range tcases {
		var tc = tcases[i]
		t.Run(tc.desc, func(t *testing.T) {
			var f, err = payment.NewFilterByType(tc.args.cmp, tc.args.val)
			tc.assert(t, tc, f, err)
		})
	}
}
