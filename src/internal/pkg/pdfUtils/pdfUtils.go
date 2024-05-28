package pdf_utils

import (
	"bytes"

	"github.com/unidoc/unipdf/v3/model"
)

func GetPdfPageCount(pdfData []byte) (int, error) {

	pdf, err := model.NewPdfReader(bytes.NewReader(pdfData))
	if err != nil {
		return -1, err
	}

	pageCount, err := pdf.GetNumPages()
	if err != nil {
		return -1, err
	}
	return pageCount, err
}
