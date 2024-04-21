package menus

import (
	document_ui "annotater/tech_ui/utils/document_req"
	document_view "annotater/tech_ui/view/document"
	"fmt"
	"log"

	"github.com/dixonwille/wmenu/v5"
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

	fmt.Println("Enter the path to the folder fot the output:")
	fmt.Scan(&outputFolderPath)

	resp, err := document_ui.CheckDocument(clientEntity.Client, filePath, m.jwt)
	if err != nil {
		return err
	}

	res, err := document_view.GetCheckDocumentResult(resp, outputFolderPath)
	if err != nil {
		return err
	}
	fmt.Print(res)
	return nil
}

func (m *Menu) LoadingDocument(opt wmenu.Opt) error {
	clientEntity, ok := opt.Value.(ClientEntity)
	if !ok {
		log.Fatal("Could not cast option's value to ClientEntity")
	}
	var filePath string
	fmt.Println("Enter path to the pdf document to load:")
	fmt.Scan(&filePath)

	resp, err := document_ui.LoadDocument(clientEntity.Client, filePath, m.jwt)
	if err != nil {
		return err
	}
	res, err := document_view.GetLoadDocumentResult(resp)
	if err != nil {
		return err
	}
	fmt.Print(res)
	return nil
}
