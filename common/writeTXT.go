package common

import (
	"log"
	"os"
)

func TXTWriter(content string) {
	file, err := os.OpenFile("test.txt", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0744)
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()
	file.WriteString(content)
}
