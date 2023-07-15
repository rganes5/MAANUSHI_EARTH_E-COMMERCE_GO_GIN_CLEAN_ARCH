package repository

import (
	"context"
	"errors"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	domain "github.com/rganes5/maanushi_earth_e-commerce/pkg/domain"
	"github.com/rganes5/maanushi_earth_e-commerce/pkg/utils"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// The function takes a testing.T parameter t, which is used for test assertions and reporting.
func TestFindByEmail(t *testing.T) {
	// db, mock, err := sqlmock.New()
	//Creates a new SQL mock database connection and assigns it to the db variable. It also creates a mock object (mock) that can be used to set expectations on SQL queries.
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//Initializes a mock GORM database session (gormDB) using the mock database connection (db) created earlier.
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when initializing a mock db session", err)
	}
	AdminRepository := NewAdminRepository(gormDB)
	//Slice of test cases
	tests := []struct {
		name           string
		input          string
		expectedOutput domain.Admin
		buildStub      func(mock sqlmock.Sqlmock)
		//The buildStub field in the test case struct is a function that is responsible for setting up the expected behavior of the mock database connection.
		expectedErr error
	}{
		{
			name:           "non-existing email",
			input:          "nonexisting@gmail.com",
			expectedOutput: domain.Admin{},
			buildStub: func(mock sqlmock.Sqlmock) {
				query := `SELECT * FROM admins WHERE email=$1`
				mock.ExpectQuery(query).WithArgs("nonexisting@gmail.com").WillReturnError(errors.New("invalid Email"))
			},
			expectedErr: errors.New("invalid Email"),
		},
		{
			name:  "valid email",
			input: "ganeshrko007@gmail.com",
			expectedOutput: domain.Admin{
				ID:        1,
				FirstName: "Ganesh",
				LastName:  "R",
				Email:     "ganeshrko007@gmail.com",
				PhoneNum:  "9746226152",
				Password:  "Admin@123",
			},
			buildStub: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "password", "phone_num"}).
					AddRow(1, "Ganesh", "R", "ganeshrko007@gmail.com", "Admin@123", "9746226152")
				query := `SELECT * FROM admins WHERE email=$1`
				mock.ExpectQuery(query).WithArgs("ganeshrko007@gmail.com").WillReturnRows(rows)
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.buildStub(mock)
			actualOutput, actualErr := AdminRepository.FindByEmail(context.TODO(), tt.input)
			//The returned output and error are compared with the expected output and error using assertions from the testify/assert package.
			if tt.expectedErr == nil {
				assert.NoError(t, actualErr)
			} else {
				assert.Equal(t, tt.expectedErr, actualErr)
			}

			assert.Equal(t, tt.expectedOutput, actualOutput)

			err = mock.ExpectationsWereMet()
			if err != nil {
				t.Errorf("unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestSignUpAdmin(t *testing.T) {

	//New() method from sqlmock package create sqlmock database connection and a mock to manage expectations.
	db, mock, err := sqlmock.New()
	//db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	//close the mock db connection after testing.
	defer db.Close()

	//initialize a mock db session
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	//create NewUserRepository mock by passing a pointer to gorm.DB
	adminRepository := NewAdminRepository(gormDB)

	tests := []struct {
		testName       string
		inputField     utils.AdminSignUp
		expectedOutput domain.Admin
		buildStub      func(mock sqlmock.Sqlmock)
		expectedError  error
	}{
		{ // test case for creating a new admin
			testName: "create admin successfull",
			inputField: utils.AdminSignUp{
				FirstName: "Ganesh",
				LastName:  "R",
				Email:     "ganesh@gmail.com",
				PhoneNum:  "9746226152",
				Password:  "Admin@123",
			},
			expectedOutput: domain.Admin{
				//ID:       1,
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
				FirstName: "Ganesh",
				LastName:  "R",
				Email:     "ganesh@gmail.com",
				PhoneNum:  "9746226152",
				Password:  "Admin@123",
			},

			buildStub: func(mock sqlmock.Sqlmock) {
				//simulate the result rows that the query is expected to return.
				rows := sqlmock.NewRows([]string{"first_name", "last_name", "email", "phone_num", "password"}).
					AddRow("Ganesh", "R", "ganesh@gmail.com", "9746226152", "Admin@123")
				//actually above is correct without using quotemeta, regexp.QuoteMeta returns a string that escapes all regular expression metacharacters inside the argument text; the returned string is a regular expression matching the literal text. so https://pkg.go.dev/regexp#QuoteMeta, in the ofcicial documentation we can check
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO admins(first_name,last_name,email,phone_num,password,created_at,updated_at)VALUES($1,$2,$3,$4,$5,NOW(),NOW()) RETURNING *;`)).
					WithArgs("Ganesh", "R", "ganesh@gmail.com", "9746226152", "Admin@123").
					WillReturnRows(rows)
			},

			expectedError: nil,
		},

		{ // test case for creating a new admin with duplicate phone number
			testName: "duplicate phone",
			inputField: utils.AdminSignUp{
				FirstName: "Ganesh",
				LastName:  "R",
				Email:     "ganesh@gmail.com",
				PhoneNum:  "9746226152",
				Password:  "Admin@123",
			},
			expectedOutput: domain.Admin{},

			buildStub: func(mock sqlmock.Sqlmock) {

				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO admins(first_name,last_name,email,phone_num,password,created_at,updated_at)VALUES($1,$2,$3,$4,$5,NOW(),NOW()) RETURNING *;`)).
					WithArgs("Ganesh", "R", "ganesh@gmail.com", "9746226152", "Admin@123").
					WillReturnError(errors.New("phone number already exists- value violates unique constraint 'phone_num'"))

			},

			expectedError: errors.New("phone number already exists- value violates unique constraint 'phone_num'"),
		},

		{ // test case for creating a new admin with duplicate email
			testName: "duplicate email",
			inputField: utils.AdminSignUp{
				FirstName: "Ganesh",
				LastName:  "R",
				Email:     "ganesh@gmail.com",
				PhoneNum:  "9746226152",
				Password:  "Admin@123",
			},
			expectedOutput: domain.Admin{},

			buildStub: func(mock sqlmock.Sqlmock) {

				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO admins(first_name,last_name,email,phone_num,password,created_at,updated_at)VALUES($1,$2,$3,$4,$5,NOW(),NOW()) RETURNING *;`)).
					WithArgs("Ganesh", "R", "ganesh@gmail.com", "9746226152", "Admin@123").
					WillReturnError(errors.New("email already exists- value violates unique constraint 'email'"))

			},

			expectedError: errors.New("email already exists- value violates unique constraint 'email'"),
		},
	}

	for _, tt := range tests {

		t.Run(tt.testName, func(t *testing.T) {

			tt.buildStub(mock)
			actualOutput, actualError := adminRepository.SignUpAdmin(context.Background(), tt.inputField)

			/* This is by using assert from testify package
			if tt.expectedError == nil {
				assert.NoError(t, actualError)
			} else {
				assert.Equal(t, tt.expectedError, actualError)
			}

			if !reflect.DeepEqual(tt.expectedOutput, actualOutput) {
				t.Errorf("got %v, but want %v", actualOutput, tt.expectedOutput)
			}
			*/

			//without using testify assert package, using default testing package
			if tt.expectedError == nil {
				if actualError != nil {
					t.Errorf("expected no error, but got: %v", actualError)
				}
			} else {
				if tt.expectedError.Error() != actualError.Error() {
					t.Errorf("expected error: %v, but got: %v", tt.expectedError, actualError)
				}
			}

			if !reflect.DeepEqual(tt.expectedOutput, actualOutput) {
				t.Errorf("got %+v, but want %+v", actualOutput, tt.expectedOutput)
			}

			// Check that all expectations were met
			err = mock.ExpectationsWereMet()
			if err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}
