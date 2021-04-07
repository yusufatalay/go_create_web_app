package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// This struct represents a page in a wiki-web-site
type Page struct { // this struct shows how a page is stored in memory
	Title string
	Body  []byte
}

func (p *Page) save() error { //save method can be used for saing the page to the persisten memory
	filename := p.Title + ".txt"                    // name the save file same with the page title
	return ioutil.WriteFile(filename, p.Body, 0600) // give read permission to current user only
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, error := ioutil.ReadFile(filename) // read the file which given as parameter and store the content to body variable
	if error != nil {
		log.Fatal(error)
		return nil, error

	} else {

		return &Page{Title: title, Body: body}, nil // return a pointer to the Page readed from the file
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]                         // assign  the part after view  to the title
	p, _ := loadPage(title)                                     //ignore the error for now
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body) // render the page in a reasonable format
}

func main() {

	//	p1 := &Page{Title: "Test1", Body: []byte("body for test1 page")}
	//	p1.save()
	//	p2, _ := loadPage("Test1")
	//	fmt.Println(string(p2.Body))

	http.HandleFunc("/view/", viewHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
