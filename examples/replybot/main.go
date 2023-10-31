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

	"go.elara.ws/salix"
)

//go:embed replybot.salix.txt
var tmpl string

var matches = [][]string{
	{"https://i.imgur.com/aaaaaaa.jpg", "aaaaaaa.jpg"},
	{"https://imgur.com/a/bbbbbbb.jpg", "a/bbbbbbb.jpg"},
}

func main() {
	vars := map[string]any{
		"type":    "post",
		"matches": matches,
	}

	t, err := salix.New().
		ParseString("replybot.salix.txt", tmpl)
	if err != nil {
		log.Fatalln(err)
	}

	err = t.WithVarMap(vars).Execute(os.Stdout)
	if err != nil {
		log.Fatalln(err)
	}
}
