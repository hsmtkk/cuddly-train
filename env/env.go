package env

import (
	"fmt"
	"os"
	"strconv"
)

func RequiredVar(name string) (string, error) {
	val := os.Getenv(name)
	if val == "" {
		return "", fmt.Errorf("env var %s is not defined", name)
	}
	return val, nil
}

func Port() (int, error) {
	portStr, err := RequiredVar("PORT")
	if err != nil {
		return 0, err
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 0, fmt.Errorf("strconv.Atoi failed; %w", err)
	}
	return port, nil
}
