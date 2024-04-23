package menus

import (
	response "annotater/internal/lib/api"
	markup_req "annotater/tech_ui/utils/markup"
	markup_view "annotater/tech_ui/view/markup"
	"fmt"
	"log"

	"github.com/dixonwille/wmenu/v5"
)

func (m *Menu) AddingAnotattion(opt wmenu.Opt) error {
	clientEntity, ok := opt.Value.(ClientEntity)
	if !ok {
		log.Fatal("Could not cast option's value to ClientEntity")
	}
	var filePath string
	fmt.Println("Enter path to the png image to load")
	fmt.Scan(&filePath)

	bbs := make([]float32, 4)

	fmt.Print("Enter the coords os the bounding boxes(xyxy format):")
	for i := range bbs {
		fmt.Scan(&bbs[i])
	}

	var classLabel uint64
	fmt.Print("Enter the classLabel:")
	fmt.Scan(&classLabel)

	err := markup_req.AddMarkup(clientEntity.Client, filePath, bbs, classLabel, m.jwt)
	if err != nil {
		return err
	}
	fmt.Print(response.StatusOK)
	return nil
}

func (m *Menu) DeletingAnotattion(opt wmenu.Opt) error {
	clientEntity, ok := opt.Value.(ClientEntity)
	if !ok {
		log.Fatal("Could not cast option's value to ClientEntity")
	}
	var markupID uint64
	fmt.Println("Enter the ID of the markup to delete")
	fmt.Scan(&markupID)

	err := markup_req.DeletingMarkups(clientEntity.Client, markupID, m.jwt)
	if err != nil {
		return err
	}
	fmt.Print(response.StatusOK)
	return nil
}

func (m *Menu) GettingAnotattion(opt wmenu.Opt) error {
	clientEntity, ok := opt.Value.(ClientEntity)
	if !ok {
		log.Fatal("Could not cast option's value to ClientEntity")
	}
	var filePath string
	fmt.Println("Enter path to the folder where results will be stored")
	fmt.Scan(&filePath)

	resp, err := markup_req.GetYourMarkups(clientEntity.Client, m.ID, m.jwt)
	if err != nil {
		return err
	}
	err = markup_view.DrawBbsOnMarkups(resp, filePath)

	if err != nil {
		return err
	}

	fmt.Print(response.StatusOK)
	return nil
}
