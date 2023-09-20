package user

import (
	"demerzel-events/internal/db"
	"demerzel-events/internal/models"
	"errors"
)

func SignUpUser(user models.User) (models.UserResponse, string, int, error) {
	// check if user already exists
	_, err := getUserFromDB(user.Email)
	if err == nil {
		return models.UserResponse{}, "user already exist", 403, errors.New("user already exist in database")
	}

	// logic to signup user

	return models.UserResponse{}, "successfully created user", 0, nil
}

// func LoginUser(userLoginObject models.UserLogin) (models.UserResponse, string, int, error) {
// 	_, err := getUserFromDB(userLoginObject.Email)
// 	if err != nil {
// 		return models.UserResponse{}, "user does not exist", 404, err
// 	}

// 	// logic to login user

// 	return models.UserResponse{}, "", 0, nil
// }

func getUserFromDB(email string) (models.User, error) {
	
	// get user from db
	var user models.User
	err:=db.DB.Preload("Events").Where("Email=?",email).Find(&user).Error
	if err != nil {
        return models.User{}, err
    }
	return user, nil
}