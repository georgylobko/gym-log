package mappers

import "github.com/georgylobko/gym-log/internal/database"

type User struct {
	ID     int32  `json:"id"`
	Name   string `json:"name"`
	Gender string `json:"gender"`
	Role   string `json:"role"`
	Email  string `json:"email"`
}

func DatabaseUserToUser(dbUser database.User) User {
	return User{
		ID:     dbUser.ID,
		Name:   dbUser.Name,
		Gender: dbUser.Gender.String,
		Role:   dbUser.Role,
		Email:  dbUser.Email,
	}
}
