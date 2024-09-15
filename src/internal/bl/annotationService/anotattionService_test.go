package service_test

import (
	service "annotater/internal/bl/annotationService"
	mock_repository "annotater/internal/mocks/bl/annotationService/annotattionRepo"
	"annotater/internal/models"
	"bytes"
	"image"
	"image/png"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

const (
	TEST_BASIC_ID uint64 = 20
)

var TEST_MARKUP models.Markup = models.Markup{
	ErrorBB:    []float32{0.0, 1.0, 0.0, 1.0},
	ClassLabel: 1,
}

var TEST_VALID_PNG_IMG *image.RGBA = image.NewRGBA(image.Rect(0, 0, 100, 100))

func createPNGBuffer(img *image.RGBA) []byte {
	if img == nil {
		return nil
	}
	pngBuf := new(bytes.Buffer)
	png.Encode(pngBuf, img)
	return pngBuf.Bytes()
}

func TestAnottationService_AddAnnotation(t *testing.T) {
	type fields struct {
		repo *mock_repository.MockIAnotattionRepository
	}
	type args struct {
		annotation *models.Markup
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
			name: "Add no error",
			prepare: func(f *fields) {
				f.repo.EXPECT().AddAnottation(&models.Markup{
					ErrorBB: []float32{
						1.0, 1.0, 1.0, 1.0,
					},
					PageData: createPNGBuffer(TEST_VALID_PNG_IMG)}).Return(nil)
			},
			wantErr: false,
			errStr:  nil,
			args: args{annotation: &models.Markup{
				ErrorBB: []float32{
					1.0, 1.0, 1.0, 1.0,
				},
				PageData: createPNGBuffer(TEST_VALID_PNG_IMG),
			}},
		},
		{
			name: "Add with repository error",
			prepare: func(f *fields) {
				f.repo.EXPECT().AddAnottation(&models.Markup{PageData: createPNGBuffer(TEST_VALID_PNG_IMG)}).Return(errors.New(""))
			},
			wantErr: true,
			args:    args{annotation: &models.Markup{PageData: createPNGBuffer(TEST_VALID_PNG_IMG)}},
			errStr:  errors.New(service.ADDING_ANNOT_ERR_STR + ": "),
		},
		{
			name:    "Add with invalid markup BBs",
			wantErr: true,
			args: args{annotation: &models.Markup{
				ErrorBB: []float32{
					-1.0, 1.0, 1.0, 1.0,
				},
				PageData: createPNGBuffer(TEST_VALID_PNG_IMG),
			}},
			errStr: errors.New(service.INVALID_BBS_ERR_STR),
		},
		{
			name:    "Add with invalid page",
			wantErr: true,
			args: args{annotation: &models.Markup{
				ErrorBB: []float32{
					1.0, 1.0, 1.0, 1.0,
				},
				PageData: createPNGBuffer(nil),
			}},
			errStr: errors.New(service.INVALID_FILE_ERR_STR + ": image: unknown format"),
		},
	}
	image.RegisterFormat("png", "\x89PNG\r\n\x1a\n", png.Decode, png.DecodeConfig)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				repo: mock_repository.NewMockIAnotattionRepository(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			s := service.NewAnnotattionService(f.repo)
			err := s.AddAnottation(tt.args.annotation)
			if tt.wantErr {
				require.Equal(t, tt.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func Test_areBBsValid(t *testing.T) {
	type args struct {
		slice []float32
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid slice",
			args: args{slice: []float32{
				1.0, 0.0, 0.0, 1.0,
			}},
			want: true,
		},
		{
			name: "invalid neg slice",
			args: args{slice: []float32{
				-1.0, 0.0, 0.0, 1.0,
			}},
			want: false,
		},
		{
			name: "invalid bigger than 1 slice",
			args: args{slice: []float32{
				1.0, 0.0, 0.0, 1.1,
			}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := service.AreBBsValid(tt.args.slice); got != tt.want {
				t.Errorf("areBBsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnotattionService_DeleteAnotattion(t *testing.T) {
	type fields struct {
		repo *mock_repository.MockIAnotattionRepository
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
			name: "Delete no error",
			prepare: func(f *fields) {
				f.repo.EXPECT().DeleteAnotattion(TEST_BASIC_ID).Return(nil)
			},
			wantErr: false,
			errStr:  nil,
			args:    args{id: TEST_BASIC_ID},
		},
		{
			name: "Delete with repository error",
			prepare: func(f *fields) {
				f.repo.EXPECT().DeleteAnotattion(TEST_BASIC_ID).Return(errors.New(""))
			},
			wantErr: true,
			args:    args{id: TEST_BASIC_ID},
			errStr:  errors.New(service.DELETING_ANNOT_ERR_STR + ": "),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				repo: mock_repository.NewMockIAnotattionRepository(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			s := service.NewAnnotattionService(f.repo)
			err := s.DeleteAnotattion(tt.args.id)
			if tt.wantErr {
				require.Equal(t, tt.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestAnotattionService_GetAnottationByID(t *testing.T) {
	type fields struct {
		repo *mock_repository.MockIAnotattionRepository
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
		want    *models.Markup
	}{
		{
			name: "Get no error",
			prepare: func(f *fields) {
				f.repo.EXPECT().GetAnottationByID(TEST_BASIC_ID).Return(&TEST_MARKUP, nil)
			},
			wantErr: false,
			errStr:  nil,
			args:    args{id: TEST_BASIC_ID},
			want:    &TEST_MARKUP,
		},
		{
			name: "Get with repository error",
			prepare: func(f *fields) {
				f.repo.EXPECT().GetAnottationByID(TEST_BASIC_ID).Return(nil, errors.New(""))
			},
			wantErr: true,
			args:    args{id: TEST_BASIC_ID},
			errStr:  errors.New(service.GETTING_ANNOT_ERR_STR + ": "),
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				repo: mock_repository.NewMockIAnotattionRepository(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			s := service.NewAnnotattionService(f.repo)
			markup, err := s.GetAnottationByID(tt.args.id)
			if tt.wantErr {
				require.Equal(t, tt.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
				require.Equal(t, markup, tt.want)
			}
		})
	}
}

func Test_checkPngFile(t *testing.T) {
	type args struct {
		pngFile []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			args:    args{pngFile: createPNGBuffer(TEST_VALID_PNG_IMG)},
			wantErr: false,
		},
		{
			args:    args{pngFile: createPNGBuffer(nil)},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := service.CheckPngFile(tt.args.pngFile); (err != nil) != tt.wantErr {
				t.Errorf("checkPngFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
