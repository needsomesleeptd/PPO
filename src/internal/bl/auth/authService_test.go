package service_test

import (
	service "annotater/internal/bl/auth"
	mock_repository "annotater/internal/mocks/bl/userService/userRepo"
	mock_auth_utils "annotater/internal/mocks/pkg/authUtils"
	"annotater/internal/models"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

const (
	TEST_HASH_KEY       = "test"
	TEST_VALID_LOGIN    = "login"
	TEST_VALID_PASSWORD = "smth"
)

var VALID_USER = models.User{
	Login:    TEST_VALID_LOGIN,
	Password: TEST_VALID_PASSWORD,
}

func TestAuthService_Auth(t *testing.T) {
	type fields struct {
		userRepo       *mock_repository.MockIUserRepository
		passwordHasher *mock_auth_utils.MockIPasswordHasher
		tokenizer      *mock_auth_utils.MockITokenHandler
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
		want    string
	}{
		{
			name: "Valid Created User",
			prepare: func(f *fields) {
				f.userRepo.EXPECT().GetUserByLogin(VALID_USER.Login).Return(nil)
				f.userRepo.EXPECT().CreateUser(&VALID_USER).Return(nil)
			},
			args:    args{VALID_USER},
			wantErr: false,
			errStr:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				userRepo:       mock_repository.NewMockIUserRepository(ctrl),
				passwordHasher: mock_auth_utils.NewMockIPasswordHasher(ctrl),
				tokenizer:      mock_auth_utils.NewMockITokenHandler(ctrl),
				key:            TEST_HASH_KEY,
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			s := service.NewAuthService(f.userRepo, f.passwordHasher, f.tokenizer, f.key)
			err := s.Auth(&tt.args.candidate)
			if tt.wantErr {
				require.Equal(t, tt.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
			}
		})
	}
}
