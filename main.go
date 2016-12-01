package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"fmt"

	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

const header string =
`#pragma once

#include "../lib/ab_common.h"

`

const templateStart string = "const uint8_t %s[] PROGMEM = {\n"
const templateEnd string =
`};

const ab_Image %s PROGMEM = {%d, %d, %s};

`

func main() {
	var directory string
	flag.StringVar(&directory, "d", "", "Directory to load images from")
	flag.Parse()

	if directory == "" {
		flag.PrintDefaults()
		return
	}

	files, err := ioutil.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(header)
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		reader, err := os.Open(path.Join(directory, file.Name()))
		if err != nil {
			log.Fatal(err)
		}

		img, _, err := image.Decode(reader)
		if err != nil {
			log.Fatal(err)
		}

		imgName := strings.Replace(file.Name(), ".", "_", -1)


		fmt.Printf(templateStart, imgName)
		bounds := img.Bounds()
		for y := bounds.Min.Y; y < bounds.Max.Y; y += 8 {
			fmt.Printf("    ")
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				var value uint8
				var i uint8
				for i = 0; i < 8; i++ {
					yy := y + int(i)
					_, _, _, a := img.At(x, yy).RGBA()

					if yy >= bounds.Max.Y {
						a = 0
					}

					if a > 128 {
						value |= (1 << i)
					}
				}

				fmt.Printf("0x%0.2X, ", value)
			}
			fmt.Printf("\n")
		}

		width := (uint8(bounds.Max.X - bounds.Min.X + 7) / uint8(8)) * 8
		height := (uint8(bounds.Max.Y - bounds.Min.Y + 7) / uint8(8)) * 8
		fmt.Printf(templateEnd, "img_" + imgName, width, height, imgName)
	}
}
