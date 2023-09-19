package user

import (
	"demerzel-events/internal/db"
	"demerzel-events/internal/models"
	"errors"
	"fmt"
	"net/http"
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

func LoginUser(userLoginObject models.UserLogin) (models.UserResponse, string, int, error) {
	_, err := getUserFromDB(userLoginObject.Email)
	if err != nil {
		return models.UserResponse{}, "user does not exist", 404, err
	}

	// logic to login user

	return models.UserResponse{}, "", 0, nil
}

func ResetPassword(reqBody models.PasswordReset) (int, error) {
	_, err := getUserFromDB(reqBody.Email)
	if err != nil {
		return 404, fmt.Errorf("user does not exist: %s", err.Error())
	}

	// logic to reset password

	return 0, nil
}

func ForgotPassword(reqBody models.ForgotPassword) (int, error) {
	_, err := getUserFromDB(reqBody.Email)
	if err != nil {
		return 404, fmt.Errorf("user does not exist: %s", err.Error())
	}

	// logic to reset password

	return http.StatusOK, nil
}

func getUserFromDB(email string) (models.User, error) {
	
	// get user from db
	var user models.User
	err:=db.DB.Preload("Events").Where("Email=?",email).Find(&user).Error
	if err != nil {
        return models.User{}, err
    }
	return user, nil
}


func UpdateUserService(user models.UpdateUser, userId interface{}) (int, error) {
	return 0, nil
}


func CheckPasswords() error {
	
	// logic check password

	return errors.New("")
}
