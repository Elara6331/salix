package main

import (
	_ "embed"
	"log"
	"os"
	"time"

	"go.elara.ws/salix"
)

//go:embed readme.salix.html
var tmpl string

type User struct {
	Name           string
	LoggedIn       bool
	IsAdmin        bool
	RegisteredTime time.Time
}

var users = []User{
	{
		Name:           "Elara",
		LoggedIn:       true,
		IsAdmin:        true,
		RegisteredTime: time.Date(2023, time.January, 10, 10, 10, 10, 0, time.UTC),
	},
	{
		Name:           "Karl Marx",
		LoggedIn:       false,
		IsAdmin:        true,
		RegisteredTime: time.Date(1818, time.May, 5, 10, 10, 10, 0, time.UTC),
	},
	{
		Name:           "Jane Doe",
		LoggedIn:       true,
		IsAdmin:        false,
		RegisteredTime: time.Date(2050, time.March, 6, 10, 10, 10, 0, time.UTC),
	},
	{
		Name:           "John Doe <script>alert('hello!');</script>",
		LoggedIn:       false,
		IsAdmin:        false,
		RegisteredTime: time.Date(1700, time.November, 20, 10, 10, 10, 0, time.UTC),
	},
}

type Page struct {
	Title string
}

func main() {
	vars := map[string]any{
		"users": users,
		"page":  Page{Title: "Users"},
	}

	t, err := salix.New().
		WithVarMap(vars).
		WithEscapeHTML(true).
		ParseString("readme.salix.html", tmpl)
	if err != nil {
		log.Fatalln(err)
	}

	err = t.Execute(os.Stdout)
	if err != nil {
		log.Fatalln(err)
	}
}
