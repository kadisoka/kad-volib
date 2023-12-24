package email

import (
	"regexp"
	"strings"

	"github.com/alloyzeus/go-azfl/v2/azcore"
	"github.com/alloyzeus/go-azfl/v2/errors"
)

type Address struct {
	localPart  string
	domainPart string
	rawInput   string
}

var _ azcore.ValueObject = Address{}
var _ azcore.ValueObjectAssert[Address] = Address{}

func AddressFromString(input string) (Address, error) {
	parts := strings.SplitN(input, "@", 2) //TODO: from the back
	if len(parts) < 2 {
		return Address{}, errors.ArgMsg("input", "malformed")
	}
	//TODO(exa): normalize localPart and domainPart
	if parts[0] == "" {
		return Address{}, errors.ArgMsg("input", "local part empty")
	}
	if parts[1] == "" || !addressDomainRE.MatchString(parts[1]) {
		return Address{}, errors.ArgMsg("input", "domain part malformed")
	}
	//TODO(exa): perform more extensive checking

	return Address{
		localPart:  parts[0],
		domainPart: strings.ToLower(parts[1]),
		rawInput:   input,
	}, nil
}

// TODO: at least we check for common address convention
func (addr Address) IsStaticallyValid() bool {
	return addr.localPart != "" && addr.domainPart != ""
}

func (addr Address) Clone() Address { return addr }

func (addr Address) Equal(other interface{}) bool {
	return addr.Equals(other)
}

func (addr Address) Equals(other interface{}) bool {
	//TODO: compare the normalized representations
	if o, ok := other.(Address); ok {
		return strings.EqualFold(o.domainPart, addr.domainPart) &&
			o.localPart == addr.localPart
	}
	if o, _ := other.(*Address); o != nil {
		return strings.EqualFold(o.domainPart, addr.domainPart) &&
			o.localPart == addr.localPart
	}
	return false
}

func (addr Address) String() string {
	return addr.localPart + "@" + addr.domainPart
}

func (addr Address) LocalPart() string {
	return addr.localPart
}

func (addr Address) DomainPart() string {
	return addr.domainPart
}

func (addr Address) RawInput() string {
	return addr.rawInput
}

// RawOrFormatted returns a string which prefers raw input with formatted
// string as the default.
func (addr Address) RawOrFormatted() string {
	if addr.rawInput != "" {
		return addr.rawInput
	}
	return addr.String()
}

// NOTE: actually, it's not recommended to use regex to
// identify if a string is an email address:
// https://www.regular-expressions.info/email.html
var addressRE = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
var addressDomainRE = regexp.MustCompile("^[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func IsValidAddress(str string) bool {
	return addressRE.MatchString(str)
}
