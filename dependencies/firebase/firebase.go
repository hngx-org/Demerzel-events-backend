package firebase

import (
	"context"
	gFirebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
	"os"
)

type Firebase struct {
	app *gFirebase.App
}

var Inner Firebase

func Initialize() {
	opt := option.WithCredentialsJSON([]byte(os.Getenv("FIREBASE_SERVICE_ACCOUNT")))
	app, err := gFirebase.NewApp(context.Background(), nil, opt)

	if err != nil {
		panic(err)
	}

	Inner = Firebase{
		app: app,
	}
}

func (f *Firebase) VerifyIDToken(idToken string) (map[string]interface{}, error) {
	ctx := context.Background()
	auth, err := f.app.Auth(ctx)
	if err != nil {
		return nil, err
	}

	token, err := auth.VerifyIDToken(ctx, idToken)

	if err != nil {
		return nil, err
	}

	return token.Claims, nil
}
