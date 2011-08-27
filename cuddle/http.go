package cuddle

import (
	"appengine"
	"fmt"
	"http"
)

func init() {
	http.HandleFunc("/", root)
}

func root(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/favicon.ico" {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}

	c := appengine.NewContext(r)

	name := r.URL.Path[1:]
	if !validName.MatchString(name) {
		name = randName()
		http.Redirect(w, r, "/"+name, http.StatusFound)
		return
	}

	room, err := getRoom(c, name)
	if err != nil {
		http.Error(w, err.String(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Room %q", room.Name)
}
