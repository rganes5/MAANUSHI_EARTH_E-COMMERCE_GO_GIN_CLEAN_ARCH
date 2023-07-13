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

func TestFindByEmail(t *testing.T) {
	// db, mock, err := sqlmock.New()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when initializing a mock db session", err)
	}
	AdminRepository := NewAdminRepository(gormDB)

	tests := []struct {
		name           string
		input          string
		expectedOutput domain.Admin
		buildStub      func(mock sqlmock.Sqlmock)
		expectedErr    error
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
