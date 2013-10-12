package main

import (
	"bytes"
	"fmt"
	"github.com/hoisie/web"
	"html/template"
	"log"
)

func errorPage() string {
	return "We encountered an error. Sorry for the inconvenience."
}

func enumParams(cxt *web.Context) string {
	buf := bytes.NewBuffer(nil)
	for k, v := range cxt.Params {
		fmt.Fprintf(buf, "Got param %s == %s<br />", k, v)
	}
	return buf.String()
}

func printPage(name string) string {
	temp, err := template.ParseFiles("thanks.html")
	if err != nil {
		log.Println(err)
		return errorPage()
	}
	result := bytes.NewBuffer(nil)
	temp.ExecuteTemplate(result, "thanks.html", name)
	return result.String()
}

func main() {
	web.Get("/favicon.ico", func() string { return "" })
	web.Get("/params", enumParams)
	web.Get("/(.*)", printPage)
	web.Run("0.0.0.0:8080")
}
