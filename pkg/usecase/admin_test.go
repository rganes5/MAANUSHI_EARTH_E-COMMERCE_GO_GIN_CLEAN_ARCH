package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rganes5/maanushi_earth_e-commerce/pkg/domain"
	"github.com/rganes5/maanushi_earth_e-commerce/pkg/repository/mockRepo"
	"github.com/stretchr/testify/assert"
)

func TestFindByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	adminRepo := mockRepo.NewMockAdminRepository(ctrl)
	adminUseCase := NewAdminUseCase(adminRepo)
	tests := []struct {
		name           string
		input          string
		expectedOutput domain.Admin
		buildStub      func(adminRepo mockRepo.MockAdminRepository)
		expectedError  error
	}{
		{
			name:  "valid user",
			input: "ganeshrko007@gmail.com",
			expectedOutput: domain.Admin{
				ID:        1,
				FirstName: "Ganesh",
				LastName:  "R",
				Email:     "ganeshrko007@gmail.com",
				PhoneNum:  "9746226152",
				Password:  "Admin@123",
			},
			buildStub: func(adminRepo mockRepo.MockAdminRepository) {
				adminRepo.EXPECT().FindByEmail(gomock.Any(), "ganeshrko007@gmail.com").Times(1).Return(domain.Admin{
					ID:        1,
					FirstName: "Ganesh",
					LastName:  "R",
					Email:     "ganeshrko007@gmail.com",
					PhoneNum:  "9746226152",
					Password:  "Admin@123",
				}, nil)
			},
			expectedError: nil,
		},
		{
			name:           "non-existing user",
			input:          "nonexisting@gmail.com",
			expectedOutput: domain.Admin{},
			buildStub: func(adminRepo mockRepo.MockAdminRepository) {
				adminRepo.EXPECT().FindByEmail(gomock.Any(), "nonexisting@gmail.com").Times(1).Return(domain.Admin{}, errors.New("non-existing user"))
			},
			expectedError: errors.New("non-existing user"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.buildStub(*adminRepo)
			actualOutput, actualError := adminUseCase.FindByEmail(context.TODO(), tt.input)
			assert.Equal(t, tt.expectedOutput, actualOutput)
			assert.Equal(t, tt.expectedError, actualError)
		})
	}
}

// func TestSignUpAdmin(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	adminRepo := mockRepo.NewMockAdminRepository(ctrl)
// 	adminUseCase := NewAdminUseCase(adminRepo)
// 	fmt.Println(adminUseCase)
// 	tests := []struct {
// 		testName       string
// 		inputField     utils.AdminSignUp
// 		expectedOutput domain.Admin
// 		buildStub      func(adminRepo mockRepo.MockAdminRepository)
// 		expectedErr    error
// 	}{}
// }
