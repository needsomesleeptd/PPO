package menus

import (
	response "annotater/internal/lib/api"
	"annotater/tech_ui/utils/document_req"
	document_ui "annotater/tech_ui/utils/document_req"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/dixonwille/wmenu/v5"
	"github.com/google/uuid"
	"github.com/olekukonko/tablewriter"
)

func (m *Menu) checkingDocumentOpt(opt wmenu.Opt) error {
	clientEntity, ok := opt.Value.(ClientEntity)
	if !ok {
		log.Fatal("Could not cast option's value to ClientEntity")
	}
	var filePath string
	var outputFolderPath string
	fmt.Println("Enter path to the pdf document to check:")
	fmt.Scan(&filePath)

	fmt.Println("Enter the path to the folder for the output:")
	fmt.Scan(&outputFolderPath)

	err := document_req.ReportDocument(clientEntity.Client, filePath, outputFolderPath, m.jwt)
	if err != nil {
		return err
	}

	fmt.Print(response.StatusOK)
	return nil
}

func (m *Menu) GettingDocumentByID(opt wmenu.Opt) error {
	clientEntity, ok := opt.Value.(ClientEntity)
	if !ok {
		log.Fatal("Could not cast option's value to ClientEntity")
	}
	var filePath string
	var idStr string

	var err error
	var id uuid.UUID

	fmt.Println("Enter the id of the document to load:")
	fmt.Scan(&idStr)

	fmt.Println("Enter path to the output directory:")
	fmt.Scan(&filePath)

	id, err = uuid.Parse(idStr)
	if err != nil {
		return err
	}
	err = document_ui.GetDocument(clientEntity.Client, filePath, m.jwt, id)

	if err != nil {
		return err
	}

	fmt.Print(response.StatusOK)
	return nil
}

func (m *Menu) GettingReportByID(opt wmenu.Opt) error {
	clientEntity, ok := opt.Value.(ClientEntity)
	if !ok {
		log.Fatal("Could not cast option's value to ClientEntity")
	}
	var filePath string
	var idStr string

	var err error
	var id uuid.UUID

	fmt.Println("Enter the id of the report to load(report id is equal to document id):")
	fmt.Scan(&idStr)

	fmt.Println("Enter path to the output directory:")
	fmt.Scan(&filePath)

	id, err = uuid.Parse(idStr)
	if err != nil {
		return err
	}
	err = document_ui.GetReport(clientEntity.Client, filePath, m.jwt, id)

	if err != nil {
		return err
	}

	fmt.Print(response.StatusOK)
	return nil
}

func (m *Menu) GettingReportsMetaData(opt wmenu.Opt) error {
	clientEntity, ok := opt.Value.(ClientEntity)
	if !ok {
		log.Fatal("Could not cast option's value to ClientEntity")
	}

	resp, err := document_ui.GetDocumentsMetaData(clientEntity.Client, m.jwt)
	if err != nil {
		return err
	}
	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader([]string{"id", "page_cnt", "name", "creator_id", "creation_time"})
	for _, documentmeta := range resp.DocumentsMetaData {
		table.Append([]string{documentmeta.ID.String(),
			strconv.FormatInt(int64(documentmeta.PageCount), 10),
			documentmeta.DocumentName,
			strconv.FormatInt(int64(documentmeta.CreatorID), 10),
			documentmeta.CreationTime.String(),
		})
	}
	table.Render()
	return nil
}
