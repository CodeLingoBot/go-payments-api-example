package payment

// SortDir is the type which represents the direction when sorting values.
type SortDir uint8

// The list of valid SortDir values
const (
	SortUnspecified SortDir = iota
	SortAscending
	SortDescending
)

// Sort specifies the allowed fields for sorting a list of payments.
//
// The payment's fields which aren't present aren't allowed to be used for
// sorting.
type Sort struct {
	Type       SortDir
	ID         SortDir
	Version    SortDir
	OrgID      SortDir
	Attributes SortAttributes
}

// SortAttributes is the type of the Attributes field of the Sort type.
type SortAttributes struct {
	Amount             SortDir
	Currency           SortDir
	Reference          SortDir
	EndToEndReference  SortDir
	NumericReference   SortDir
	PaymentID          SortDir
	ProcessingDate     SortDir
	ChargesInformation SortChargesInformation
}

// SortChargesInformation is the type of the ChargesInformation field of the
// SortAttributes type.
type SortChargesInformation struct {
	SenderChargesAmount   SortDir
	ReceiverChargesAmount SortDir
}
