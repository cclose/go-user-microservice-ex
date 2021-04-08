package models

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/cclose/go-user-microservice-ex/user-service/src/database"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"strings"
)

type UserModel struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	Password   string `json:"password,omitempty"` //omit on empty so we avoid sending this field to the client
	FirstName  string `json:"firstname"`
	MiddleName string `json:"middlename"`
	LastName   string `json:"lastname"`
	Email      string `json:"email"`
	Telephone  string `json:"telephone"`
}

const UserSchema string = `
CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	username TEXT NOT NULL UNIQUE,
	password_hash TEXT NOT NULL,
	firstname TEXT NOT NULL,
	middlename TEXT,
	lastname TEXT NOT NULL,
	email TEXT NOT NULL UNIQUE,
	telephone TEXT NOT NULL
);
`

// hashPassword safely converts a plaintext password into a salted, hashed, base64 value
func hashPassword(password string) (string, error) {
	// Hash the password with bcrypt Note: bcrypt autosalts!
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// Base64 encode the hash so we can transport it more easily
	encodedPasswordHash := base64.URLEncoding.EncodeToString(passwordHash)

	return encodedPasswordHash, nil
}

// CheckPassword compares the specified plaintext password with the stared hash to see if they match
func CheckPassword(hashedPassword, password string) bool {
	//We store our passwords as base64, so decode this
	decodedHash, _ := base64.URLEncoding.DecodeString(hashedPassword)
	err := bcrypt.CompareHashAndPassword(
		decodedHash, []byte(password))
	return err == nil
}

// handlePassword encapsulates all the logic to validate and hash the password
func (user *UserModel) handlePassword() error {
	if !ValidatePassword(user.Password) {
		return errors.New("UserModel failed validation:\n\t- Password must be between 8 and 25 characters, " +
			"contain upper and lower case, at least one number, and at least one symbol from " + PASS_SPECIAL_CHARS)
	}

	var err error
	user.Password, err = hashPassword(user.Password)
	if err != nil {
		return errors.New("failed to Encode Password")
	}

	return nil
}

// validateEmail verifies that the email address is a valid email
func validateEmail(emailAddress string) bool {
	// TODO: Improve validation. Currently only tests if it has an @ sign. Email Validation is tricky. Better get a lib
	emailParts := strings.Split(emailAddress, "@")
	return len(emailParts) == 2
}

var US_PHONE_REGEX = regexp.MustCompile("^\\(\\d{3}\\) \\d{3}-\\d{4}(\\s?x\\d{1,5})?$")

//var INT_PHONE_REGEX = regexp.MustCompile("")
// validateTelephone verifies that the Telephone number is a valid phone number
func validateTelephone(phone string) bool {
	// Test US Phone numbers of the form (###) ###-####[x#####]
	if US_PHONE_REGEX.MatchString(phone) {
		return true
	}

	//TODO add support for international numbers
	//// Test International Phone numbers
	//if INT_PHONE_REGEX.MatchString(phone) {
	//	return true
	//}

	// we didn't match any of our regexs, so this is invalid
	return false
}

var USERNAME_REGEX = regexp.MustCompile("[^a-zA-Z0-9]")

// validateUsername verifies that the Username is valid per the standards
func validateUsername(username string) bool {
	if len(username) >= 5 && len(username) <= 25 {
		// This regex tests for invalid characters (non alphanumeric)
		if !USERNAME_REGEX.MatchString(username) {
			return true
		}
	}

	return false
}

var PASS_HAS_UPPERCASE = regexp.MustCompile("[A-Z]")
var PASS_HAS_LOWERCASE = regexp.MustCompile("[a-z]")
var PASS_HAS_NUMBER = regexp.MustCompile("[0-9]")

const PASS_SPECIAL_CHARS string = "!@#$%&?+.,\\*\\^-_=<>\\[\\]\\(\\)\\{\\}" //@#$%^&*<>?-_=+~"
var PASS_HAS_SPECIAL = regexp.MustCompile(fmt.Sprintf("[%s]", PASS_SPECIAL_CHARS))
var PASS_HAS_ILLEGAL = regexp.MustCompile(fmt.Sprintf("[^a-zA-Z0-9%s]", PASS_SPECIAL_CHARS))

// validatePassword verifies that the plain text password meets requirements
func ValidatePassword(password string) bool {
	// passwords must be between 8 and 25 characters
	if len(password) < 8 || len(password) > 25 {
		return false
	}
	// must have a capitol letter
	if !PASS_HAS_UPPERCASE.MatchString(password) {
		return false
	}
	// must have a lower case letter
	if !PASS_HAS_LOWERCASE.MatchString(password) {
		return false
	}
	// must have a number
	if !PASS_HAS_NUMBER.MatchString(password) {
		return false
	}
	// must possess a special character
	if !PASS_HAS_SPECIAL.MatchString(password) {
		return false
	}
	// must not possess any other characters
	if PASS_HAS_ILLEGAL.MatchString(password) {
		return false
	}

	return true
}

