package service_test

import (
	service "annotater/internal/bl/userService"
	mock_repository "annotater/internal/mocks/bl/userService/userRepo"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var TEST_VALID_LOGIN string = "valid"
var TEST_INVALID_LOGIN string = "invalid"

func TestUserService_ChangeUserRoleByLogin(t *testing.T) {
	type fields struct {
		userRepo *mock_repository.MockIUserRepository
	}
	type args struct {
		login string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		prepare func(f *fields)
		errStr  error
	}{
		{
			name: "Update user no err",
			prepare: func(f *fields) {
				f.userRepo.EXPECT().UpdateUserByLogin(TEST_VALID_LOGIN).Return(nil)
			},
			args:    args{login: TEST_VALID_LOGIN},
			wantErr: false,
			errStr:  nil,
		},
		{
			name: "Update user err",
			prepare: func(f *fields) {
				f.userRepo.EXPECT().UpdateUserByLogin(TEST_VALID_LOGIN).Return(errors.New(""))
			},
			args:    args{login: TEST_VALID_LOGIN},
			wantErr: true,
			errStr:  errors.New(service.ERROR_CHANGE_ROLE_STR + ": "),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				userRepo: mock_repository.NewMockIUserRepository(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			s := service.NewUserService(f.userRepo)
			err := s.ChangeUserRoleByLogin(tt.args.login)
			if tt.wantErr {
				require.Equal(t, tt.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
			}

		})
	}
}
