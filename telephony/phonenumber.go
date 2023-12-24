package telephony

import (
	"strconv"
	"strings"

	"github.com/alloyzeus/go-azfl/v2/azcore"
	"github.com/nyaruka/phonenumbers"
)

// PhoneNumber represents a phone number as we need.
type PhoneNumber struct {
	countryCode    int32
	nationalNumber int64
	rawInput       string
	isValid        bool
}

var _ azcore.ValueObject = PhoneNumber{}
var _ azcore.ValueObjectAssert[PhoneNumber] = PhoneNumber{}

func NewPhoneNumber(countryCode int32, nationalNumber int64) PhoneNumber {
	return PhoneNumber{countryCode: countryCode, nationalNumber: nationalNumber}
}

func PhoneNumberFromString(input string) (PhoneNumber, error) {
	// Check if the country code is doubled
	if parts := strings.Split(input, "+"); len(parts) == 3 {
		// We assume that the first part was automatically added at the client
		input = "+" + parts[2]
	}

	parsedPhoneNumber, err := phonenumbers.Parse(input, "")
	if err != nil {
		return PhoneNumber{}, err
	}

	pn := PhoneNumber{
		countryCode:    *parsedPhoneNumber.CountryCode,
		nationalNumber: int64(*parsedPhoneNumber.NationalNumber),
		rawInput:       input,
		isValid:        phonenumbers.IsValidNumber(parsedPhoneNumber),
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
		return o.countryCode == pn.countryCode &&
			o.nationalNumber == pn.nationalNumber
	}
	if o, _ := other.(*PhoneNumber); o != nil {
		return o.countryCode == pn.countryCode &&
			o.nationalNumber == pn.nationalNumber
	}
	return false
}

func (pn PhoneNumber) CountryCode() int32    { return pn.countryCode }
func (pn PhoneNumber) NationalNumber() int64 { return pn.nationalNumber }
func (pn PhoneNumber) RawInput() string      { return pn.rawInput }

// RawOrFormatted returns a string which prefers raw input with formatted
// string as the default.
func (pn PhoneNumber) RawOrFormatted() string {
	if pn.rawInput != "" {
		return pn.rawInput
	}
	return pn.String()
}

// TODO: get E.164 string
// TODO: consult the standards
func (pn PhoneNumber) String() string {
	if pn.countryCode == 0 && pn.nationalNumber == 0 {
		return "+"
	}
	return "+" + strconv.FormatInt(int64(pn.countryCode), 10) +
		strconv.FormatInt(pn.nationalNumber, 10)
}
