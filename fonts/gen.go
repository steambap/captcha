package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"path/filepath"
)

// This program generates a go file for Comismsh font

func main() {
	src, err := ioutil.ReadFile("Comismsh.ttf")
	if err != nil {
		log.Fatal(err)
	}

	buf := new(bytes.Buffer)
	fmt.Fprint(buf, "// DO NOT EDIT. This file is generated.\n\n")
	fmt.Fprint(buf, "package captcha\n\n")
	fmt.Fprint(buf, "// The following is Comismsh TrueType font data.\n")
	fmt.Fprint(buf, "var ttf = []byte{")
	for i, x := range src {
		if i&15 == 0 {
			buf.WriteByte('\n')
		}
		fmt.Fprintf(buf, "%#02x,", x)
	}
	fmt.Fprint(buf, "\n}\n")

	dst, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile(filepath.Join("../", "font.go"), dst, 0666); err != nil {
		log.Fatal(err)
	}
}
