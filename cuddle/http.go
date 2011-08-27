package cuddle

import (
	"appengine"
	"http"
	"appengine-go-backports/template"
)

const (
	defaultNameLen = 4
	clientIdLen    = 40
)

func init() {
	http.HandleFunc("/", root)
	http.HandleFunc("/post", post)
}

var rootTmpl = template.Must(template.ParseFile("tmpl/root.html"))

func root(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/favicon.ico" {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}

	c := appengine.NewContext(r)

	name := r.URL.Path[1:]
	if !validName.MatchString(name) {
		name = randName(defaultNameLen)
		http.Redirect(w, r, "/"+name, http.StatusFound)
		return
	}

	room, err := getRoom(c, name)
	if err != nil {
		http.Error(w, err.String(), http.StatusInternalServerError)
		return
	}

	// generate random identifier for the client
	token, err := room.AddClient(c, randName(clientIdLen))
	if err != nil {
		http.Error(w, err.String(), http.StatusInternalServerError)
		return
	}

	tmplData := struct { Room, Token string }{room.Name, token}
	err = rootTmpl.Execute(w, tmplData)
	if err != nil {
		http.Error(w, err.String(), http.StatusInternalServerError)
	}
}

func post(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	room, err := getRoom(c, r.FormValue("room"))
	if err != nil {
		http.Error(w, err.String(), http.StatusInternalServerError)
		return
	}

	err = room.Send(c, r.FormValue("msg"))
	if err != nil {
		http.Error(w, err.String(), http.StatusInternalServerError)
	}
}
