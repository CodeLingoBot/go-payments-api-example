package payment

import (
	"math/big"

	"github.com/gofrs/uuid"
	"go.fraixed.es/errors"
)

// Pymt contains all the information which a single payment has.
type Pymt struct {
	PymtUpsert
	ID      uuid.UUID `json:"id"`
	Version uint32    `json:"version"`
}

// PymtUpsert contains the information required to create or update a payment.
type PymtUpsert struct {
	Type       string    `json:"type"`
	OrgID      uuid.UUID `json:"organisation_id"`
	Attributes Attrs     `json:"attributes"`
}

// Validate validates that the input payment contains all the required values
// and their values respect the requirments of the business domain.
func (p PymtUpsert) Validate() error {
	// TODO: won't be implemented
	// Some other fields should be also validated.

	if p.OrgID == uuid.Nil {
		return errors.New(ErrInvalidPaymentOrgID, ErrMDField("OrgID", p.OrgID))
	}

	if p.Type != "Payment" {
		return errors.New(ErrInvalidPaymentType, ErrMDField("Type", p.Type))
	}

	if err := p.Attributes.Validate(); err != nil {
		var c, _ = errors.GetCode(err)
		return errors.Wrap(err, c, ErrMDField("Attributes", p.Attributes))
	}

	return nil
}

// Attrs contains the information of the attributes attached to a payment.
type Attrs struct {
	Amount               *big.Float `json:"amount"`
	Currency             string     `json:"currency"`
	Reference            string     `json:"reference"`
	EndToEndReference    string     `json:"end_to_end_reference"`
	NumericReference     string     `json:"numeric_reference"`
	PaymentID            string     `json:"payment_id"`
	PaymentPurpose       string     `json:"payment_purpose"`
	PaymentScheme        string     `json:"payment_scheme"`
	PaymentType          string     `json:"payment_type"`
	ProcessingDate       string     `json:"processing_date"`
	SchemePaymentSubType string     `json:"scheme_payment_sub_type"`
	SchemePaymentType    string     `json:"scheme_payment_type"`
	BeneficiaryParty     Party      `json:"beneficiary_party"`
	DebtorParty          Party      `json:"debtor_party"`
	SponsorParty         Party      `json:"sponsor_party"`
	ChargesInformation   struct {
		BearerCode    string `json:"bearer_code"`
		SenderCharges []struct {
			Amount   *big.Float `json:"amount"`
			Currency string     `json:"currency"`
		} `json:"sender_charges"`
		ReceiverChargesAmount   *big.Float `json:"receiver_charges_amount"`
		ReceiverChargesCurrency string     `json:"receiver_charges_currency"`
	} `json:"charges_information"`
	Fx struct {
		ContractReference string `json:"contract_reference"`
		ExchangeRate      string `json:"exchange_rate"`
		OriginalAmount    string `json:"original_amount"`
		OriginalCurrency  string `json:"original_currency"`
	} `json:"fx"`
}

// Validate validates that the attributes contains alls the required values and
// their values respect the requirements of the business domain.
func (a Attrs) Validate() error {
	// TODO: won't be implemented
	// Some other fields should be also validated.

	if a.PaymentID == "" {
		return errors.New(ErrInvalidPaymentAttrPaymentID, ErrMDField("PaymentID", ""))
	}

	return nil
}

// Party contains the information of each single party involved in a payment.
type Party struct {
	AccountName       string `json:"account_name"`
	AccountNumber     string `json:"account_number"`
	AccountNumberCode string `json:"account_number_code"`
	AccountType       int    `json:"account_type"`
	Address           string `json:"address"`
	BankID            string `json:"bank_id"`
	BankIDCode        string `json:"bank_id_code"`
	Name              string `json:"name"`
}

// Validate valides that the party contains all the required values and their
// values respect the requirments of the business domain.
func (p Party) Validate() error {
	// TODO: won't be implemented
	return nil
}
