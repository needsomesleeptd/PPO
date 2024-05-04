package models_da //stands for data_acess

import "annotater/internal/models"

type User struct {
	ID       uint64      `gorm:"primaryKey;column:id"`
	Login    string      `gorm:"unique;column:login"`
	Password string      `gorm:"column:password"`
	Name     string      `gorm:"column:name"`
	Surname  string      `gorm:"column:surname"`
	Role     models.Role `gorm:"embedded;column:role"`
	Group    string      `gorm:"embedded;column:group"` // in case it is a controller it will have work group, in case of user, his group
}

func FromDaUser(userDa *User) models.User {
	return models.User{
		ID:       userDa.ID,
		Name:     userDa.Name,
		Login:    userDa.Login,
		Password: userDa.Password,
		Surname:  userDa.Surname,
		Role:     models.Role(userDa.Role),
		Group:    userDa.Group,
	}

}

func FromDaUserSlice(usersDa []User) []models.User {

	users := make([]models.User, len(usersDa))

	for i, userDA := range usersDa {
		users[i] = FromDaUser(&userDA)
	}
	return users
}

func ToDaUser(user models.User) *User {
	return &User{
		ID:       user.ID,
		Name:     user.Name,
		Login:    user.Login,
		Password: user.Password,
		Surname:  user.Surname,
		Role:     models.Role(user.Role),
		Group:    user.Group,
	}
}
