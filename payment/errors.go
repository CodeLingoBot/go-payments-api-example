package payment

import (
	"fmt"

	"go.fraixed.es/errors"
)

type code uint8

// The list of error codes that the payment service can return.
const (
	ErrAbortedOperation code = iota + 1

	ErrInvalidArgFilterCmpNotExists
	ErrInvalidArgFilterCmpNotSupported
	ErrInvalidArgFilterLeafNoValSet
	ErrInvalidArgFilterLogicalOpNotExists
	ErrInvalidArgFilterNodeEmpty
	ErrInvalidArgFilterValue

	ErrInvalidArgVersionMismatch

	ErrInvalidPaymentID
	ErrInvalidPaymentOrgID
	ErrInvalidPaymentType
	ErrInvalidPaymentAttrPaymentID

	ErrNotFound

	ErrUnexpectedOSError
	ErrUnexpectedStoreError
	ErrUnexpectedSysError
)

func (c code) String() string {
	switch c {
	case ErrAbortedOperation:
		return "AbortedOperation"
	case ErrInvalidArgFilterCmpNotExists:
		return "InvalidArgFilterCmpNotExists"
	case ErrInvalidArgFilterCmpNotSupported:
		return "InvalidArgFilterCmpNotSupported"
	case ErrInvalidArgFilterLeafNoValSet:
		return "InvalidArgFilterLeafNoValSet"
	case ErrInvalidArgFilterLogicalOpNotExists:
		return "InvalidArgFilterLogicalOp"
	case ErrInvalidArgFilterNodeEmpty:
		return "InvalidArgFilterNodeEmpty"
	case ErrInvalidArgFilterValue:
		return "InvalidArgFilterValue"
	case ErrInvalidArgVersionMismatch:
		return "InvalidArgVersionMismatch"
	case ErrInvalidPaymentID:
		return "InvalidPaymentID"
	case ErrInvalidPaymentOrgID:
		return "InvalidPaymentOgID"
	case ErrInvalidPaymentType:
		return "InvalidPaymentType"
	case ErrInvalidPaymentAttrPaymentID:
		return "InvalidPaymentAttrPaymentID"
	case ErrNotFound:
		return "NotFound"
	case ErrUnexpectedOSError:
		return "UnexpectedOSError"
	case ErrUnexpectedStoreError:
		return "UnexpectedStoreError"
	case ErrUnexpectedSysError:
		return "UnexpectedSysError"
	}

	return ""
}

func (c code) Message() string {
	switch c {
	case ErrAbortedOperation:
		return "the operation has been aborted"
	case ErrInvalidArgFilterCmpNotExists:
		return "The filter comparison operator doesn't exist"
	case ErrInvalidArgFilterCmpNotSupported:
		return "The filter comparison operator isn't supported for this type leaf node"
	case ErrInvalidArgFilterLeafNoValSet:
		return "The filter leaf node is invalid because its value isn't set"
	case ErrInvalidArgFilterLogicalOpNotExists:
		return "The filter logical operator doesn't exist"
	case ErrInvalidArgFilterNodeEmpty:
		return "The filter node cannot be an empty"
	case ErrInvalidArgFilterValue:
		return "The filter value isn't a valid one"
	case ErrInvalidArgVersionMismatch:
		return "The provided version doesn't match with the current one"
	case ErrInvalidPaymentID:
		return "invalid payment because its ID is not valid"
	case ErrInvalidPaymentOrgID:
		return "Invalid payment because its organisation ID is not valid"
	case ErrInvalidPaymentType:
		return "Invalid payment because its type value is not valid"
	case ErrInvalidPaymentAttrPaymentID:
		return "Invalid payment because the payment ID value of its attributes is not valid"
	case ErrNotFound:
		return "The entity was not found"
	case ErrUnexpectedOSError:
		return "an unexpected error has been returned when performing an operative system operation"
	case ErrUnexpectedStoreError:
		return "the store has returned an unexpected error"
	case ErrUnexpectedSysError:
		return "an unexpected general system error has happened"
	}

	return ""
}

// ErrMDArg creates a new metadata from an function argument which is related
// with the error to create.
func ErrMDArg(name string, val interface{}) errors.MD {
	return errors.MD{
		K: fmt.Sprintf("arg:%s", name),
		V: val,
	}
}

// ErrMDVar creates a new metdata from a variable whose name and value are
// relevant for the error to create.
// When the variable is not exposed, the name should be meaningful to the
// user/developer/ops when reading the verbose version of the error.
func ErrMDVar(name string, val interface{}) errors.MD {
	return errors.MD{
		K: fmt.Sprintf("var:%s", name),
		V: val,
	}
}

// ErrMDFnCall creates a new metadata to inform the function which has been
// called and with which arguments.
// This metadata is intended to be used for internal function calls, so the
// user isn't aware of those and when they return an error code which isn't
// enough concrete to let inform the user what specifically happened.
func ErrMDFnCall(fname string, args ...interface{}) errors.MD {
	return errors.MD{
		K: fmt.Sprintf("func:%s", fname),
		V: fmt.Sprintf("%+v", args),
	}
}

// ErrMDField creates a new metadata from a struct field whose name and value
// are relevant for the error to create.
func ErrMDField(name string, val interface{}) errors.MD {
	return errors.MD{
		K: fmt.Sprintf("field:%s", name),
		V: val,
	}
}
