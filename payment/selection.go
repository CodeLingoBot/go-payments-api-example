package payment

// Selection specifies the fields of a single payment to be retrieved.
// The payment's fields which aren't present are always retrieved.
// Each file is a boolean, when it's true, the value is retrieved otherwise it
// won't be.
type Selection struct {
	Version    bool
	Type       bool
	OrgID      bool
	Attributes SelectionAttributes
}

// SelectionAttributes is the type of the Attributes field of the Selection type.
type SelectionAttributes struct {
	Amount               bool
	Currency             bool
	Reference            bool
	EndToEndReference    bool
	NumericReference     bool
	PaymentID            bool
	PaymentPurpose       bool
	PaymentScheme        bool
	PaymentType          bool
	ProcessingDate       bool
	SchemePaymentSubType bool
	SchemePaymentType    bool
	BeneficiaryParty     bool
	DebtorParty          bool
	SponsorParty         bool
	ChargesInformation   bool
	Fx                   bool
}
