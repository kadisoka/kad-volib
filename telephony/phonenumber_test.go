package telephony_test

import (
	"testing"

	. "github.com/kadisoka/kad-volib/telephony"
)

func TestConstructor(t *testing.T) {
	testCases := []struct {
		countryCode    int32
		nationalNumber int64
		err            error
	}{
		{0, 0, nil},
		{1, 5552323, nil},
	}

	for idx, tc := range testCases {
		pn, err := PhoneNumberFromCountryCodeAndNationalNumber(tc.countryCode, tc.nationalNumber)
		if err != tc.err {
			t.Errorf("#%d: returned %v, want %v",
				idx+1, err, tc.err)
		}
		if pn.CountryCode() != tc.countryCode {
			t.Errorf("#%d: CountryCode() returned %v, want %v",
				idx+1, pn.CountryCode(), tc.countryCode)
		}
		if pn.NationalNumber() != tc.nationalNumber {
			t.Errorf("#%d: NationalNumber() returned %v, want %v",
				idx+1, pn.NationalNumber(), tc.nationalNumber)
		}
	}
}

func TestEquality(t *testing.T) {
	zero, _ := PhoneNumberFromCountryCodeAndNationalNumber(0, 0)
	one, _ := PhoneNumberFromCountryCodeAndNationalNumber(1, 1)
	zeroOne, _ := PhoneNumberFromCountryCodeAndNationalNumber(0, 1)
	oneZero, _ := PhoneNumberFromCountryCodeAndNationalNumber(1, 0)

	if !zero.Equals(zero) {
		t.Error("!zero.Equals(zero)")
	}
	if !zero.Equals(&zero) {
		t.Error("!zero.Equals(&zero)")
	}
	if zero.Equals(nil) {
		t.Error("zero.Equals(nil)")
	}
	if !zero.Equal(zero) {
		t.Error("zero.Equal(zero)")
	}
	if zero.Equals(one) {
		t.Error("zero.Equals(one)")
	}
	if zero.Equal(one) {
		t.Error("zero.Equal(one)")
	}
	if zero.Equals(zeroOne) {
		t.Error("zero.Equals(zeroOne)")
	}
	if zero.Equals(oneZero) {
		t.Errorf("zero.Equals(oneZero)")
	}
}

func TestChangeCountryCode(t *testing.T) {
	zero, _ := PhoneNumberFromCountryCodeAndNationalNumber(0, 0)

	oneZero := zero.WithCountryCode(1)
	if zero.CountryCode() != 0 || zero.NationalNumber() != 0 {
		t.Error("zero should not be affected")
	}
	if oneZero.CountryCode() != 1 {
		t.Errorf("oneZero.CountryCode() returned %v, want 1",
			oneZero.CountryCode())
	}
	if oneZero.NationalNumber() != 0 {
		t.Errorf("oneZero.NationalNumber() returned %v, want 0",
			oneZero.NationalNumber())
	}
	if oneZero.EqualsPhoneNumber(zero) {
		t.Errorf("oneZero.EqualsPhoneNumber(zero) returned true, want false")
	}

	zeroZero := oneZero.WithCountryCode(0)
	if zeroZero.Equals(oneZero) {
		t.Error("zeroZero.Equals(oneZero)")
	}
	if !zeroZero.Equals(zero) {
		t.Error("zeroZero.Equal(zero)")
	}
}

func TestChangeNationalNumber(t *testing.T) {
	zero, _ := PhoneNumberFromCountryCodeAndNationalNumber(0, 0)

	zeroOne := zero.WithNationalNumber(1)
	if zero.CountryCode() != 0 || zero.NationalNumber() != 0 {
		t.Error("zero should not be affected")
	}
	if zeroOne.CountryCode() != 0 {
		t.Errorf("oneZero.CountryCode() returned %v, want 0",
			zeroOne.CountryCode())
	}
	if zeroOne.NationalNumber() != 1 {
		t.Errorf("oneZero.NationalNumber() returned %v, want 1",
			zeroOne.NationalNumber())
	}
	if zeroOne.EqualsPhoneNumber(zero) {
		t.Errorf("oneZero.EqualsPhoneNumber(zero) returned true, want false")
	}

	zeroZero := zeroOne.WithNationalNumber(0)
	if zeroZero.Equals(zeroOne) {
		t.Error("zeroZero.Equals(oneZero)")
	}
	if !zeroZero.Equals(zero) {
		t.Error("zeroZero.Equal(zero)")
	}
}
