package service_test

import (
	service "annotater/internal/bl/documentService"
	mock_nn "annotater/internal/mocks/bl/NN"
	mock_repository "annotater/internal/mocks/bl/documentService/documentRepo"
	"annotater/internal/models"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestDocumentService_LoadDocument(t *testing.T) {
	type fields struct {
		repo          *mock_repository.MockIDocumentRepository
		neuralNetwork *mock_nn.MockINeuralNetwork
	}
	type args struct {
		document models.Document
	}
	tests := []struct {
		name    string
		fields  fields
		prepare func(f *fields)
		args    args
		wantErr bool
		errStr  error
	}{
		{
			name: "SuccessFul Load Document",
			prepare: func(f *fields) {
				f.repo.EXPECT().AddDocument(&models.Document{}).Return(nil)
			},
			args:    args{document: models.Document{ChecksCount: 0}},
			wantErr: false,
			errStr:  nil,
		},
		{
			name: "Impossible to add to Repo",
			prepare: func(f *fields) {
				f.repo.EXPECT().AddDocument(&models.Document{}).Return(errors.Errorf("%s", ""))
			},
			args:    args{document: models.Document{ChecksCount: 0}},
			wantErr: true,
			errStr:  errors.New("Error in loading document: "),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				repo:          mock_repository.NewMockIDocumentRepository(ctrl),
				neuralNetwork: mock_nn.NewMockINeuralNetwork(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			s := service.NewDocumentService(f.repo, f.neuralNetwork)
			err := s.LoadDocument(tt.args.document)
			if tt.wantErr {
				require.Equal(t, tt.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
			}

		})
	}
}

func TestDocumentService_CheckDocument(t *testing.T) {
	type fields struct {
		repo          *mock_repository.MockIDocumentRepository
		neuralNetwork *mock_nn.MockINeuralNetwork
	}
	type args struct {
		document models.Document
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*models.Markup
		wantErr bool
		prepare func(f *fields)
		errStr  error
	}{
		{
			name: "Valid Check one markup",
			prepare: func(f *fields) {
				f.neuralNetwork.EXPECT().Predict(models.Document{}).Return([]*models.Markup{
					{ErrorBB: []float32{0.4, 0.3, 0.2},
						ClassLabel: 1,
					},
				}, nil)
			},
			args:    args{document: models.Document{}},
			wantErr: false,
			errStr:  nil,
			want: []*models.Markup{
				{ErrorBB: []float32{0.4, 0.3, 0.2},
					ClassLabel: 1,
				},
			},
		},
		{
			name: "Valid check zero markups",
			prepare: func(f *fields) {
				f.neuralNetwork.EXPECT().Predict(models.Document{}).Return(nil, nil)
			},
			args:    args{document: models.Document{}},
			wantErr: false,
			errStr:  nil,
			want:    nil,
		},
		{
			name: "Check error",
			prepare: func(f *fields) {
				f.neuralNetwork.EXPECT().Predict(models.Document{}).Return(nil, errors.New(""))
			},
			args:    args{document: models.Document{}},
			wantErr: true,
			errStr:  errors.New(service.ERROR_CHECKING_DOCUMENT + ": "),
			want:    nil,
		},
		{
			name: "Valid check numerous markups",
			prepare: func(f *fields) {
				f.neuralNetwork.EXPECT().Predict(models.Document{}).Return([]*models.Markup{
					{ErrorBB: []float32{0.4, 0.3, 0.2},
						ClassLabel: 1,
					},
					{ErrorBB: []float32{0.1, 0.2, 0.1},
						ClassLabel: 2,
					},
				}, nil)
			},
			args:    args{document: models.Document{}},
			wantErr: false,
			errStr:  nil,
			want: []*models.Markup{
				{ErrorBB: []float32{0.4, 0.3, 0.2},
					ClassLabel: 1,
				},
				{ErrorBB: []float32{0.1, 0.2, 0.1},
					ClassLabel: 2,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				repo:          mock_repository.NewMockIDocumentRepository(ctrl),
				neuralNetwork: mock_nn.NewMockINeuralNetwork(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			s := service.NewDocumentService(f.repo, f.neuralNetwork)
			markups, err := s.CheckDocument(tt.args.document)
			if tt.wantErr {
				require.Equal(t, tt.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
				require.Equal(t, markups, tt.want)
			}

		})
	}
}
