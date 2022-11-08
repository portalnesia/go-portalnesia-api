package webauthn

import (
	"testing"

	"portalnesia.com/api/app"
	"portalnesia.com/api/config"
	"portalnesia.com/api/models"
)

func TestBeginLogin(t *testing.T) {
	app.Initialization()

	db := config.DB
	var user models.User

	if err := db.First(&user, 2).Error; err != nil {
		t.Fatal(err)
	}

	_, _, err := BeginLogin(user)
	if err != nil {
		t.Fatal(err)
	}
}

func TestBeginRegister(t *testing.T) {
	app.Initialization()

	db := config.DB
	var user models.User

	if err := db.First(&user, 2).Error; err != nil {
		t.Error(err)
	}

	_, _, err := BeginRegister(user)
	if err != nil {
		t.Error(err)
	}
}
