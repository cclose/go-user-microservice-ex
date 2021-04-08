## Prerequisites

You must have docker installed, compatible with docker compose v3.7

## Set up

Create a `.env` file next to `docker-compose.yml`. Replace `<username>` and `<password>` with your desired values:
```
PG_USER_USER=<username>
PG_USER_PASSWORD=<password>
```

This file is read by the Dockerfile to distribute credentials into the containers safely

## Build

Build the App stack with:

    docker-composer up --build

## Tests 

### Unit Test
Unit tests for the User model's validation functions can be found in the `models` folder. With go installed run:

    cd user-service/src/models; go test

If attempting to run the unit test within the docker container, you must first install gcc with:

    apk add --no-cache gcc musl-dev

gcc is required for `go test` to run, but is left out of the base image to keep it lightweight. It is best to run `go test` from your development workspace

### Integration Test

Visit http://localhost:8080/test to see a Javascript Integration Test.

## API Specification

### JSON Schema

The User model JSON format is as such:

```json
{
  "id": "",
  "username":"",
  "firstname":"",
  "middlename":"",
  "lastname":"",
  "email":"",
  "telephone":"",
  "password":""
}
```

#### Password Field
The `password` field is never returned by any of the routes.
Create (POST) and update (PUT) routes accept the password in plaintext, but all passwords are hashed and encoded before being stored.

Note: Ideally the server should ONLY accept plaintext passwords when running in HTTPS, but I did not have SSL certs to use at this time.

#### Id Field
The `id` field is not accepted as part of the CREATE route.


#### Validation

The User model provides validation on the following fields. If validation fails, code 400 is returned and json explaining the failure:

Field | Validation
----- | ----------
username | must be unique, be between 5 and 25 characters, and only contain alphanumerics
email | must be unique and vaguely look like an email address (contain 1 and only 1 @ )
firstname | is required
lastname | is required
telephone | is required and must be of the form (###) ###-####[ x#####]. Extension is optional, max length of 5, with an optional space before the x
password | password must be between 8 and 25 characters, contain at least 1 of: lower case, upper case, number, and special character

### Routes

#### Test App

Route: `/test` Method: `GET`

Serves an HTML and Vanilla Javascript Integration test to demonstrate the functionality of the service easily

#### Get All Users
Route: `/api/v1/user` Method: `GET` Returns: `json`

Fetches all users in the database, adjust by limit and offset

Query Parameters:

Key | Type | Description
--- | ---- | ---------
limit | integer | Limits the number of results. Default 100
offset | integer | sets an offset for when to begin serving results. Only works with limit. Default 0

Response Codes:

Code | Reason
---- | ------
200  | Success. If empty, returns empty array
400  | Query parameters invalid
500  | an error occurred with the service


#### Get User By ID
Route: `/api/v1/user/{id}` Method: `GET` Returns: `json`

Fetches a single user by the id in the route

Route Parameters:

Key | Type | Description
--- | ---- | ---------
id | integer | The id of the user to get

Response Codes: 

Code | Reason
---- | ------
200  | Success
404  | user with that id not found
500  | an error occurred with the service


#### Create User
Route: `/api/v1/user` Method: `POST` Accepts: `json` Returns `json`

Fetches all users in the database, adjust by limit and offset

Response Codes:

Code | Reason
---- | ------
201  | Success
400  | The results is malformed or fails validation
409  | A uniqueness constraint was violated (username, email)
415  | Wrong content-type (Json only)
500  | an error occurred with the service


#### Update User
Route: `/api/v1/user/{id}` Method: `PUT` Accepts: `json` Returns `json`

Updates the specified user with the PUT Json.
If password is omitted, it is left unchanged. If included, it will be hashed before storage

The submitted id must match the id in the url. Both are required.

Route Parameters:

Key | Type | Description
--- | ---- | ---------
id | integer | The id of the user to update

Response Codes:

Code | Reason
---- | ------
200  | Success
400  | The results is malformed or fails validation
404  | No user with that id exists
409  | A uniqueness constraint was violated (username, email)
415  | Wrong content-type (Json only)
500  | an error occurred with the service


#### Delete User
Route: `/api/v1/user/{id}` Method: `DELETE` Returns `json`

Deletes the user with the specified ID

Route Parameters:

Key | Type | Description
--- | ---- | ---------
id | integer | The id of the user to update

Response Codes:

Code | Reason
---- | ------
200  | Success. User with that id deleted
404  | No user with that id exists
500  | an error occurred with the service


#### Authenticate User
Route: `/api/v1/user/auth` Method: `POST` Returns `json`

Using Basic HTTP Authentication Header, tests if the provided username and password match

Response Codes:

Code | Reason
---- | ------
200  | Success. Credentials Valid
401  | Failed: credentials invalid or unparsable
500  | an error occurred with the service