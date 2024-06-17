package exception

import (
	"fmt"
)

func Print(err error, message string) {
	if err != nil {
		fmt.Printf("%s : %v", message, err)
	}
}
