package server

import (
	"log"
	"testing"
)

func TestPixivPicID(t *testing.T) {
	message := "htrbhspid0few456tgrhbr"
	pid, err := PixivPicID(message)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(pid)
}
