package config

import (
	"context"

	firebase "firebase.google.com/go/v4"
)

var (
	Firebase         *firebase.App
	FirebaseAppCheck *AppCheckClient
)

func SetupFirebase() error {
	var err error
	Firebase, err = firebase.NewApp(context.Background(), nil)

	if err != nil {
		return err
	}

	FirebaseAppCheck, err = NewAppCheck(context.Background())
	if err != nil {
		return err
	}

	return nil
}
