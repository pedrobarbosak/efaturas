package session

import (
	"context"

	"efaturas-xtreme/internal/api/middlewares"
	"efaturas-xtreme/pkg/errors"
)

func GetUserID(ctx context.Context) (string, error) {
	return get(ctx, middlewares.UserID)
}

func GetCredentials(ctx context.Context) (string, string, error) {
	uname, err := get(ctx, middlewares.Username)
	if err != nil {
		return "", "", err
	}

	pword, err := get(ctx, middlewares.Password)
	if err != nil {
		return "", "", err
	}

	return uname, pword, nil
}

func GetUser(ctx context.Context) (string, string, string, error) {
	userID, err := GetUserID(ctx)
	if err != nil {
		return "", "", "", err
	}

	uname, pword, err := GetCredentials(ctx)
	if err != nil {
		return "", "", "", err
	}

	return userID, uname, pword, nil
}

func get(ctx context.Context, key string) (string, error) {
	v, ok := ctx.Value(key).(string)
	if !ok || v == "" {
		return "", errors.NewUnauthorized("failed to get:", key)
	}

	return v, nil
}
