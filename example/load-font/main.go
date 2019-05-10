package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/steambap/captcha"
	"golang.org/x/image/font/gofont/goregular"
	"github.com/golang/freetype/truetype"
)

var ttf *truetype.Font

func main() {
	var err error

	ttf, err = captcha.LoadFont(goregular.TTF)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", indexHandle)
	http.HandleFunc("/captcha", captchaHandle)

	fmt.Println("Server start at port 8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func indexHandle(w http.ResponseWriter, _ *http.Request) {
	doc, err := template.ParseFiles("index.html")
	if err != nil {
		fmt.Fprint(w, err.Error())

		return
	}

	err = doc.Execute(w, nil)
	if err != nil {
		fmt.Println(err.Error())

		return
	}
}

func captchaHandle(w http.ResponseWriter, _ *http.Request) {
	img, err := captcha.New(150, 50, func(options *captcha.Options) {
		options.FontScale = 0.8
		options.TTFFont = ttf
	})
	if err != nil {
		fmt.Fprint(w, nil)
		fmt.Println(err.Error())

		return
	}

	err = img.WriteImage(w)
	if err != nil {
		fmt.Println(err.Error())

		return
	}
}
