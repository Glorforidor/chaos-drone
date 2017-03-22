package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"

	"github.com/bieber/barcode"
)

func main() {
	// Load file data
	data, err := ioutil.ReadFile("demo.jpg")
	if err != nil {
		fmt.Println("Error 1")
		log.Fatal(err)
	}
	// Decode image
	m, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		fmt.Println("Error 2")
		log.Fatal(err)
	}
	img := barcode.NewImage(m)
	scanner := barcode.NewScanner().SetEnabledAll(true)

	symbols, _ := scanner.ScanImage(img)
	for _, s := range symbols {
		fmt.Println(s.Type.Name(), s.Data, s.Quality, s.Boundary)
	}
}
