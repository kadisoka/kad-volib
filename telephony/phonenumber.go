package telephony

import (
	"strings"

	"github.com/alloyzeus/go-azfl/v2/azcore"
	"github.com/nyaruka/phonenumbers"
)

// PhoneNumber represents a phone number as we need.
//
// Note that phonenumber is pretty complex. We can see it
// in "github.com/nyaruka/phonenumbers".PhoneNumber. We will
// add more details as needed.
type PhoneNumber struct {
	countryCode    int32
	nationalNumber int64
	rawInput       string
	isValid        bool
}

var _ azcore.ValueObject = PhoneNumber{}
var _ azcore.ValueObjectAssert[PhoneNumber] = PhoneNumber{}

func NewPhoneNumber(countryCode int32, nationalNumber int64) PhoneNumber {
	if nationalNumber < 0 {
		panic("national number must be positive number")
	}
	nn := uint64(nationalNumber)
	parsed := &phonenumbers.PhoneNumber{
		CountryCode:    &countryCode,
		NationalNumber: &nn,
	}
	return PhoneNumber{
		countryCode:    *parsed.CountryCode,
		nationalNumber: int64(*parsed.NationalNumber),
		rawInput:       "",
		isValid:        phonenumbers.IsValidNumber(parsed),
	}
}

func PhoneNumberFromString(input string) (PhoneNumber, error) {
	// Check if the country code is doubled
	if parts := strings.Split(input, "+"); len(parts) == 3 {
		// We assume that the first part was automatically added at the client
		input = "+" + parts[2]
	}

	parsed, err := phonenumbers.Parse(input, "")
	if err != nil {
		return PhoneNumber{}, err
	}

	pn := PhoneNumber{
		countryCode:    *parsed.CountryCode,
		nationalNumber: int64(*parsed.NationalNumber),
		rawInput:       input,
		isValid:        phonenumbers.IsValidNumber(parsed),
	}

	return pn, nil
}

func (pn PhoneNumber) IsStaticallyValid() bool { return pn.isValid }

func (pn PhoneNumber) Clone() PhoneNumber { return pn }

func (pn PhoneNumber) Equal(other interface{}) bool {
	return pn.Equals(other)
}

func (pn PhoneNumber) Equals(other interface{}) bool {
	if o, ok := other.(PhoneNumber); ok {
		return pn.EqualsPhoneNumber(o)
	}
	if o, _ := other.(*PhoneNumber); o != nil {
		return pn.EqualsPhoneNumber(*o)
	}
	return false
}
func (pn PhoneNumber) EqualsPhoneNumber(other PhoneNumber) bool {
	return other.countryCode == pn.countryCode &&
		other.nationalNumber == pn.nationalNumber
}

func (pn PhoneNumber) CountryCode() int32 { return pn.countryCode }
func (pn PhoneNumber) WithCountryCode(countryCode int32) PhoneNumber {
	out := &PhoneNumber{
		countryCode:    countryCode,
		nationalNumber: pn.nationalNumber,
		rawInput:       "", // Clear because it's now not representative
	}
	out.revalidate()
	return *out
}

func (pn PhoneNumber) NationalNumber() int64 { return pn.nationalNumber }
func (pn PhoneNumber) WithNationalNumber(nationalNumber int64) PhoneNumber {
	out := &PhoneNumber{
		countryCode:    pn.countryCode,
		nationalNumber: nationalNumber,
		rawInput:       "", // Clear because it's now not representative
	}
	out.revalidate()
	return *out
}

func (pn PhoneNumber) RawInput() string { return pn.rawInput }

// RawOrFormatted returns a string which prefers raw input with formatted
// string as the default.
func (pn PhoneNumber) RawOrFormatted() string {
	if pn.rawInput != "" {
		return pn.rawInput
	}
	return pn.String()
}

// TODO: consult the standards
func (pn PhoneNumber) String() string {
	nn := uint64(pn.nationalNumber)
	parsed := &phonenumbers.PhoneNumber{
		CountryCode:    &pn.countryCode,
		NationalNumber: &nn,
	}
	return phonenumbers.Format(parsed, phonenumbers.E164)
}

func (pn *PhoneNumber) revalidate() bool {
	if pn.nationalNumber < 0 {
		panic("national number must be positive number")
	}
	nn := uint64(pn.nationalNumber)
	parsed := &phonenumbers.PhoneNumber{
		CountryCode:    &pn.countryCode,
		NationalNumber: &nn,
	}
	pn.isValid = phonenumbers.IsValidNumber(parsed)
	return pn.isValid
}
