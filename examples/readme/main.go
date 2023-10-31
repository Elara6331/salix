/*
 * Salix - Go templating engine
 * Copyright (C) 2023 Elara Musayelyan
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

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

	t, err := salix.New().ParseString("readme.salix.html", tmpl)
	if err != nil {
		log.Fatalln(err)
	}

	err = t.WithVarMap(vars).
		WithEscapeHTML(true).
		Execute(os.Stdout)
	if err != nil {
		log.Fatalln(err)
	}
}
