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
	"net/http"

	"go.elara.ws/salix"
)

//go:embed post.html
var tmpl string

type Post struct {
	ID       int
	Message  string
	Children []*Post
}

var posts = Post{
	ID:      0,
	Message: "Hello",
	Children: []*Post{
		{
			ID:      12,
			Message: "Hi!",
			Children: []*Post{
				{
					ID:      13234,
					Message: "Test 2",
				},
			},
		},
		{
			ID:      13,
			Message: "Test",
		},
	},
}

func main() {
	t, err := salix.New().ParseString("post.html", tmpl)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err = t.WithVarMap(map[string]any{"post": posts}).Execute(w)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	log.Println("Starting HTTP server on port 8080")

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalln(err)
	}
}
