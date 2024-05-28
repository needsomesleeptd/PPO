package rep_data_repo

import (
	"annotater/internal/models"

	"github.com/google/uuid"
)

type IReportDataRepository interface {
	AddReport(doc *models.ErrorReport) error
	DeleteReportByID(id uuid.UUID) error
	GetDocumentByID(id uuid.UUID) (*models.ErrorReport, error)
}
