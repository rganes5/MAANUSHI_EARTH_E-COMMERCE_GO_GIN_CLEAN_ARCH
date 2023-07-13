package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	domain "github.com/rganes5/maanushi_earth_e-commerce/pkg/domain"
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
				Model:     gorm.Model{ID: 1},
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
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when initializing a mock db session", err)
	}

	AdminRepository := NewAdminRepository(gormDB)

	tests := []struct {
		name        string
		input       domain.Admin
		buildStub   func(mock sqlmock.Sqlmock, admin domain.Admin)
		expectedErr error
	}{
		{
			name: "SignUp admin success",
			input: domain.Admin{
				FirstName: "Ganesh",
				LastName:  "R",
				Email:     "ganeshrko007@gmail.com",
				PhoneNum:  "9746226152",
				Password:  "Admin@123",
			},
			buildStub: func(mock sqlmock.Sqlmock, admin domain.Admin) {
				mock.ExpectExec(`INSERT INTO admins`).
					WithArgs(admin.FirstName, admin.LastName, admin.Email, admin.PhoneNum, admin.Password).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedErr: nil,
		},
		{
			name: "Failed to signup admin due to constraint violation",
			input: domain.Admin{
				FirstName: "Ganesh",
				LastName:  "R",
				Email:     "ganeshrko007@gmail.com",
				PhoneNum:  "9746226152",
				Password:  "Admin@123",
			},
			buildStub: func(mock sqlmock.Sqlmock, admin domain.Admin) {
				mock.ExpectExec(`INSERT INTO admins`).
					WithArgs(admin.FirstName, admin.LastName, admin.Email, admin.PhoneNum, admin.Password).
					WillReturnError(errors.New("unique constraint violation"))
			},
			expectedErr: errors.New("unique constraint violation"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.buildStub(mock, tt.input)

			actualErr := AdminRepository.SignUpAdmin(context.TODO(), tt.input)

			if tt.expectedErr == nil {
				assert.NoError(t, actualErr)
			} else {
				assert.Equal(t, tt.expectedErr, actualErr)
			}

			err = mock.ExpectationsWereMet()
			if err != nil {
				t.Errorf("unfulfilled expectations: %s", err)
			}
		})
	}
}
