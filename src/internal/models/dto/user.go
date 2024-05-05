package models_dto // stands for data_transfer_objects

import (
	"annotater/internal/models"
	"encoding/json"
)

type User struct {
	ID       uint64      `json:"id"`
	Login    string      `json:"login"`
	Password string      `json:"password"`
	Name     string      `json:"name"`
	Surname  string      `json:"surname"`
	Role     models.Role `json:"role"`
	Group    string      `json:"group"` // in case it is a controller it will have work group, in case of user, his group
}

func (u *User) ToJSON() (string, error) {
	jsonBytes, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

func FromDtoUser(userDa *User) models.User {
	return models.User{
		ID:       userDa.ID,
		Name:     userDa.Name,
		Login:    userDa.Login,
		Password: userDa.Password,
		Surname:  userDa.Surname,
		Role:     userDa.Role,
		Group:    userDa.Group,
	}
}

func ToDtoUser(user models.User) *User {
	return &User{
		ID:       user.ID,
		Name:     user.Name,
		Login:    user.Login,
		Password: user.Password,
		Surname:  user.Surname,
		Role:     user.Role,
		Group:    user.Group,
	}
}

func ToDtoUserSlice(users []models.User) []User {
	usersDTO := make([]User, len(users))
	for i, user := range users {
		usersDTO[i] = *ToDtoUser(user)
	}
	return usersDTO
}
