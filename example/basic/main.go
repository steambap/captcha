package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/steambap/captcha"
)

func main() {
	http.HandleFunc("/", indexHandle)
	http.HandleFunc("/captcha-default", captchaHandle)
	http.HandleFunc("/captcha-math", mathHandle)

	fmt.Println("Server start at port 8080")
	err := http.ListenAndServe(":8080", nil)
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
	img, err := captcha.New(150, 50)
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

func mathHandle(w http.ResponseWriter, _ *http.Request) {
	img, err := captcha.NewMathExpr(150, 50)
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
