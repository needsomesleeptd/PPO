package service_test

import (
	service "annotater/internal/bl/userService"
	mock_repository "annotater/internal/mocks/bl/userService/userRepo"
	"annotater/internal/models"
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
		role  models.Role
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
			name: "Changing tole no err",
			prepare: func(f *fields) {
				f.userRepo.EXPECT().GetUserByLogin(TEST_VALID_LOGIN).Return(&models.User{Role: models.Admin}, nil)
				f.userRepo.EXPECT().UpdateUserByLogin(TEST_VALID_LOGIN, &models.User{Role: models.Controller}).Return(nil)
			},
			args:    args{login: TEST_VALID_LOGIN, role: models.Controller},
			wantErr: false,
			errStr:  nil,
		},
		{
			name: "Changing role getting err",
			prepare: func(f *fields) {
				f.userRepo.EXPECT().GetUserByLogin(TEST_VALID_LOGIN).Return(nil, errors.New(""))
			},
			args:    args{login: TEST_VALID_LOGIN, role: models.Admin},
			wantErr: true,
			errStr:  errors.New(service.ERROR_CHANGE_ROLE_STR + ": "),
		},
		{
			name: "Changing role update err",
			prepare: func(f *fields) {
				f.userRepo.EXPECT().GetUserByLogin(TEST_VALID_LOGIN).Return(&models.User{Role: models.Sender}, nil)
				f.userRepo.EXPECT().UpdateUserByLogin(TEST_VALID_LOGIN, &models.User{Role: models.Controller}).Return(errors.New(""))
			},
			args:    args{login: TEST_VALID_LOGIN, role: models.Controller},
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
			err := s.ChangeUserRoleByLogin(tt.args.login, tt.args.role)
			if tt.wantErr {
				require.Equal(t, tt.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
			}

		})
	}
}
