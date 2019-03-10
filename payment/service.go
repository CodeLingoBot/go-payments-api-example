package payment

import (
	"context"

	"github.com/gofrs/uuid"
)

// Service is the interface which any specific implementation of a payment
// service must satisfy.
// All the implementation must also return the error codes documented in this
// interface (general ones and on each method).
//
// All the methods can return, a part of their specific ones which are
// documented on them, the following error codes:
//
// * ErrOperationAborted - This could happen due to a ctx deadline; ctx.Err must
//	   be checked to find out that the error is returned because of ctx.
//
// * ErrUnexpectedOSError
//
// * ErrUnexpectedStoreError
//
// * ErrUnexpectedSysError
type Service interface {
	// Create creates a new payment returning its ID.
	//
	// This method can return any of the errors returned by p.Validate.
	Create(ctx context.Context, p PymtUpsert) (uuid.UUID, error)

	// Delete deletes the payment which has associated the passed ID.
	//
	// The following error codes can be returned:
	//
	// * ErrInvalidPaymentID
	//
	// * ErrNotFound
	Delete(context.Context, uuid.UUID) error

	// Find retrieve list of payments which fulfill f, sorted by o and chunked by
	// c. Each payment only contains the fields indicated by s.
	// If there is not payments which fulfill f or the chunk specified by p is out
	// of range, an empty list and nil error are returned.
	Find(ctx context.Context, f Filter, s Selection, o Sort, c Chunk) ([]Pymt, error)

	// Get retrieves the payment which has associated the passed ID. The payment
	// only contains the fields indicated by s.
	//
	// The following error codes can be returned:
	//
	// * ErrInvalidPaymentID
	//
	// * ErrNotFound
	Get(ctx context.Context, id uuid.UUID, s Selection) (Pymt, error)

	// Update updates the payment with the associated ID, if its version matches
	// with version. The payment version is incremented if the update succeeds.
	//
	// The following error codes can be returned:
	//
	// * ErrInvalidArgVersionMismatch - When the version doesn't match with the
	// current payment version for avoiding to override the payment concurrently.
	//
	// * ErrNotFound
	Update(ctx context.Context, id uuid.UUID, version uint32, p PymtUpsert) error
}
