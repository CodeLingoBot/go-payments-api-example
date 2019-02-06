package payment_test

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/ifraixedes/go-payments-api-example/payment"
	"github.com/ifraixedes/go-payments-api-example/payment/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPymtUpsert_Validate(t *testing.T) {
	var tcases = []struct {
		desc   string
		p      payment.PymtUpsert
		assert func(*testing.T, payment.PymtUpsert, error)
	}{
		{
			desc: "Valid",
			p: payment.PymtUpsert{
				OrgID: func() uuid.UUID {
					id, err := uuid.NewV4()
					require.NoError(t, err)

					return id
				}(),
				Type: "Payment",
				Attributes: payment.Attrs{
					PaymentID: "some-id",
				},
			},
			assert: func(t *testing.T, _ payment.PymtUpsert, err error) {
				assert.NoError(t, err)
			},
		},
		{
			desc: "Invalid: OrgID",
			p: payment.PymtUpsert{
				Type: "Payment",
				Attributes: payment.Attrs{
					PaymentID: "some-id",
				},
			},
			assert: func(t *testing.T, p payment.PymtUpsert, err error) {
				testutil.AssertError(t, err, payment.ErrInvalidPaymentOrgID, payment.ErrMDField("OrgID", p.OrgID))
			},
		},
		{
			desc: "Invalid: Type",
			p: payment.PymtUpsert{
				OrgID: func() uuid.UUID {
					id, err := uuid.NewV4()
					require.NoError(t, err)

					return id
				}(),
				Type: "Transfer",
				Attributes: payment.Attrs{
					PaymentID: "some-id",
				},
			},
			assert: func(t *testing.T, p payment.PymtUpsert, err error) {
				testutil.AssertError(t, err, payment.ErrInvalidPaymentType, payment.ErrMDField("Type", p.Type))
			},
		},
		{
			desc: "Invalid: Attributes",
			p: payment.PymtUpsert{
				OrgID: func() uuid.UUID {
					id, err := uuid.NewV4()
					require.NoError(t, err)

					return id
				}(),
				Type: "Payment",
			},
			assert: func(t *testing.T, p payment.PymtUpsert, err error) {
				testutil.AssertError(t, err, payment.ErrInvalidPaymentAttrPaymentID, payment.ErrMDField("Attributes", p.Attributes))
			},
		},
	}

	for i := range tcases {
		var tc = tcases[i]
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			tc.assert(t, tc.p, tc.p.Validate())
		})
	}
}

func TestAttrs_Validate(t *testing.T) {
	var tcases = []struct {
		desc   string
		a      payment.Attrs
		assert func(*testing.T, payment.Attrs, error)
	}{
		{
			desc: "Valid",
			a: payment.Attrs{
				PaymentID: "some-id",
			},
			assert: func(t *testing.T, _ payment.Attrs, err error) {
				assert.NoError(t, err)
			},
		},
		{
			desc: "Invalid: PaymentID",
			a:    payment.Attrs{},
			assert: func(t *testing.T, a payment.Attrs, err error) {
				testutil.AssertError(t, err, payment.ErrInvalidPaymentAttrPaymentID, payment.ErrMDField("PaymentID", a.PaymentID))
			},
		},
	}

	for i := range tcases {
		var tc = tcases[i]
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			tc.assert(t, tc.a, tc.a.Validate())
		})
	}
}
