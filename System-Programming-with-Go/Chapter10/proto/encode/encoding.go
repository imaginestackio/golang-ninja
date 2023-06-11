package main

import (
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/imaginedevops/Hands-On-Systems-Programming-with-Go/Chapter10/proto/gen"
)

func main() {
	var char = gen.Character{
		Name:        "George",
		Surname:     "Gammell Angell",
		YearOfBirth: 1834,
		Job:         "professor emeritus",
	}
	b := proto.NewBuffer(nil)
	if err := b.EncodeMessage(&char); err != nil {
		log.Fatalln(err)
	}
	log.Printf("%q", b.Bytes())
}