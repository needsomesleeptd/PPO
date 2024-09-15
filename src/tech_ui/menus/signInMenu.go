package menus

import (
	auth_service "annotater/internal/bl/auth"
	response "annotater/internal/lib/api"
	"annotater/internal/models"
	auth_utils "annotater/internal/pkg/authUtils"
	auth_ui "annotater/tech_ui/utils/auth"
	"fmt"
	"log"

	"github.com/dixonwille/wmenu/v5"
)

func (m *Menu) SignInMenu(opt wmenu.Opt) error {
	client, ok := opt.Value.(ClientEntity)
	if !ok {
		log.Fatal("Could not cast option's value to ClientEntity")
	}
	var login string
	var passwd string
	fmt.Println("Enter login:")
	fmt.Scan(&login)
	fmt.Println("Enter password:")
	fmt.Scan(&passwd)
	jwt, err := auth_ui.SignIn(client.Client, login, passwd)
	if err != nil {
		return err
	}
	m.jwt = jwt
	payload, err := auth_utils.JWTTokenHandler{}.ParseToken(jwt, auth_service.SECRET)
	if err != nil {
		return err
	}
	m.ID = payload.ID
	m.role = payload.Role
	switch m.role {
	case models.Sender:
		m.RunUserMenu(client.Client)
	case models.Controller:
		m.RunControllerMenu(client.Client)
	case models.Admin:
		m.RunAdminMenu(client.Client)
	}
	fmt.Println(response.StatusOK)
	return nil
}

func (m *Menu) SignUpMenu(opt wmenu.Opt) error {
	client, ok := opt.Value.(ClientEntity)
	if !ok {
		log.Fatal("Could not cast option's value to ClientEntity")
	}
	var login string
	var passwd string
	fmt.Println("Enter login:")
	fmt.Scan(&login)
	fmt.Println("Enter password:")
	fmt.Scan(&passwd)
	_, err := auth_ui.SignUp(client.Client, login, passwd)

	if err != nil {
		return err
	}
	fmt.Println(response.StatusOK)
	return nil
}
