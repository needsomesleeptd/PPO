package service_test

import (
	mock_repository "annotater/internal/mocks/bl/userService/userRepo"
	mock_auth_utils "annotater/internal/mocks/pkg/authUtils"
	"annotater/internal/models"
	"testing"
)

func TestAuthService_Auth(t *testing.T) {
	type fields struct {
		userRepo       mock_repository.MockIUserRepository
		passwordHasher mock_auth_utils.MockIPasswordHasher
		tokenizer      mock_auth_utils.MockITokenHandler
		key            string
	}
	type args struct {
		candidate models.User
	}
	tests := []struct {
		name    string
		fields  fields
		prepare func(f *fields)
		args    args
		wantErr bool
		errStr  error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serv := &AuthService{
				userRepo:       tt.fields.userRepo,
				passwordHasher: tt.fields.passwordHasher,
				tokenizer:      tt.fields.tokenizer,
				key:            tt.fields.key,
			}
			if err := serv.Auth(tt.args.candidate); (err != nil) != tt.wantErr {
				t.Errorf("AuthService.Auth() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
