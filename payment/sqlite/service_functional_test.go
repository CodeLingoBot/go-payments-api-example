package sqlite_test

import (
	"context"
	"math/big"
	"math/rand"
	"testing"

	"github.com/bxcodec/faker"
	"github.com/gofrs/uuid"
	"github.com/ifraixedes/go-payments-api-example/payment"
	"github.com/ifraixedes/go-payments-api-example/payment/internal/testutil"
	"github.com/ifraixedes/go-payments-api-example/payment/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_Create_Get_Delete(t *testing.T) {
	var svc, err = sqlite.New(testingDB)
	require.NoError(t, err)

	var ctx = context.Background()

	// Create payment
	var npymt = payment.PymtUpsert{
		Type:  "Payment",
		OrgID: testutil.NewUUID(t),
	}
	faker.FakeData(&npymt.Attributes)
	npymt.Attributes.Amount = big.NewFloat(rand.Float64())

	pid, err := svc.Create(ctx, npymt)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, pid)

	// Get payment with all the fields
	pymt, err := svc.Get(ctx, pid, payment.SelectAll())
	require.NoError(t, err)
	assert.Equal(t, payment.Pymt{ID: pid, PymtUpsert: npymt}, pymt)

	// Get payment with only a few fields
	pymt, err = svc.Get(ctx, pid, payment.Selection{
		Type: true,
		Attributes: payment.SelectionAttributes{
			Amount: true,
		},
	})
	require.NoError(t, err)
	assert.Equal(t, payment.Pymt{
		ID: pid,
		PymtUpsert: payment.PymtUpsert{
			Type: npymt.Type,
			Attributes: payment.Attrs{
				Amount: npymt.Attributes.Amount,
			},
		},
	}, pymt)

	// Delete payment
	err = svc.Delete(ctx, pid)
	require.NoError(t, err)

	_, err = svc.Get(ctx, pid, payment.SelectAll())
	testutil.AssertError(t, err, payment.ErrNotFound, payment.ErrMDArg("id", pid))
}

func TestService_Create_Get_Update_Delete(t *testing.T) {
	var svc, err = sqlite.New(testingDB)
	require.NoError(t, err)

	var ctx = context.Background()

	// Create payment
	var npymt = payment.PymtUpsert{
		Type:  "Payment",
		OrgID: testutil.NewUUID(t),
	}
	faker.FakeData(&npymt.Attributes)
	npymt.Attributes.Amount = big.NewFloat(rand.Float64())

	pid, err := svc.Create(ctx, npymt)
	require.NoError(t, err)

	// Get payment
	pymt, err := svc.Get(ctx, pid, payment.SelectAll())
	require.NoError(t, err)

	var upymt = payment.PymtUpsert{
		Type:  pymt.Type,
		OrgID: pymt.OrgID,
	}
	faker.FakeData(&upymt.Attributes)
	upymt.Attributes.Amount = pymt.Attributes.Amount

	// Update with the wrong version number
	var iv = uint32(rand.Int31()) + pymt.Version
	err = svc.Update(ctx, pid, iv, upymt)
	testutil.AssertError(t, err, payment.ErrInvalidArgVersionMismatch, payment.ErrMDArg("version", iv))

	// Update with the right version
	err = svc.Update(ctx, pid, pymt.Version, upymt)
	assert.NoError(t, err)

	pymtu, err := svc.Get(ctx, pid, payment.SelectAll())
	require.NoError(t, err)

	assert.Equal(t, pymt.Version+1, pymtu.Version)
	assert.Equal(t, pymtu.PymtUpsert, upymt)

	// Delete payment
	err = svc.Delete(ctx, pid)
	require.NoError(t, err)

	// Update unexisting payment
	err = svc.Update(ctx, pid, pymt.Version, upymt)
	testutil.AssertError(t, err, payment.ErrNotFound, payment.ErrMDArg("id", pid))
}

