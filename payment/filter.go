package payment

import (
	"github.com/gofrs/uuid"
	"go.fraixed.es/errors"
)

// FilterCmp specifies the comparison operation to perform.
type FilterCmp uint8

// The list of valid FilterCmp values.
const (
	filterCmpNone FilterCmp = iota
	FilterCmpEqual
	FilterCmpNotEqual
	FilterCmpGreaterThan
	FilterCmpGreaterOrEqualThan
	FilterCmpLessThan
	FilterCmpLessOrEqualThan
	FilterCmpMatch
)

// FilterLogical specifies the logical operation to perform when more than one
// comparison operation is involved.
type FilterLogical uint8

// The list of valid FilterLogical values.
const (
	filterLogicalNoop FilterLogical = iota
	FilterLogicalAnd
	FilterLogicalOr
)

// FilterNodeType represents the type that a filter node can have.
type FilterNodeType uint8

// The list of valid FilterNodeType values.
const (
	FilterNodeTypeEmpty FilterNodeType = iota
	FilterNodeTypeLeaf
	FilterNodeTypeNonLeaf
)

// FilterLeaf is the interface which any type which can be a leaf node of a
// Filter must satisifies
type FilterLeaf interface {
	Filter() (op FilterCmp, val interface{})
	IsSet() bool
}

// Filter is an abstraction of a specified filter over a set of different values.
//
// It represents a binary tree where when it's a leaf node, it contains the
// comparison operator and value to apply, meanwhile when it isn't a left, it
// contains the 2 children nodes (which are Filter) and the logical operator
// applied between the both of them.
//
// You can see it as any boolean logical operation used in conditionals in the
// majority of programming languages.
//
// The zero value is a noop filter.
type Filter struct {
	left  *Filter
	right *Filter
	l     FilterLeaf
	op    FilterLogical
}

// NodeType returns the type of this node. Zero value returns
// FilterNodeTypeEmpty.
func (f Filter) NodeType() FilterNodeType {
	switch {
	case f.op == filterLogicalNoop && f.l == nil:
		return FilterNodeTypeEmpty
	case f.op != filterLogicalNoop:
		return FilterNodeTypeNonLeaf
	default:
		return FilterNodeTypeLeaf
	}
}

// Leaf returns the FilterLeaf value. The value is nil when it's node type is not
// NodeTypeLeaf.
func (f Filter) Leaf() FilterLeaf {
	return f.l
}

// Nodes returns the operator applied over the 2 children nodes and those nodes
// left and right.
//
// FilterLogical is an invalid value and the 2 filter will be of node type
// NodeTypeEmpty, when f is not an NodeTypeNonLeaf.
func (f Filter) Nodes() (FilterLogical, Filter, Filter) {
	if f.NodeType() == FilterNodeTypeNonLeaf {
		return f.op, *f.left, *f.right
	}

	return filterLogicalNoop, Filter{}, Filter{}
}

// NewFilter creates a filter of NodeTypeNonLeaf composed by 2 other Filters.
func NewFilter(op FilterLogical, left Filter, right Filter) (Filter, error) {
	if err := validateLogicalOp(op); err != nil {
		return Filter{}, err
	}

	if left.NodeType() == FilterNodeTypeEmpty {
		return Filter{}, errors.New(ErrInvalidArgFilterNodeEmpty, ErrMDArg("left", left))
	}

	if right.NodeType() == FilterNodeTypeEmpty {
		return Filter{}, errors.New(ErrInvalidArgFilterNodeEmpty, ErrMDArg("right", right))
	}

	return Filter{
		op:    op,
		left:  &left,
		right: &right,
	}, nil
}

// FilterLeafID is the FilterLeaf for filtering payments by ID.
type FilterLeafID struct {
	val uuid.UUID
	cmp FilterCmp
}

// NewFilterByID creates a new Filer leaf node of a FilterLeafID with the
// specified cmp and val.
//
// The following error codes can be returned (declared in errs sub package):
//
// * InvalidArgFilterCmpNotExists
//
// * InvalidArgFilterCmpNotSupported - when cmp isn't FilterCmpEqual nor
// FilterCmpNotEqual
func NewFilterByID(cmp FilterCmp, val uuid.UUID) (Filter, error) {
	if err := validatepCmp(cmp); err != nil {
		return Filter{}, err
	}

	switch cmp {
	case FilterCmpEqual, FilterCmpNotEqual:
	default:
		return Filter{}, errors.New(ErrInvalidArgFilterCmpNotSupported, ErrMDArg("cmp", cmp))
	}

	return newFilterLeaf(FilterLeafID{
		val: val,
		cmp: cmp,
	})
}

// Filter returns the operation and ID value which has been set.
func (f FilterLeafID) Filter() (FilterCmp, interface{}) {
	return f.cmp, f.val
}

// IsSet returns true when the filter is set, otherwise none.
func (f FilterLeafID) IsSet() bool {
	return f.cmp != filterCmpNone
}

// FilterLeafType is the FilterLeaf for filtering payments by type.
type FilterLeafType struct {
	filterLeafString
}

