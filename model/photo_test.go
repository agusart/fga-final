package model

import (
	"github.com/asaskevich/govalidator"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestPhotoValidationError(t *testing.T) {
	errMsg := []string{
		"photo url is empty",
		"title is empty",
	}

	photo := Photo{}

	valid, err := govalidator.ValidateStruct(photo)
	assert.False(t, valid)
	assert.Error(t, err)
	assert.Equal(t, errMsg, strings.Split(err.Error(), ";"))
}

func TestPhotoValidationSuccess(t *testing.T) {
	photo := Photo{
		PhotoURL: "http://test.com",
		Title:    "test title",
	}

	valid, err := govalidator.ValidateStruct(photo)
	assert.True(t, valid)
	assert.NoError(t, err)
}
