package cuddle

import (
	"appengine"
	"http"
	"appengine-go-backports/template"
)

const (
	defaultNameLen = 4 // the length of randomly generated room names
	clientIdLen    = 40
)

func init() {
	// Register our handlers with the http package.
	http.HandleFunc("/", root)
	http.HandleFunc("/post", post)
}

// rootTmpl is the main (and only) HTML template.
var rootTmpl = template.Must(template.ParseFile("tmpl/root.html"))

// root is an http handler that joins or creates a Room,
// creates a new Client, and writes the HTML response.
func root(w http.ResponseWriter, r *http.Request) {
	// Get the name from the request URL.
	name := r.URL.Path[1:]

	// If no valid name is provided,
	// generate a new one and redirect there.
	if !ValidName.MatchString(name) {
		name = RandName(defaultNameLen)
		http.Redirect(w, r, "/"+name, http.StatusFound)
		return
	}

	c := appengine.NewContext(r)

	// Get or create the Room.
	room, err := getRoom(c, name)
	if err != nil {
		http.Error(w, err.String(), 500)
		return
	}

	// Create a new Client, getting the channel token.
	token, err := room.AddClient(c, RandName(clientIdLen))
	if err != nil {
		http.Error(w, err.String(), 500)
		return
	}

	// Render the HTML template.
	data := struct{ Room, Token string }{room.Name, token}
	err = rootTmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.String(), 500)
	}
}

// post broadcasts a message to a specified Room.
func post(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	// Get the room.
	room, err := getRoom(c, r.FormValue("room"))
	if err != nil {
		http.Error(w, err.String(), 500)
		return
	}

	// Send the message to the clients in the room.
	err = room.Send(c, r.FormValue("msg"))
	if err != nil {
		http.Error(w, err.String(), 500)
	}
}
