package main

import (
	"fmt"
	"gee"
	"net/http"
)

//part3: gee初步框架
func index (w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w,"URL.Paht=%q\n", r.URL.Path)
}

func helloFunc (w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w,"hello URL.Paht=%q\n", r.URL.Path)
}

func main () {
	r := gee.New()
	r.GET("/", index)
	r.GET("/hello", helloFunc)
	r.Run(":9999")
}