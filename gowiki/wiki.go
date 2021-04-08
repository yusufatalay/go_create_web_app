package main

import (
	"html/template"
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
	body, err := ioutil.ReadFile(filename) // read the file which given as parameter and store the content to body variable
	if err != nil {
		log.Fatal(err)
		return nil, err

	} else {

		return &Page{Title: title, Body: body}, nil // return a pointer to the Page readed from the file
	}
}

// since we are repeating the same code in edit/view functions its better to create a new function
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, err := template.ParseFiles(tmpl + ".html")
	// http.Error sends the specified http response code
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):] // assign  the part after view  to the title
	p, err := loadPage(title)           // html/template make that stuff error free

	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound) // if requested page is not found then redirect to edit page to create it
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):] // take the path after the edit
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	// added a form field for editing the page
	//instead of hard-coding the html I've used an standalone html file
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func main() {

	//	p1 := &Page{Title: "Test1", Body: []byte("body for test1 page")}
	//	p1.save()
	//	p2, _ := loadPage("Test1")
	////	fmt.Println(string(p2.Body))

	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil)) //start the server
}
