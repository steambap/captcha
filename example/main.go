package main

import (
	"fmt"
	"github.com/steambap/captcha"
	"html/template"
	"net/http"
)

func main() {
	http.HandleFunc("/", IndexHandle)
	http.HandleFunc("/captcha", CaptchaHandle)
	fmt.Println("Server start at port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func IndexHandle(w http.ResponseWriter, _ *http.Request) {
	doc, err := template.ParseFiles("index.html")
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	doc.Execute(w, nil)
}

func CaptchaHandle(w http.ResponseWriter, _ *http.Request) {
	img, err := captcha.New(150, 50)
	if err != nil {
		fmt.Fprint(w, nil)
		fmt.Println(err.Error())
		return
	}
	img.WriteTo(w)
}
