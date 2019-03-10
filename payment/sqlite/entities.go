package sqlite

import (
	"encoding/json"

	"github.com/ifraixedes/go-payments-api-example/payment"
	"go.fraixed.es/errors"
)

type pymtAttrs struct {
	Amount               float64   `json:"amount"`
	Currency             string    `json:"currency"`
	Reference            string    `json:"reference"`
	EndToEndReference    string    `json:"end_to_end_reference"`
	NumericReference     string    `json:"numeric_reference"`
	PaymentID            string    `json:"payment_id"`
	PaymentPurpose       string    `json:"payment_purpose"`
	PaymentScheme        string    `json:"payment_scheme"`
	PaymentType          string    `json:"payment_type"`
	ProcessingDate       string    `json:"processing_date"`
	SchemePaymentSubType string    `json:"scheme_payment_sub_type"`
	SchemePaymentType    string    `json:"scheme_payment_type"`
	BeneficiaryParty     pymtParty `json:"beneficiary_party"`
	DebtorParty          pymtParty `json:"debtor_party"`
	SponsorParty         pymtParty `json:"sponsor_party"`
	ChargesInformation   struct {
		BearerCode    string `json:"bearer_code"`
		SenderCharges []struct {
			Amount   float64 `json:"amount"`
			Currency string  `json:"currency"`
		} `json:"sender_charges"`
		ReceiverChargesAmount   float64 `json:"receiver_charges_amount"`
		ReceiverChargesCurrency string  `json:"receiver_charges_currency"`
	} `json:"charges_information"`
	Fx struct {
		ContractReference string `json:"contract_reference"`
		ExchangeRate      string `json:"exchange_rate"`
		OriginalAmount    string `json:"original_amount"`
		OriginalCurrency  string `json:"original_currency"`
	} `json:"fx"`
}

type pymtParty struct {
	AccountName       string `json:"account_name"`
	AccountNumber     string `json:"account_number"`
	AccountNumberCode string `json:"account_number_code"`
	AccountType       int    `json:"account_type"`
	Address           string `json:"address"`
	BankID            string `json:"bank_id"`
	BankIDCode        string `json:"bank_id_code"`
	Name              string `json:"name"`
}

// pymtData groups several payment fields which are stored as blob in the DB.
type pymtData struct {
	Type string `json:"type"`
	pymtAttrs
}

// Serialize returns blob to store pd in the DB.
func (pd *pymtData) Serialize() ([]byte, error) {
	b, err := json.Marshal(pd)
	if err != nil {
		return nil, errors.Wrap(err, payment.ErrUnexpectedSysError)
	}

	return b, nil
}

// Deserialize initializes pd from b. b is usually the bob stored in the DB.
func (pd *pymtData) Deserialize(b []byte) error {
	err := json.Unmarshal(b, pd)
	if err != nil {
		return errors.Wrap(err, ErrInvalidFormatBlob)
	}

	return nil
}

// Init initializes the pd from p.
func (pd *pymtData) Init(p payment.PymtUpsert) {
	pd.Type = p.Type
	//pd.Attrs = p.Attributes
	// TODO: use reflection
}

// Set sets p with the values hold by pd. If p is nill, it panics.
func (pd *pymtData) Set(p *payment.PymtUpsert) {
	p.Type = pd.Type
	// p.Attributes = pd.Attrs
	// TODO: use reflection
}

// pymt is how a payment is stored in the DB.
type pymt struct {
	ID      string
	Version uint32
	OrgID   string
	Data    []byte
}
