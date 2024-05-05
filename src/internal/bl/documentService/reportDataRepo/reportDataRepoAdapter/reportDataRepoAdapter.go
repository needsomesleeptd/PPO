package rep_data_repo_adapter

import (
	rep_data_repo "annotater/internal/bl/documentService/reportDataRepo"
	"annotater/internal/models"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type ReportDataRepositoryAdapter struct {
	root          string
	fileExtension string
}

func NewDocumentRepositoryAdapter(rootSrc string, ext string) rep_data_repo.IReportDataRepository {
	return &ReportDataRepositoryAdapter{
		root:          rootSrc,
		fileExtension: ext,
	}
}

func (repo *ReportDataRepositoryAdapter) MakeDir(dir string) error {
	dirPath := fmt.Sprintf("%s/%s", repo.root, dir) + repo.fileExtension
	return os.MkdirAll(dirPath, 0755)
}

func (repo *ReportDataRepositoryAdapter) Exists(path string) bool {
	fullPath := fmt.Sprintf("%s/%s", repo.root, path) + repo.fileExtension
	_, err := os.Stat(fullPath)

	return !os.IsNotExist(err)
}

func (repo *ReportDataRepositoryAdapter) AddReport(rep *models.ErrorReport) error {
	if !repo.Exists(repo.root) {
		err := repo.MakeDir(repo.root)
		if err != nil {
			return errors.Wrap(err, "error in saving document data")
		}
	}

	filePath := fmt.Sprintf("%s/%s", repo.root, rep.DocumentID) + repo.fileExtension

	err := os.WriteFile(filePath, rep.ReportData, 0644)
	if err != nil {
		return errors.Wrap(err, "error in saving document data")
	}

	return nil
}

func (repo *ReportDataRepositoryAdapter) DeleteReportByID(id uuid.UUID) error {
	filePath := fmt.Sprintf("%s/%s", repo.root, id) + repo.fileExtension
	err := os.Remove(filePath)
	if err != nil {
		return errors.Wrap(err, "error in deleting document data")
	}

	return nil
}

func (repo *ReportDataRepositoryAdapter) GetDocumentByID(id uuid.UUID) (*models.ErrorReport, error) {
	filePath := fmt.Sprintf("%s/%s", repo.root, id) + repo.fileExtension
	fileBytes, err := os.ReadFile(filePath)

	if err == os.ErrNotExist {
		return nil, models.ErrNotFound
	}
	if err != nil {
		return nil, errors.Wrap(err, "error getting file")
	}

	report := &models.ErrorReport{
		ReportData: fileBytes,
		DocumentID: id,
	}
	return report, nil
}
