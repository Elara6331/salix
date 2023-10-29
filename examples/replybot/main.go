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
		WithVarMap(vars).
		ParseString("replybot.salix.txt", tmpl)
	if err != nil {
		log.Fatalln(err)
	}

	err = t.Execute(os.Stdout)
	if err != nil {
		log.Fatalln(err)
	}
}
