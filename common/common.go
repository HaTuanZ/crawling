package common

import (
	"fmt"
	"os"
)

func GetCurrentPath() string {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return pwd
}