// validate Validates that all fields are included and container proper values
// returns the list of validation errors
// Does not validate the password as this is handled independently
func (user UserModel) Validate() (errs []string) {
	if user.FirstName == "" {
		errs = append(errs, "FirstName is not specified!")
	}

	if user.LastName == "" {
		errs = append(errs, "LastName is not specified!")
	}

	if user.Email == "" {
		errs = append(errs, "Email is not specified!")
	} else if !validateEmail(user.Email) {
		errs = append(errs, "Invalid Email specified!")
	}

	if user.Telephone == "" {
		errs = append(errs, "Telephone is not specified!")
	} else if !validateTelephone(user.Telephone) {
		errs = append(errs, "Invalid Telephone specified! Accepts (###) ###-####[[ ]x#####]")
	}

	if user.Username == "" {
		errs = append(errs, "Username is not specified!")
	} else if !validateUsername(user.Username) {
		errs = append(errs, "Invalid Username: must be between 5 and 25 characters and be only alphanumeric")
	}

	return
}

func (user *UserModel) Create(db *database.PostGresDB) error {
	if user.ID != 0 {
		return errors.New("ID must be null when creating a User")
	}

	valErrors := user.Validate()
	if len(valErrors) > 0 {
		return errors.New(
			fmt.Sprintf("UserModel failed validation:\n\t- %s", strings.Join(valErrors, "\n\t- ")))
	}

	err := user.handlePassword()
	if err != nil {
		return err
	}

	insertStmt := `INSERT INTO users (username, password_hash, firstname, middlename, lastname, email, telephone)
		VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	err = db.PgDbSession.QueryRow(insertStmt, user.Username, user.Password, user.FirstName, user.MiddleName,
		user.LastName, user.Email, user.Telephone).Scan(&user.ID)
	if err != nil {
		if database.DuplicateKeyError(err) {
			return database.ErrDuplicateKey
		}
		return err
	}

	return nil
}

func (user *UserModel) Delete(db *database.PostGresDB) error {
	if user.ID == 0 {
		return sql.ErrNoRows
	}

	deleteStmt := `DELETE FROM users WHERE id = $1`
	res, err := db.PgDbSession.Exec(deleteStmt, user.ID)
	if err != nil {
		return err
	}

	var rows int64
	rows, _ = res.RowsAffected()
	// If we didnt delete anything, return the no rows error to tell the controller to 404
	if int(rows) == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (user *UserModel) Update(db *database.PostGresDB) error {
	if user.ID == 0 {
		return sql.ErrNoRows
	}

	valErrors := user.Validate()
	if len(valErrors) > 0 {
		return errors.New(
			fmt.Sprintf("UserModel failed validation:\n\t- %s", strings.Join(valErrors, "\n\t- ")))
	}

	updateStmt := `UPDATE users SET username = $2, firstname = $3, middlename = $4, lastname = $5, email = $6, telephone = $7`

	// If the password is blank, the user isn't updating it at this time
	if user.Password != "" {
		err := user.handlePassword()
		if err != nil {
			return errors.New("Bad Password: " + user.Password + ":: " + err.Error())
		}
		updateStmt += `, password_hash = $8`
	}

	updateStmt += ` WHERE id = $1`
	var res sql.Result
	var err error
	if user.Password != "" {
		res, err = db.PgDbSession.Exec(updateStmt, user.ID, user.Username, user.FirstName, user.MiddleName,
			user.LastName, user.Email, user.Telephone, user.Password)
	} else {
		res, err = db.PgDbSession.Exec(updateStmt, user.ID, user.Username, user.FirstName, user.MiddleName,
			user.LastName, user.Email, user.Telephone)
	}
	if err != nil {
		return err
	}

	var rows int64
	rows, _ = res.RowsAffected()
	// If we didnt update anything, return the no rows error to tell the controller to 404
	if int(rows) == 0 {
		return sql.ErrNoRows
	}

	return nil
}

const USER_GET_FIELDLIST string = "id, username, firstname, middlename, lastname, email, telephone"

func GetUsers(db *database.PostGresDB, field, value string, limit, offset int) (users []UserModel, err error) {
	var params []interface{}

	selectStmt := `SELECT ` + USER_GET_FIELDLIST + ` FROM users`
	if field != "id" && field != "username" && field != "firstname" && field != "middlename" && field != "lastname" &&
		field != "email" && field != "telephone" && field != "all" {
		err = errors.New(fmt.Sprintf("Unsupported search field |%s|", field))
		return
	}
	if field != "all" {
		params = append(params, value)
		selectStmt += " WHERE " + field + " = $1 "
		if offset > 0 {
			selectStmt += " LIMIT $2"
			params = append(params, limit)
		} else if limit > 0 {
			selectStmt += " LIMIT $2 OFFSET $3"
			params = append(params, limit)
			params = append(params, offset)
		}
	} else {
		if offset > 0 {
			selectStmt += " LIMIT $1"
			params = append(params, limit)
		} else if limit > 0 {
			selectStmt += " LIMIT $1 OFFSET $2"
			params = append(params, limit)
			params = append(params, offset)
		}
	}

	var rows *sql.Rows
	rows, err = db.PgDbSession.Query(selectStmt, params...)
	if err != nil {
		return
	}

	for rows.Next() {
		var user UserModel
		err = rows.Scan(&user.ID, &user.Username, &user.FirstName, &user.MiddleName, &user.LastName, &user.Email,
			&user.Telephone)
		if err != nil {
			return
		}
		users = append(users, user)
	}

	return
}

// GetPassword fetchs the password hash for the user from the Database
func GetUserCredentials(db *database.PostGresDB, username string) (hashedPassword string, err error) {
	selectStmt := `SELECT password_hash FROM users WHERE username = $1`
	row := db.PgDbSession.QueryRow(selectStmt, username)
	if row == nil {
		err = sql.ErrNoRows
		return
	}

	err = row.Scan(&hashedPassword)

	return
}
