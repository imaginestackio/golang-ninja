package main

import (
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/imaginedevops/Hands-On-Systems-Programming-with-Go/Chapter10/proto/gen"
)

func main() {
	b := proto.NewBuffer([]byte(
		"/\n\x06George\x12\x0eGammell Angell" +
			"\x1a\x12professor emeritus \xaa\x0e",
	))
	var char gen.Character
	if err := b.DecodeMessage(&char); err != nil {
		log.Fatalln(err)
	}
	log.Printf("%+v", char)
}
