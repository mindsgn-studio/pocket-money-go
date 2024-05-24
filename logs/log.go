package logs

import (
	"fmt"
	"time"
)

func LogError(message string) {
	fmt.Println(time.Now().UTC(), message)
}
