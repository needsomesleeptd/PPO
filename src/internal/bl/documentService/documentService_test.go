package service_test

import (
	service "annotater/internal/bl/documentService"
	mock_nn "annotater/internal/mocks/bl/NN"
	mock_repository "annotater/internal/mocks/bl/documentService/documentRepo"
	"annotater/internal/models"
	"bytes"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/signintech/gopdf"
	"github.com/stretchr/testify/require"
)

var TEST_VALID_PDF *gopdf.GoPdf = &gopdf.GoPdf{}

func createPDFBuffer(pdf *gopdf.GoPdf) []byte {
	if pdf == nil {
		return []byte{1}
	}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	var buf bytes.Buffer
	pdf.WriteTo(&buf)

	return buf.Bytes()
}

func TestDocumentService_LoadDocument(t *testing.T) {
	type fields struct {
		repo          *mock_repository.MockIDocumentRepository
		neuralNetwork *mock_nn.MockINeuralNetwork
	}
	type args struct {
		document models.DocumentMetaData
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
			name: "Successful Load Document",
			prepare: func(f *fields) {
				f.repo.EXPECT().AddDocument(&models.DocumentMetaData{
					DocumentData: createPDFBuffer(TEST_VALID_PDF),
				}).Return(nil)
			},
			args:    args{document: models.DocumentMetaData{DocumentData: createPDFBuffer(TEST_VALID_PDF)}},
			wantErr: false,
			errStr:  nil,
		},
		{
			name: "Impossible to add to Repo",
			prepare: func(f *fields) {
				f.repo.EXPECT().AddDocument(&models.DocumentMetaData{
					DocumentData: createPDFBuffer(TEST_VALID_PDF),
				}).Return(errors.Errorf("%s", ""))
			},
			args:    args{document: models.DocumentMetaData{DocumentData: createPDFBuffer(TEST_VALID_PDF)}},
			wantErr: true,
			errStr:  errors.New(service.LOADING_DOCUMENT_ERR_STRF + ": "),
		},
		{
			name:    "invalid file format",
			args:    args{document: models.DocumentMetaData{DocumentData: createPDFBuffer(nil)}},
			wantErr: true,
			errStr:  errors.New(service.DOCUMENT_FORMAT_ERR_STR),
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
		document models.DocumentMetaData
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.Markup
		wantErr bool
		prepare func(f *fields)
		errStr  error
	}{
		{
			name: "Valid Check one markup",
			prepare: func(f *fields) {
				f.neuralNetwork.EXPECT().Predict(models.DocumentMetaData{
					DocumentData: createPDFBuffer(TEST_VALID_PDF)}).Return([]models.Markup{
					{ErrorBB: []float32{0.4, 0.3, 0.2},
						ClassLabel: 1,
					},
				}, nil)
			},
			args: args{document: models.DocumentMetaData{
				DocumentData: createPDFBuffer(TEST_VALID_PDF)}},
			wantErr: false,
			errStr:  nil,
			want: []models.Markup{
				{ErrorBB: []float32{0.4, 0.3, 0.2},
					ClassLabel: 1,
				},
			},
		},
		{
			name: "Valid check zero markups",
			prepare: func(f *fields) {
				f.neuralNetwork.EXPECT().Predict(models.DocumentMetaData{
					DocumentData: createPDFBuffer(TEST_VALID_PDF)}).Return(nil, nil)
			},
			args: args{document: models.DocumentMetaData{
				DocumentData: createPDFBuffer(TEST_VALID_PDF)}},
			wantErr: false,
			errStr:  nil,
			want:    nil,
		},
		{
			name: "Check error",
			prepare: func(f *fields) {
				f.neuralNetwork.EXPECT().Predict(models.DocumentMetaData{
					DocumentData: createPDFBuffer(TEST_VALID_PDF)}).Return(nil, errors.New(""))
			},
			args:    args{document: models.DocumentMetaData{DocumentData: createPDFBuffer(TEST_VALID_PDF)}},
			wantErr: true,
			errStr:  errors.New(service.CHECKING_DOCUMENT_ERR_STRF + ": "),
			want:    nil,
		},
		{
			name: "Valid check numerous markups",
			prepare: func(f *fields) {
				f.neuralNetwork.EXPECT().Predict(models.DocumentMetaData{DocumentData: createPDFBuffer(TEST_VALID_PDF)}).Return([]models.Markup{
					{ErrorBB: []float32{0.4, 0.3, 0.2},
						ClassLabel: 1,
					},
					{ErrorBB: []float32{0.1, 0.2, 0.1},
						ClassLabel: 2,
					},
				}, nil)
			},
			args:    args{document: models.DocumentMetaData{DocumentData: createPDFBuffer(TEST_VALID_PDF)}},
			wantErr: false,
			errStr:  nil,
			want: []models.Markup{
				{ErrorBB: []float32{0.4, 0.3, 0.2},
					ClassLabel: 1,
				},
				{ErrorBB: []float32{0.1, 0.2, 0.1},
					ClassLabel: 2,
				},
			},
		},
		{
			name: "Invalid File Format",
			args: args{document: models.DocumentMetaData{
				DocumentData: createPDFBuffer(nil)}},
			wantErr: true,
			errStr:  errors.New(service.DOCUMENT_FORMAT_ERR_STR),
			want:    nil,
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
