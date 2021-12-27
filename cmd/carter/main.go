package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	
	bt "github.com/pratikdeoghare/brashtag"
)

var (
	cardTmpl = template.Must(template.ParseFiles("./card-template2.html"))
)

type card struct {
	front string
	back  string
}

type deck []card

func (d deck) handler(w http.ResponseWriter, _ *http.Request) {
	type htmlCard struct {
		Front template.HTML
		Back  template.HTML
	}
	
	card := d[rand.Intn(len(d))]
	
	c := htmlCard{
		Front: template.HTML(card.front),
		Back:  template.HTML(card.back),
	}
	
	err := cardTmpl.Execute(w, c)
	if err != nil {
		log.Fatal(err)
	}
	
}

func main() {
	text, err := ioutil.ReadFile("./cards.bt")
	if err != nil {
		panic(err)
	}
	
	tree, err := bt.Parse(string(text))
	if err != nil {
		panic(err)
	}
	
	var d deck
	for _, kid := range tree.(bt.Bag).Kids() {
		switch x := kid.(type) {
		case bt.Bag:
			d = append(d, makeCard(x))
		}
	}
	
	http.HandleFunc("/", d.handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func makeCard(t bt.Bag) card {
	var front []string
	for _, kid := range t.Kids() {
		switch x := kid.(type) {
		case bt.Blob:
			front = append(front, x.Text())
		case bt.Code:
			front = append(front, code(x))
		case bt.Bag:
			t = x
		}
	}
	
	var back []string
	for _, kid := range t.Kids() {
		switch x := kid.(type) {
		case bt.Code:
			back = append(back, code(x))
		case bt.Blob:
			back = append(back, x.Text())
		}
	}
	
	return card{
		front: strings.Join(front, ""),
		back:  strings.Join(back, ""),
	}
}

func code(x bt.Code) string {
	
	// for math
	trimmed := strings.TrimSpace(x.Text())
	if trimmed != "" {
		if trimmed[0] == '$' {
			return x.Text()
		}
	}
	
	if len(x.Tag()) == 1 {
		return fmt.Sprintf("<code>%s</code>", x.Text())
	}
	return fmt.Sprintf("<pre>%s</pre>", x.Text())
}
