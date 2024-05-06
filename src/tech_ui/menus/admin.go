package menus

import (
	response "annotater/internal/lib/api"
	"annotater/internal/models"
	annot_type_req "annotater/tech_ui/utils/markupType"
	role_req "annotater/tech_ui/utils/role"
	user_req "annotater/tech_ui/utils/user"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/dixonwille/wmenu/v5"
	"github.com/olekukonko/tablewriter"
)

func (m *Menu) RunAdminMenu(client *http.Client) error {
	m.aMenu = wmenu.NewMenu("Enter your option:")
	m.AddOptionsAdmin(client)

	for {
		err := m.aMenu.Run()
		fmt.Println()
		if err != nil {
			if errors.Is(err, errExit) {
				break
			}

			fmt.Printf("ERROR: %v\n\n", err)
		}
	}

	fmt.Printf("Exited menu.\n")
	return nil
}

func (m *Menu) ChangeUserRole(opt wmenu.Opt) error {
	clientEntity, ok := opt.Value.(ClientEntity)
	if !ok {
		log.Fatal("Could not cast option's value to ClientEntity")
	}
	var login string
	var role models.Role
	fmt.Println("Enter the login of the user to change perms")
	fmt.Scan(&login)

	fmt.Print("Enter the wanted role (0-sender,1-controller,2-admin)")
	fmt.Scan(&role)
	if role < 1 || role > 2 {
		return errors.New("invalid error number")
	}
	err := role_req.ChangeUserRole(clientEntity.Client, login, role, m.jwt)
	if err != nil {
		return err
	}

	return nil
}

func (m *Menu) GettingAllUsers(opt wmenu.Opt) error {
	clientEntity, ok := opt.Value.(ClientEntity)
	if !ok {
		log.Fatal("Could not cast option's value to ClientEntity")
	}

	users, err := user_req.GetAllUsers(clientEntity.Client, m.jwt)
	if err != nil {
		return err
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"id", "login", "role"})

	for _, user := range users {
		table.Append([]string{strconv.FormatUint(user.ID, 10), user.Login, user.Role.ToString()})
	}
	table.Render()

	return nil
}

func (m *Menu) DeletingAnotattionType(opt wmenu.Opt) error {
	clientEntity, ok := opt.Value.(ClientEntity)
	if !ok {
		log.Fatal("Could not cast option's value to ClientEntity")
	}

	var id uint64
	fmt.Println("Enter the id of the annotattion type to delete")
	fmt.Scan(&id)

	err := annot_type_req.DeleteMarkupType(clientEntity.Client, id, m.jwt)
	if err != nil {
		return err
	}
	fmt.Print(response.StatusOK)
	return nil
}
