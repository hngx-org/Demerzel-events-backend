package firebase

import (
	"context"
	gFirebase "firebase.google.com/go/v4"
	"fmt"
	"google.golang.org/api/option"
)

type Firebase struct {
	app *gFirebase.App
}

var Inner Firebase

func Initialize() *Firebase {
	opt := option.WithCredentialsFile("service-account.json")
	app, err := gFirebase.NewApp(context.Background(), nil, opt)

	if err != nil {
		panic(err)
	}

	return &Firebase{
		app: app,
	}
}

func (f *Firebase) VerifyIDToken(ctx context.Context, idToken string) (map[string]interface{}, error) {
	auth, err := f.app.Auth(ctx)
	if err != nil {
		return nil, err
	}

	token, err := auth.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, err
	}

	fmt.Println(token)
	fmt.Println(token.Claims)

	return token.Claims, nil
}
