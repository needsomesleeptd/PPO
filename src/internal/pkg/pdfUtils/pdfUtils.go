package pdf_utils

import (
	"bytes"

	"github.com/unidoc/unipdf/v3/processor"
)

func GetPdfPageCount(pdfData []byte) (int, error) {
	pdf, err := processor.NewPdfReader(bytes.NewReader(pdfData))
	if err != nil {
		return -1, err
	}
	defer pdf.Close()

	pageCount, err := pdf.GetNumPages()
	if err != nil {
		return -1, err
	}
	return pageCount, err
}
