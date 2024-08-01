package helper

import (
	"fmt"
	"strconv"
)

func GenerateUniqueID(userID string, timestamp string) (string, error) {
	// Validate the timestamp
	_, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return "", fmt.Errorf("invalid timestamp format")
	}

	return fmt.Sprintf("%s-%s", userID, timestamp), nil
}
