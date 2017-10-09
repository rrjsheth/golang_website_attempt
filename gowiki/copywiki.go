package main

import (
	"fmt"
	"regexp"
	"io/ioutil"
	"log"
	"net/http"
	"html/template"
)

type Page struct {
	Title string
	Body []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page,error) {
	filename := title+ ".txt"
	body,err := ioutil.ReadFile(filename)
	if err != nil{
		return nil,err
	}
	return &Page{Title: title, Body: body},nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p,err := loadPage(title)
	
	if err!=nil{
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, title, p) 
}

func editHandler( w http.ResponseWriter, r *http.Request, title string){
	//check if password is true or not
	passwordCorrect := true
	if username:= r.FormValue("username"); username != "ravisheth"{
		fmt.Println(username)
		passwordCorrect = false
	}
	if password:= r.FormValue("password"); password != "ravisheth"{
		passwordCorrect = false
	}

	//if password is true then let them edit otherwise redirect to login page
	if passwordCorrect {
		p, err:= loadPage(title)
		if err != nil {
			p = &Page{Title: title}
		}
		renderTemplate(w, "edit", p)
	} else{
		http.Redirect(w,r,"/login/"+title, http.StatusFound)
	}
}

//start of changes

func loginHandler( w http.ResponseWriter, r *http.Request, title string){
	p,_:=loadPage(title)
	renderTemplate(w, "login",p)
}

//end of changes

var templates = template.Must(template.ParseFiles("edit.html","login.html","aboutMe.html",
				"CoverLetter.html","gitHubPage.html", "resumePage.html", "linkedInPage.html"))

func renderTemplate( w http.ResponseWriter, tmpl string, p *Page){
	err := templates.ExecuteTemplate(w,tmpl+ ".html", p)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func saveHandler( w http.ResponseWriter, r *http.Request, title string){
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err!= nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func docsHandler( w http.ResponseWriter, r *http.Request, title string){ 
	fmt.Println(title)
	fmt.Println(r.URL.Path)
	path := r.URL.Path+ ".doc"
	fmt.Println(path)
	w.Header().Set("Content-Disposition", "attachment; filename=resume.doc")
w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	http.ServeFile(w, r, path)
}	
var validPath = regexp.MustCompile("^/(edit|save|view|login|docs)/([a-zA-Z0-9]+)$")


func makeHandler( fn func( http.ResponseWriter, *http.Request, string)) http.HandlerFunc{
	return func( w http.ResponseWriter, r *http.Request ){
		//extract page title from request
		//and call the provided handler 'fn'
		fmt.Println(r.URL.Path)
		m := validPath.FindStringSubmatch(r.URL.Path)
		fmt.Println(m)
		if m == nil{
			http.NotFound(w,r)
			return
		}
		fn(w,r,m[2])
	}

}
func main(){
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/login/", makeHandler(loginHandler))
	http.HandleFunc("/docs/", makeHandler(docsHandler))
	log.Fatal(http.ListenAndServe(":8080",nil))
}
