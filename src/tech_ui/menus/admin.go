package menus

import (
	"annotater/internal/models"
	role_req "annotater/tech_ui/utils/role"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/dixonwille/wmenu/v5"
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

	err := role_req.ChangeUserRole(clientEntity.Client, login, role, m.jwt)
	if err != nil {
		return err
	}

	return nil
}
