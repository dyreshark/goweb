package main

import (
	"bytes"
	"crypto/tls"
	"errors"
	"github.com/hoisie/web"
	"html/template"
	"io/ioutil"
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

func tryAddMessage(ctx *web.Context) {
	msg, ok := ctx.Params["message"]
	if !ok {
		return
	}

	user, err := ctx.Request.Cookie("user")
	if err != nil {
		return
	}

	err = saveMessage(user.Value, msg)
	if err != nil {
		log.Println("Failed to save message", err)
	}
}

func homePage(ctx *web.Context) string {
	username, err := ctx.Request.Cookie("user")

	if ctx.Request.Method == "POST" {
		tryAddMessage(ctx)
	}

	// No username cookie? Cool. Redirect to login.
	if err != nil || username.Value == "" {
		ctx.SetCookie(web.NewCookie("user", "", -1))
		ctx.Redirect(302, "/login")
		return "Redirecting to /login..."
	}

	result := bytes.NewBuffer(nil)

	params := map[string]interface{}{
		"Name":  username.Value,
		"Posts": getMessages(),
	}

	err = pages.ExecuteTemplate(result, "home.html", params)
	if err != nil {
		log.Println("Error rendering /home:", err)
	}

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
	web.Post("/home", homePage)
	web.Get("/home", homePage)

	// HTTPS
	pkey, err := ioutil.ReadFile("privkey.key")
	if err != nil {
		log.Fatal("Could not read privkey", err)
	}

	cert, err := ioutil.ReadFile("cert.cert")
	if err != nil {
		log.Fatal("Could not read cert", err)
	}

	keypair, err := tls.X509KeyPair(cert, pkey)
	if err != nil {
		log.Fatal("Unable to create key pair", err)
	}

	cfg := tls.Config{
		Time:         nil,
		Certificates: []tls.Certificate{keypair},
	}

	web.RunTLS("0.0.0.0:8081", &cfg)
}
