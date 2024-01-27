package email

import (
	"regexp"
	"strings"

	"github.com/alloyzeus/go-azfl/v2/azcore"
	"github.com/alloyzeus/go-azfl/v2/errors"
)

// TODO:
// - This package should provide a set of Normalization and Validation
//   profiles
//   - e.g., some some use cases might require high interoperability, e.g.
//     they accept @ or " character in local-part. While some others are
//     based on common practices.
// - Users of this package could also define their own Normalization and
//   Validation profiles

type Address struct {
	localPart  string
	domainPart string
	rawInput   string
}

var _ azcore.ValueObject = Address{}
var _ azcore.ValueObjectAssert[Address] = Address{}

func AddressFromLocalAndDomainParts(localPart, domainPart string) (Address, error) {
	return Address{
		localPart:  localPart,
		domainPart: domainPart,
	}, nil
}

func AddressFromString(input string) (Address, error) {
	localPart, domainPart, _, err := normalizeAddressString(input)
	if err != nil {
		return Address{}, err
	}

	return Address{
		localPart:  localPart,
		domainPart: domainPart,
		rawInput:   input,
	}, nil
}

// TODO: at least we check for common address convention. Also, by profile.
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

func normalizeAddressString(input string) (localPart, domainPart, normalized string, err error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", "", "", errors.ArgMsg("input", "empty")
	}

	lastAt := strings.LastIndex(input, "@")
	if lastAt == -1 {
		return "", "", "", errors.ArgMsg("input", "malformed")
	}

	localPart = input[:lastAt]
	domainPart = input[lastAt+1:]

	if localPart == "" {
		return "", "", "", errors.ArgMsg("input", "local part empty")
	}
	if domainPart == "" || !addressDomainRE.MatchString(domainPart) {
		return "", "", "", errors.ArgMsg("input", "domain part malformed")
	}

	// TODO(exa): domain normalization
	// - expect punycode (use IDNA)
	// - normalize, e.g., lowercase ASCII characters
	domainPart = strings.ToLower(domainPart)

	// TODO(exa): local part normalization
	// - domain-specific?
	// - for now, let's lowercase everything
	localPart = strings.ToLower(localPart)

	return localPart, domainPart, localPart + "@" + domainPart, nil
}
