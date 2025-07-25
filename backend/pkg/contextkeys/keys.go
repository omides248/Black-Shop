package contextkeys

import (
	"context"
	"fmt"
)

type contextKey string

const UserIDKey = contextKey("userID")

func GetUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok || userID == "" {
		return "", fmt.Errorf("user ID not found in context")
	}
	return userID, nil
}
