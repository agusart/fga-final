package model

import (
	"github.com/asaskevich/govalidator"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestSocialMediaValidationError(t *testing.T) {
	errMsg := []string{
		"name is empty",
		"social media url is empty",
	}

	sosmed := SocialMedia{}

	valid, err := govalidator.ValidateStruct(sosmed)
	assert.False(t, valid)
	assert.Error(t, err)
	assert.Equal(t, errMsg, strings.Split(err.Error(), ";"))
}

func TestSosialeMediaValidationSuccess(t *testing.T) {
	sosmed := SocialMedia{
		SocialMediaURL: "http://test.com",
		Name:           "test name",
	}

	valid, err := govalidator.ValidateStruct(sosmed)
	assert.True(t, valid)
	assert.NoError(t, err)
}
