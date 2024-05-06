package menus

import (
	response "annotater/internal/lib/api"
	annot_type_req "annotater/tech_ui/utils/markupType"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/dixonwille/wmenu/v5"
	"github.com/olekukonko/tablewriter"
)

func (m *Menu) GettingAnotattionType(opt wmenu.Opt) error {
	clientEntity, ok := opt.Value.(ClientEntity)
	if !ok {
		log.Fatal("Could not cast option's value to ClientEntity")
	}

	markupTypes, err := annot_type_req.GetMarkupTypesCreatorID(clientEntity.Client, m.jwt)
	if err != nil {
		return err
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"id", "name", "description"})
	for _, markupType := range markupTypes {
		table.Append([]string{strconv.FormatUint(markupType.ID, 10), markupType.ClassName, markupType.Description})
	}
	table.Render()
	return nil
}

func (m *Menu) AddingAnotattionType(opt wmenu.Opt) error {
	clientEntity, ok := opt.Value.(ClientEntity)
	if !ok {
		log.Fatal("Could not cast option's value to ClientEntity")
	}
	var labelName string
	var description string
	var id uint64

	fmt.Println("Enter the ID of the anotattion type:")
	fmt.Scan(&id)

	fmt.Println("Enter the label name for the new type:")
	fmt.Scan(&labelName)

	fmt.Println("Enter the description of the new type:")
	fmt.Scan(&description)

	err := annot_type_req.AddMarkupTypeByCreatorID(clientEntity.Client, labelName, description, m.jwt, id)
	if err != nil {
		return err
	}
	fmt.Print(response.StatusOK)
	return nil
}

func (m *Menu) GettingAllAnottationTypes(opt wmenu.Opt) error {
	clientEntity, ok := opt.Value.(ClientEntity)
	if !ok {
		log.Fatal("Could not cast option's value to ClientEntity")
	}

	markupTypes, err := annot_type_req.GetAllMarkupTypes(clientEntity.Client, m.jwt)
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"id", "name", "description"})
	for _, markupType := range markupTypes {
		table.Append([]string{strconv.FormatUint(markupType.ID, 10), markupType.ClassName, markupType.Description})
	}
	table.Render()
	return nil
}