// NewFilterByType creates a new Filter leaf node of a FilterByType with cmp and
// val.
//
// The following error codes can be returned (declared in errs sub package):
//
// * InvalidArgFilterCmpNotExists
//
// * InvalidArgFilterCmpNotSupported - when cmp isn't FilterCmpEqual nor
// FilterCmpNotEqual
//
// * InvalidArgFilterValue - currently only "Payment" is valid value type
func NewFilterByType(cmp FilterCmp, val string) (Filter, error) {
	if err := validatepCmp(cmp); err != nil {
		return Filter{}, err
	}

	switch cmp {
	case FilterCmpEqual, FilterCmpNotEqual:
	default:
		return Filter{}, errors.New(ErrInvalidArgFilterCmpNotSupported, ErrMDArg("cmp", cmp))
	}

	if val != "Payment" {
		return Filter{}, errors.New(ErrInvalidArgFilterValue, ErrMDArg("val", val))
	}

	var f, err = newFilterLeafString(cmp, val)
	if err != nil {
		return Filter{}, err
	}

	return newFilterLeaf(FilterLeafType{f})
}

// FilterLeafAmount allows to filter payment by its amount filed.
type FilterLeafAmount struct {
	filterLeafFloat64
}

// NewFilterByAmount creates a new Filter leaf node of a FilterByAmount with the
// specified cmp and val.
//
// The following error codes can be returned (declared in errs sub package):
//
// * InvalidArgFilterCmpNotExists
//
// * InvalidArgFilterCmpNotSupported - when cmd is FilterCmpMatch
func NewFilterByAmount(cmp FilterCmp, val float64) (Filter, error) {
	var f, err = newFilterLeafFloat64(cmp, val)
	if err != nil {
		return Filter{}, err
	}

	return newFilterLeaf(FilterLeafAmount{f})
}

// newFilterLeaf creates a Filter of NodeTypeLeaf.
// It returns an error if f.IsSet returns false.
//
// The following error codes may be returned:
//
// * InvalidArgFilterLeafNoValSet
func newFilterLeaf(l FilterLeaf) (Filter, error) {
	if !l.IsSet() {
		return Filter{}, errors.New(ErrInvalidArgFilterLeafNoValSet)
	}

	return Filter{
		l: l,
	}, nil
}

// filterLeafString is a leaf node filter type for being used for filtering string
// values.
type filterLeafString struct {
	val string
	cmp FilterCmp
}

// newFilterLeafString creates a new filterLeafString with the specified cmp and
// val.
// It returns an error if cmp isn't one of the list of accepted FilterCmp values.
func newFilterLeafString(cmp FilterCmp, val string) (filterLeafString, error) {
	if err := validatepCmp(cmp); err != nil {
		return filterLeafString{}, err
	}

	return filterLeafString{
		val: val,
		cmp: cmp,
	}, nil
}

// Filter returns the operation and string value which has been set.
func (f filterLeafString) Filter() (FilterCmp, interface{}) {
	return f.cmp, f.val
}

// IsSet returns true when the filter is set, otherwise none.
func (f filterLeafString) IsSet() bool {
	return f.cmp != filterCmpNone
}

// filterLeafFloat64 is the filtering type for being used for the float64 type.
type filterLeafFloat64 struct {
	val float64
	cmp FilterCmp
}

// newFilterLeafFloat64 creates a new filterLeafFloat64 with the specified cmp
// and val.
// It returns an error if cmp isn't one of the list of accepted FilterCmp values.
func newFilterLeafFloat64(cmp FilterCmp, val float64) (filterLeafFloat64, error) {
	if err := validatepCmp(cmp); err != nil {
		return filterLeafFloat64{}, err
	}

	if cmp == FilterCmpMatch {
		return filterLeafFloat64{}, errors.New(ErrInvalidArgFilterCmpNotSupported, ErrMDArg("cmp", cmp))
	}

	return filterLeafFloat64{
		val: val,
		cmp: cmp,
	}, nil
}

// Filter returns the operation and string value which has been set.
func (f filterLeafFloat64) Filter() (FilterCmp, interface{}) {
	return f.cmp, f.val
}

// IsSet returns true when the filter is set, otherwise none.
func (f filterLeafFloat64) IsSet() bool {
	return f.cmp != filterCmpNone
}

func validateLogicalOp(op FilterLogical) error {
	switch op {
	case FilterLogicalAnd, FilterLogicalOr:
		return nil
	}

	return errors.New(ErrInvalidArgFilterLogicalOpNotExists, ErrMDArg("op", op))
}

func validatepCmp(cmp FilterCmp) error {
	switch cmp {
	case FilterCmpEqual, FilterCmpGreaterOrEqualThan, FilterCmpGreaterThan,
		FilterCmpLessOrEqualThan, FilterCmpLessThan, FilterCmpMatch, FilterCmpNotEqual:
		return nil
	}

	return errors.New(ErrInvalidArgFilterCmpNotExists, ErrMDArg("cmp", cmp))
}
