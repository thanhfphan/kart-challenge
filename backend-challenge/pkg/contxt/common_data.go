package contxt

import (
	"context"
	"fmt"
)

const KeyUID = "uid"

func GetUserID(ctx context.Context) (int64, error) {
	wrapperCtx, err := GetAppWrapper(ctx)
	if err != nil {
		return 0, fmt.Errorf("get wrapper context err=%w", err)
	}

	return wrapperCtx.GetInt64(KeyUID)
}

func GetAuthToken(ctx context.Context) (string, error) {
	wrapperCtx, err := GetAppWrapper(ctx)
	if err != nil {
		return "", fmt.Errorf("get wrapper context err=%w", err)
	}

	return wrapperCtx.GetString("token")
}

func GetAcceptLanguage(ctx context.Context) (string, error) {
	wrapperCtx, err := GetAppWrapper(ctx)
	if err != nil {
		return "", fmt.Errorf("get wrapper context err=%w", err)
	}

	return wrapperCtx.GetString("accept-language")
}

func GetClientIP(ctx context.Context) (string, error) {
	wrapperCtx, err := GetAppWrapper(ctx)
	if err != nil {
		return "", fmt.Errorf("get wrapper context err=%w", err)
	}

	return wrapperCtx.GetString("ip-address")
}
