package service_test

import (
	service "annotater/internal/bl/anotattionTypeService"
	mock_repository "annotater/internal/mocks/bl/anotattionTypeService/anottationTypeRepo"
	"annotater/internal/models"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

const TEST_BASIC_ID uint64 = 20

func TestAnotattionTypeService_AddAnottationType(t *testing.T) {
	type fields struct {
		repo *mock_repository.MockIAnotattionTypeRepository
	}
	type args struct {
		anotattionType *models.MarkupType
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		prepare func(f *fields)
		wantErr bool
		errStr  error
	}{
		{
			name: "Add no err",
			prepare: func(f *fields) {
				f.repo.EXPECT().AddAnottationType(&models.MarkupType{}).Return(nil)
			},
			wantErr: false,
			errStr:  nil,
			args:    args{anotattionType: &models.MarkupType{}},
		},
		{
			name: "Add with err",
			prepare: func(f *fields) {
				f.repo.EXPECT().AddAnottationType(&models.MarkupType{}).Return(errors.New(""))
			},
			wantErr: true,
			args:    args{anotattionType: &models.MarkupType{}},
			errStr:  errors.New(service.ADDING_ANNOTATTION_ERR_STR + ": "),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				repo: mock_repository.NewMockIAnotattionTypeRepository(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			s := service.NewAnotattionTypeService(f.repo)
			err := s.AddAnottationType(tt.args.anotattionType)
			if tt.wantErr {
				require.Equal(t, tt.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestAnotattionTypeService_DeleteAnotattionType(t *testing.T) {
	type fields struct {
		repo *mock_repository.MockIAnotattionTypeRepository
	}
	type args struct {
		id uint64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		prepare func(f *fields)
		wantErr bool
		errStr  error
	}{
		{
			name: "Delete no err",
			prepare: func(f *fields) {
				f.repo.EXPECT().DeleteAnotattionType(TEST_BASIC_ID).Return(nil)
			},
			wantErr: false,
			errStr:  nil,
			args:    args{id: TEST_BASIC_ID},
		},
		{
			name: "Delete with err",
			prepare: func(f *fields) {
				f.repo.EXPECT().DeleteAnotattionType(TEST_BASIC_ID).Return(errors.New(""))
			},
			wantErr: true,
			args:    args{id: TEST_BASIC_ID},
			errStr:  errors.New(service.DELETING_ANNOTATTION_ERR_STR + ": "),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				repo: mock_repository.NewMockIAnotattionTypeRepository(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			s := service.NewAnotattionTypeService(f.repo)
			err := s.DeleteAnotattionType(tt.args.id)
			if tt.wantErr {
				require.Equal(t, tt.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestAnotattionTypeService_GetAnottationTypeByID(t *testing.T) {
	type fields struct {
		repo *mock_repository.MockIAnotattionTypeRepository
	}
	type args struct {
		id uint64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		prepare func(f *fields)
		wantErr bool
		errStr  error
		want    *models.MarkupType
	}{
		{
			name: "Get no err",
			prepare: func(f *fields) {
				f.repo.EXPECT().GetAnottationTypeByID(TEST_BASIC_ID).Return(&models.MarkupType{ID: TEST_BASIC_ID}, nil)
			},
			wantErr: false,
			errStr:  nil,
			args:    args{id: TEST_BASIC_ID},
			want:    &models.MarkupType{ID: TEST_BASIC_ID},
		},
		{
			name: "Get with err",
			prepare: func(f *fields) {
				f.repo.EXPECT().GetAnottationTypeByID(TEST_BASIC_ID).Return(nil, errors.New(""))
			},
			wantErr: true,
			args:    args{id: TEST_BASIC_ID},
			errStr:  errors.New(service.GETTING_ANNOTATTION_STR_ERR_STR + ": "),
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				repo: mock_repository.NewMockIAnotattionTypeRepository(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			s := service.NewAnotattionTypeService(f.repo)
			res, err := s.GetAnottationTypeByID(tt.args.id)
			if tt.wantErr {
				require.NotNil(t, err)
				require.Equal(t, tt.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
				require.Equal(t, res, tt.want)
			}
		})
	}
}
