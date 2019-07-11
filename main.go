package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/BradleyJohnson/smpltrace"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func main() {
	var addr = flag.String("addr", ":8080", "The addr of the application.")
	flag.Parse()
	r := newRoom()
	r.tracer = smpltrace.New(os.Stdout)

	// call the Handle function which registers the Handler with the pattern in defaultservermux
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))

	http.Handle("/login", &templateHandler{filename: "login.html"})

	// our custom type room can be passed into the http.Handle function
	// which requires a valid Handler interface. room is a valid Handler
	// type simply because it implements the ServeHTTP function.

	// Interestingly though, we can't visit that URL directly and actually the jacascript
	// in chat.html is responsible for hitting that route and creating the socket for the client
	http.Handle("/room", r)

	go r.run()
	// start the server
	log.Println("Starting webserver on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe", err)
	}
}

// templateHandler, like room, satifies the Handler interface
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates",
			t.filename)))
	})
	t.templ.Execute(w, r)
}
