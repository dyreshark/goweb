package main

import (
	"bytes"
	"errors"
	"github.com/hoisie/web"
	"html/template"
	"log"
)

var pages *template.Template

var users = map[string]string{
	"root":   "toor",
	"george": "yolo",
}

func tryLogin(ctx *web.Context) (username string, err error) {
	username, ok := ctx.Params["user"]
	if !ok {
		return "", errors.New("Need to give a username")
	}

	password, ok := ctx.Params["pass"]
	if !ok {
		return "", errors.New("Need to give a password")
	}

	userPassword, ok := users[username]
	if !ok || userPassword != password {
		return "", errors.New("Incorrect password.")
	}
	return username, nil
}

func loginPage(ctx *web.Context) string {
	showLoginError := false
	if ctx.Request.Method == "POST" {
		name, err := tryLogin(ctx)
		if err != nil {
			showLoginError = true
		} else {
			// So secure, I know. Look into SetSecureCookie if you want
			// logins that actually work.
			ctx.SetCookie(web.NewCookie("user", name, 2400))
			ctx.Redirect(302, "/home")
			return "Redirecting you to /home..."
		}
	}

	result := bytes.NewBuffer(nil)
	pages.ExecuteTemplate(result, "login.html", showLoginError)
	return result.String()
}

func printPage(ctx *web.Context) string {
	username, err := ctx.Request.Cookie("user")

	// No username cookie? Cool. Redirect to login.
	if err != nil || username.Value == "" {
		ctx.SetCookie(web.NewCookie("user", "", -1))
		ctx.Redirect(302, "/login")
		return "Redirecting to /login..."
	}

	result := bytes.NewBuffer(nil)
	pages.ExecuteTemplate(result, "home.html", username.Value)
	return result.String()
}

func main() {
	var err error
	pages, err = template.ParseFiles("login.html", "home.html")
	if err != nil {
		log.Fatal("Could not parse template files properly:", err)
	}
	web.Get("/favicon.ico", func() string { return "" })
	web.Get("/login", loginPage)
	web.Post("/login", loginPage)
	web.Get("/home", printPage)
	web.Run("0.0.0.0:8080")
}
