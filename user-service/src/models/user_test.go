package models

import "testing"

// TestUserPassword Checks both the hashPassword function and the CheckPassword function.
func TestUserPassword(t *testing.T) {
	password := "SuperSecure123@#Pass"
	hash, err := hashPassword(password)
	if err != nil {
		t.Errorf("Caught error while hasing password |%s|: %s", password, err)
	}

	if !CheckPassword(hash, password) {
		t.Errorf("Password |%s| did not hash and then check successfully", password)
	}

	if CheckPassword(hash, "Imnotetheoriginsla") {
		t.Errorf("Incorrect password validated against the hash!!")
	}
}

func TestValidateEmail(t *testing.T) {
	//Test a good email
	goodEmail := "bob@bob.win"
	if !validateEmail(goodEmail) {
		t.Errorf("Email %s expected to pass, failed validation", goodEmail)
	}

	//Test an email w/o an @
	noDomain := "nodomain.fail"
	if validateEmail(noDomain) {
		t.Errorf("Email %s expected to fail, passed", noDomain)
	}

	//Test an email with too many @'s
	tooManyAts := "too@many@atsigns.fail"
	if validateEmail(tooManyAts) {
		t.Errorf("Email %s expected to fail, passed", tooManyAts)
	}
}

func TestValidateTelephone(t *testing.T) {
	//Test a standard US number
	basicUsNum := "(555) 555-5555"
	if !validateTelephone(basicUsNum) {
		t.Errorf("Telephone %s expected to pass, failed", basicUsNum)
	}

	//Test a standard US number with extension
	basicUsNumEx1 := "(555) 555-5555x123"
	if !validateTelephone(basicUsNumEx1) {
		t.Errorf("Telephone %s expected to pass, failed", basicUsNumEx1)
	}

	//Test a standard US number with extension with a space
	basicUsNumEx2 := "(555) 555-5555 x123"
	if !validateTelephone(basicUsNumEx2) {
		t.Errorf("Telephone %s expected to pass, failed", basicUsNumEx2)
	}

	//Test a standard US number with extension that's too long
	extTooLong := "(555) 555-5555 x123456"
	if validateTelephone(extTooLong) {
		t.Errorf("Telephone %s expected to fail, passed", extTooLong)
	}

	//Test a standard US number where the area code isn't in ( )
	noParens := "555 555-5555 x123456"
	if validateTelephone(noParens) {
		t.Errorf("Telephone %s expected to fail, passed", noParens)
	}

	//Test a bunch of numbers
	numberPile := "5555555555x123456"
	if validateTelephone(numberPile) {
		t.Errorf("Telephone %s expected to fail, passed", numberPile)
	}

	//Test random garbage
	gibberish := "123lkjdf9@#$#$%klj"
	if validateTelephone(gibberish) {
		t.Errorf("Telephone %s expected to fail, passed", gibberish)
	}
}

func TestValidateUsername(t *testing.T) {
	//Test a good username
	goodName := "johnny005"
	if !validateUsername(goodName) {
		t.Errorf("Username %s expected to pass, failed", goodName)
	}

	//Test a short username
	shortName := "jon"
	if validateUsername(shortName) {
		t.Errorf("Username %s expected to fail, passed", shortName)
	}

	//Test a long username
	longName := "jonnnnnnnnnnnnnnnnnnnnnnny"
	if validateUsername(longName) {
		t.Errorf("Username %s expected to fail, passed", longName)
	}

	//Test a name with non alphanumberics
	badName := "jo!@#$%^&*()<>"
	if validateUsername(badName) {
		t.Errorf("Username %s expected to fail, passed", badName)
	}
}

func TestValidatePassword(t *testing.T) {
	//Test a good Password
	goodPass := "goodPass034!!"
	if !ValidatePassword(goodPass) {
		t.Errorf("Password %s expected to pass, failed", goodPass)
	}

	//Test a blank password
	blankPass := ""
	if ValidatePassword(blankPass) {
		t.Errorf("Password %s expected to fail, passed", blankPass)
	}

	//Test a short password
	shortPass := "sH0rt!"
	if ValidatePassword(shortPass) {
		t.Errorf("Password %s expected to fail, passed", shortPass)
	}

	//Test a long password
	longPass := "longPassHasMoreThan25Characters!!"
	if ValidatePassword(longPass) {
		t.Errorf("Password %s expected to fail, passed", longPass)
	}

	//Test no Uppercase
	noUpperPass := "badpassword123!!"
	if ValidatePassword(noUpperPass) {
		t.Errorf("Password %s expected to fail, passed", noUpperPass)
	}

	//Test no lowercase
	noLowerPass := "BADPASSWORD123!!"
	if ValidatePassword(noLowerPass) {
		t.Errorf("Password %s expected to fail, passed", noLowerPass)
	}

	//Test no numbers
	noNumPass := "badPASSword!!"
	if ValidatePassword(noNumPass) {
		t.Errorf("Password %s expected to fail, passed", noNumPass)
	}

	//Test no symbols
	noSymPass := "badPASSword123"
	if ValidatePassword(noSymPass) {
		t.Errorf("Password %s expected to fail, passed", noSymPass)
	}

	//Test illegals
	illegalPass := "badPASSword123!!;;"
	if ValidatePassword(illegalPass) {
		t.Errorf("Password %s expected to fail, passed", illegalPass)
	}
}
