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
	"embed"
	"io/fs"
	"log"
	"net/http"
	"time"

	"go.elara.ws/salix"
)

//go:embed tmpls
var tmpls embed.FS

func main() {
	tmplsFS, err := fs.Sub(tmpls, "tmpls")
	if err != nil {
		log.Fatalln(err)
	}

	ns := salix.New().WithVarMap(map[string]any{"now": time.Now})

	err = ns.ParseFSGlob(tmplsFS, "*.html")
	if err != nil {
		log.Fatalln(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, ok := ns.GetTemplate("home.html")
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		name := r.URL.Query().Get("name")
		vars := map[string]any{"title": "Home"}
		if name != "" {
			vars["name"] = name
		}

		err = tmpl.
			WithVarMap(vars).
			Execute(w)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		tmpl, ok := ns.GetTemplate("about.html")
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = tmpl.
			WithVarMap(map[string]any{"title": "About"}).
			Execute(w)
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
