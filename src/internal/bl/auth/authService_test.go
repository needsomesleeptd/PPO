package service_test

import (
	service "annotater/internal/bl/auth"
	mock_repository "annotater/internal/mocks/bl/userService/userRepo"
	mock_auth_utils "annotater/internal/mocks/pkg/authUtils"
	"annotater/internal/models"
	"testing"

	"errors"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

const (
	TEST_HASH_KEY       = "test"
	TEST_VALID_LOGIN    = "login"
	TEST_VALID_PASSWORD = "passed"
	TEST_HASH_PASSWD    = "hashed_passwd"
	TEST_VALID_TOKEN    = "token"
)

var VALID_USER = models.User{
	Login:    TEST_VALID_LOGIN,
	Password: TEST_VALID_PASSWORD,
}

var VALID_USER_IN_DB = models.User{
	Login:    TEST_VALID_LOGIN,
	Password: TEST_HASH_PASSWD,
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
		name      string
		fields    fields
		prepare   func(f *fields)
		changeArg func(a *args)
		args      args
		wantErr   bool
		errStr    error
		want      models.User
	}{
		{
			name: "Valid Created User",
			prepare: func(f *fields) {
				f.userRepo.EXPECT().CreateUser(&VALID_USER_IN_DB).Return(nil)
				f.passwordHasher.EXPECT().GenerateHash(VALID_USER.Password).Return(TEST_HASH_PASSWD, nil)

			},
			args:    args{VALID_USER},
			want:    VALID_USER_IN_DB,
			wantErr: false,
			errStr:  nil,
		},
		{
			name:    "No login",
			args:    args{models.User{Password: TEST_VALID_PASSWORD}},
			want:    models.User{},
			wantErr: true,
			errStr:  service.ErrNoLogin,
		},
		{
			name:    "No Passwd",
			args:    args{models.User{Login: TEST_VALID_LOGIN}},
			want:    models.User{},
			wantErr: true,
			errStr:  service.ErrNoPasswd,
		},
		{
			name: "Hash error",
			prepare: func(f *fields) {
				f.passwordHasher.EXPECT().GenerateHash(VALID_USER.Password).Return("", errors.New(""))

			},
			args:    args{VALID_USER},
			want:    models.User{},
			wantErr: true,
			errStr:  errors.Join(service.ErrGeneratingHash, errors.New("")),
		},
		{
			name: "CreateUser error",
			prepare: func(f *fields) {
				f.passwordHasher.EXPECT().GenerateHash(VALID_USER.Password).Return(TEST_HASH_PASSWD, nil)
				f.userRepo.EXPECT().CreateUser(&VALID_USER_IN_DB).Return(errors.New(""))
			},
			args:    args{VALID_USER},
			want:    models.User{},
			wantErr: true,
			errStr:  errors.Join(service.ErrCreatingUser, errors.New("")),
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
			err := s.SignUp(&tt.args.candidate)
			if tt.wantErr {
				require.Equal(t, tt.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestAuthService_SignIn(t *testing.T) {
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
		name      string
		fields    fields
		prepare   func(f *fields)
		changeArg func(a *args)
		args      args
		wantErr   bool
		errStr    error
		want      string
	}{
		{
			name: "Valid SignIn User",
			prepare: func(f *fields) {
				f.userRepo.EXPECT().GetUserByLogin(VALID_USER.Login).Return(&VALID_USER_IN_DB, nil)
				f.passwordHasher.EXPECT().ComparePasswordhash(VALID_USER.Password, VALID_USER_IN_DB.Password).Return(nil)
				f.tokenizer.EXPECT().GenerateToken(VALID_USER, TEST_HASH_KEY).Return(TEST_VALID_TOKEN, nil)
			},
			args:    args{VALID_USER},
			want:    TEST_VALID_TOKEN,
			wantErr: false,
			errStr:  nil,
		},
		{
			name:    "No login",
			args:    args{models.User{Password: TEST_VALID_PASSWORD}},
			want:    "",
			wantErr: true,
			errStr:  service.ErrNoLogin,
		},
		{
			name:    "No Passwd",
			args:    args{models.User{Login: TEST_VALID_LOGIN}},
			want:    "",
			wantErr: true,
			errStr:  service.ErrNoPasswd,
		},
		{
			name: "GetUser Error",
			args: args{VALID_USER},
			prepare: func(f *fields) {
				f.userRepo.EXPECT().GetUserByLogin(VALID_USER.Login).Return(nil, errors.New(""))

			},
			want:    "",
			wantErr: true,
			errStr:  errors.Join(service.ErrWrongLogin, errors.New("")),
		},
		{
			name: "CmpHash Error",
			args: args{VALID_USER},
			prepare: func(f *fields) {
				f.userRepo.EXPECT().GetUserByLogin(VALID_USER.Login).Return(&VALID_USER_IN_DB, nil)
				f.passwordHasher.EXPECT().ComparePasswordhash(VALID_USER.Password, VALID_USER_IN_DB.Password).Return(errors.New(""))
			},
			want:    "",
			wantErr: true,
			errStr:  errors.Join(service.ErrHashPasswdMatch, errors.New("")),
		},
		{
			name: "Token generation Error",
			args: args{VALID_USER},
			prepare: func(f *fields) {
				f.userRepo.EXPECT().GetUserByLogin(VALID_USER.Login).Return(&VALID_USER_IN_DB, nil)
				f.passwordHasher.EXPECT().ComparePasswordhash(VALID_USER.Password, VALID_USER_IN_DB.Password).Return(nil)
				f.tokenizer.EXPECT().GenerateToken(VALID_USER, TEST_HASH_KEY).Return("", errors.New(""))
			},
			want:    "",
			wantErr: true,
			errStr:  errors.Join(service.ErrGeneratingToken, errors.New("")),
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
			token, err := s.SignIn(&tt.args.candidate)
			if tt.wantErr {
				require.Equal(t, tt.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
				require.Equal(t, token, tt.want)
			}
		})
	}
}
