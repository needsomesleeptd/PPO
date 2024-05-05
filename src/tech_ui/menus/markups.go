package menus

import (
	response "annotater/internal/lib/api"
	markup_req "annotater/tech_ui/utils/markup"
	markup_view "annotater/tech_ui/view/markup"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/dixonwille/wmenu/v5"
	"github.com/olekukonko/tablewriter"
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

func convertFloatSlicetoString(slice []float32) string {
	strSlice := make([]string, len(slice))
	for i, val := range slice {
		str := strconv.FormatFloat(float64(val), 'f', 2, 32) //use 32 bits and all precision
		strSlice[i] = str
	}
	return strings.Join(strSlice, ", ")
}

func (m *Menu) GettingAnotattionsByUserID(opt wmenu.Opt) error {
	clientEntity, ok := opt.Value.(ClientEntity)
	if !ok {
		log.Fatal("Could not cast option's value to ClientEntity")
	}
	var filePath string
	fmt.Println("Enter path to the folder where results will be stored")
	fmt.Scan(&filePath)

	resp, err := markup_req.GetMarkupsByID(clientEntity.Client, m.ID, m.jwt)
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"id", "error_bb", "class_label", "creator_id"})
	for _, markup := range resp.Markups {
		table.Append([]string{
			strconv.FormatUint(markup.ID, 10),
			convertFloatSlicetoString(markup.ErrorBB),
			strconv.FormatUint(markup.ClassLabel, 10),
			strconv.FormatUint(markup.CreatorID, 10),
		})
	}
	table.Render()

	err = markup_view.DrawBbsOnMarkups(resp, filePath)
	if err != nil {
		return err
	}

	fmt.Print(response.StatusOK)
	return nil
}

func (m *Menu) GettingAllAnottations(opt wmenu.Opt) error {
	clientEntity, ok := opt.Value.(ClientEntity)
	if !ok {
		log.Fatal("Could not cast option's value to ClientEntity")
	}

	var filePath string
	fmt.Println("Enter path to the folder where results will be stored")
	fmt.Scan(&filePath)

	resp, err := markup_req.GetAllMarkups(clientEntity.Client, m.jwt)
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"id", "error_bb", "class_label", "creator_id"})
	for _, markup := range resp.Markups {
		table.Append([]string{
			strconv.FormatUint(markup.ID, 10),
			convertFloatSlicetoString(markup.ErrorBB),
			strconv.FormatUint(markup.ClassLabel, 10),
			strconv.FormatUint(markup.CreatorID, 10),
		})
	}
	table.Render()

	err = markup_view.DrawBbsOnMarkups(resp, filePath)
	if err != nil {
		return err
	}

	return nil
}