func TestService_Create_Find_Delete(t *testing.T) {
	var svc, err = sqlite.New(testingDB)
	require.NoError(t, err)

	var (
		ctx = context.Background()
		a1  = rand.Float64()
		a2  = a1 + 10
		a3  = a2 + 10
	)

	var npymt = payment.PymtUpsert{
		Type:  "Payment",
		OrgID: testutil.NewUUID(t),
	}
	faker.FakeData(&npymt.Attributes)
	npymt.Attributes.Amount = big.NewFloat(a1)

	id1, err := svc.Create(ctx, npymt)
	require.NoError(t, err)
	p1, err := svc.Get(ctx, id1, payment.SelectAll())
	require.NoError(t, err)

	npymt = payment.PymtUpsert{
		Type:  "Payment",
		OrgID: testutil.NewUUID(t),
	}
	faker.FakeData(&npymt.Attributes)
	npymt.Attributes.Amount = big.NewFloat(a2)

	id2, err := svc.Create(ctx, npymt)
	require.NoError(t, err)
	p2, err := svc.Get(ctx, id2, payment.SelectAll())
	require.NoError(t, err)

	npymt = payment.PymtUpsert{
		Type:  "Payment",
		OrgID: testutil.NewUUID(t),
	}
	faker.FakeData(&npymt.Attributes)
	npymt.Attributes.Amount = big.NewFloat(a3)

	id3, err := svc.Create(ctx, npymt)
	require.NoError(t, err)
	p3, err := svc.Get(ctx, id3, payment.SelectAll())
	require.NoError(t, err)

	defer func() {
		_ = svc.Delete(ctx, id1)
		_ = svc.Delete(ctx, id2)
		_ = svc.Delete(ctx, id3)
	}()

	t.Run("find all", func(t *testing.T) {
		var pms, err = svc.Find(
			ctx, payment.Filter{}, payment.SelectAll(),
			payment.Sort{
				Attributes: payment.SortAttributes{
					Amount: payment.SortAscending,
				},
			}, payment.Chunk{},
		)
		require.NoError(t, err)

		assert.Len(t, pms, 3)
		assert.Equal(t, []payment.Pymt{p1, p2, p3}, pms)
	})

	t.Run("find some", func(t *testing.T) {
		var ft, err = payment.NewFilterByAmount(payment.FilterCmpGreaterOrEqualThan, a2)
		require.NoError(t, err)

		pms, err := svc.Find(
			ctx, ft, payment.SelectAll(),
			payment.Sort{
				Attributes: payment.SortAttributes{
					Amount: payment.SortDescending,
				},
			}, payment.Chunk{},
		)
		require.NoError(t, err)

		assert.Len(t, pms, 2)
		assert.Equal(t, []payment.Pymt{p3, p2}, pms)
	})

	t.Run("find some with some fields", func(t *testing.T) {
		var ft, err = payment.NewFilterByAmount(payment.FilterCmpGreaterOrEqualThan, a2)
		require.NoError(t, err)

		pms, err := svc.Find(
			ctx, ft, payment.Selection{OrgID: true},
			payment.Sort{
				Attributes: payment.SortAttributes{
					Amount: payment.SortDescending,
				},
			}, payment.Chunk{},
		)
		require.NoError(t, err)

		assert.Len(t, pms, 2)

		var (
			exp2 = payment.Pymt{
				ID: p2.ID,
				PymtUpsert: payment.PymtUpsert{
					OrgID: p2.OrgID,
				},
			}
			exp3 = payment.Pymt{
				ID: p3.ID,
				PymtUpsert: payment.PymtUpsert{
					OrgID: p3.OrgID,
				},
			}
		)
		assert.Equal(t, []payment.Pymt{exp3, exp2}, pms)
	})

	t.Run("find a chunk", func(t *testing.T) {
		var ftl, err = payment.NewFilterByID(payment.FilterCmpEqual, id1)
		require.NoError(t, err)

		ftr, err := payment.NewFilterByID(payment.FilterCmpEqual, id2)
		require.NoError(t, err)

		ft, err := payment.NewFilter(payment.FilterLogicalOr, ftl, ftr)
		require.NoError(t, err)

		pms, err := svc.Find(
			ctx, ft, payment.SelectAll(),
			payment.Sort{
				Attributes: payment.SortAttributes{
					Amount: payment.SortAscending,
				},
			}, payment.Chunk{
				Limit: 1,
			},
		)
		require.NoError(t, err)

		assert.Len(t, pms, 1)
		assert.Equal(t, []payment.Pymt{p1}, pms)

		pms, err = svc.Find(
			ctx, ft, payment.SelectAll(),
			payment.Sort{
				Attributes: payment.SortAttributes{
					Amount: payment.SortAscending,
				},
			}, payment.Chunk{
				Limit:  1,
				Offset: 1,
			},
		)
		require.NoError(t, err)

		assert.Len(t, pms, 1)
		assert.Equal(t, []payment.Pymt{p2}, pms)
	})

	t.Run("find any", func(t *testing.T) {
		var ft, err = payment.NewFilterByID(payment.FilterCmpEqual, testutil.NewUUID(t))
		require.NoError(t, err)

		pms, err := svc.Find(
			ctx, ft, payment.SelectAll(),
			payment.Sort{
				Attributes: payment.SortAttributes{
					Amount: payment.SortDescending,
				},
			}, payment.Chunk{},
		)
		require.NoError(t, err)

		assert.Len(t, pms, 0)
	})
}
