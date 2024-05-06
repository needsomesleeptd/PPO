package nn_adapter_test

import (
	nn_adapter "annotater/internal/bl/NN/NNAdapter"
	mock_nn_model_handler "annotater/internal/mocks/bl/NN/NNAdapter/NNmodelhandler"
	"annotater/internal/models"
	models_dto "annotater/internal/models/dto"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

const SOME_DOCUMENT_DATA = "document data"

func TestDetectionModel_Predict(t *testing.T) {
	type fields struct {
		modelHandler *mock_nn_model_handler.MockIModelHandler
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
		err     error
	}{
		{
			name: "Valid Prediction",
			prepare: func(f *fields) {
				f.modelHandler.EXPECT().GetModelResp(gomock.Any()).Return([]models_dto.Markup{
					{ErrorBB: []float32{0.1, 0.2, 0.3}, ClassLabel: 1},
					{ErrorBB: []float32{0.3, 0.2, 0.1}, ClassLabel: 2},
				}, nil)
			},
			args: args{document: models.DocumentMetaData{DocumentData: []byte(SOME_DOCUMENT_DATA)}},
			want: []models.Markup{
				{ErrorBB: []float32{0.1, 0.2, 0.3}, ClassLabel: 1},
				{ErrorBB: []float32{0.3, 0.2, 0.1}, ClassLabel: 2},
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "Error in Model Response",
			prepare: func(f *fields) {
				f.modelHandler.EXPECT().GetModelResp(gomock.Any()).Return(nil, errors.New("error in model response"))
			},
			args:    args{document: models.DocumentMetaData{DocumentData: []byte(SOME_DOCUMENT_DATA)}},
			want:    nil,
			wantErr: true,
			err:     nn_adapter.ErrInModelPrediction,
		},
		{
			name: "Empty Prediction",
			prepare: func(f *fields) {
				f.modelHandler.EXPECT().GetModelResp(gomock.Any()).Return(nil, nil)
			},
			args:    args{document: models.DocumentMetaData{DocumentData: []byte(SOME_DOCUMENT_DATA)}},
			want:    nil,
			wantErr: false,
			err:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				modelHandler: mock_nn_model_handler.NewMockIModelHandler(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			s := nn_adapter.NewDetectionModel(f.modelHandler)
			markups, err := s.Predict(tt.args.document)
			if tt.wantErr {
				require.True(t, errors.Is(err, tt.err))
			} else {
				require.Nil(t, err)
				require.Equal(t, markups, tt.want)
			}

		})
	}
}
