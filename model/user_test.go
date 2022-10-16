package model

import (
	"github.com/asaskevich/govalidator"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestUserEmptyInput(t *testing.T) {
	errMsg := []string{
		"age is empty",
		"email is blank",
		"password is empty",
		"username is empty",
	}

	user := User{}

	valid, err := govalidator.ValidateStruct(user)
	assert.False(t, valid)
	assert.Error(t, err)
	assert.Equal(t, errMsg, strings.Split(err.Error(), ";"))
}

func TestUserInvalidInput(t *testing.T) {
	errMsg := []string{
		"age is must between 8-125",
		"email is invalid",
		"minimum password must 6 char length",
		"username is empty",
	}

	user := User{
		Email:    "asdasd",
		Password: "1",
		Age:      2,
		Username: "",
	}

	valid, err := govalidator.ValidateStruct(user)
	assert.False(t, valid)
	assert.Error(t, err)
	assert.Equal(t, errMsg, strings.Split(err.Error(), ";"))
}

func TestUserValidInput(t *testing.T) {

	user := User{
		Email:    "test@gmail.com",
		Password: "123456",
		Age:      8,
		Username: "aizakami",
	}

	valid, err := govalidator.ValidateStruct(user)
	assert.True(t, valid)
	assert.NoError(t, err)
}
