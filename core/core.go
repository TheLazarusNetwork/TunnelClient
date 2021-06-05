package core

import (
	"fmt"
	"os"
)

func CheckUser(apiKey string) string {
	fmt.Println(os.Getenv("USER_API_KEY"))
	if apiKey == os.Getenv("USER_API_KEY") {
		return "valid"
	} else {
		return "invalid"
	}
}
