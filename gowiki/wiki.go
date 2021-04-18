package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

// This struct represents a page in a wiki-web-site
type Page struct { // this struct shows how a page is stored in memory
	Title string
	Body  []byte
}

// Should define the cwd for loading and saving data to spesific&different folders
var wd, _ = os.Getwd()

func (p *Page) save() error { //save method can be used for saing the page to the persisten memory
	filename := p.Title + ".txt"                // name the save file same with the page title
	path := wd + "/data/" + filename            // save the page to data folder
	return ioutil.WriteFile(path, p.Body, 0600) // give read permission to current user only
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	path := wd + "/data/" + filename   //load the page from data folder
	body, err := ioutil.ReadFile(path) // read the file which given as parameter and store the content to body variable
	if err != nil {
		return nil, err

	}

	return &Page{Title: title, Body: body}, nil // return a pointer to the Page readed from the file

}

// global variable for template caching
var templates = template.Must(template.ParseFiles(wd+"/tmpl/edit.html", wd+"/tmpl/view.html"))

// since we are repeating the same code in edit/view functions its better to create a new function
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p) // get the templates from tmpl folder
	if err != nil {
		fmt.Println("error here")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {

	p, err := loadPage(title) // html/template make that stuff error free
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound) // if requested page is not found then redirect to edit page to create it
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {

	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	// added a form field for editing the page
	//instead of hard-coding the html I've used an standalone html file
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

// redirect the user from root to FrontPage
func rootHandler(w http.ResponseWriter, r *http.Request, title string) {
	http.Redirect(w, r, "view/FrontPage", http.StatusFound)
	return

}

// to prevent path traversal bug create a regex
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

// wrapper function for these 3 handler function so they won't repeat each other
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		// Instead of editing the regex , i just put an exception here
		if m == nil {
			if r.URL.Path == "/" {
				fn(w, r, "")
				return
			}
			// no match
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func main() {

	//	p1 := &Page{Title: "Test1", Body: []byte("body for test1 page")}
	//	p1.save()
	//	p2, _ := loadPage("Test1")
	////	fmt.Println(string(p2.Body))
	http.HandleFunc("/", makeHandler(rootHandler))
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	log.Fatal(http.ListenAndServe(":8080", nil)) //start the server
}
