package util

import "github.com/nyaruka/phonenumbers"

func PhoneNumberValidator(phone string) bool {
	//validating to US phone numbers
	// might later use varingin country code
	_, err := phonenumbers.Parse(phone, "US")
	if err != nil {
		return false
	}

	return true
}
