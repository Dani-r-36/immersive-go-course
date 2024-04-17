package main

import (
	// "html/template"
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"os"
)

var htmlStart string = `
<!DOCTYPE html>
<html>`
var greeting string = `<em>Hello, world</em>`
var queryHtml string = `
<p>Query parameters:
<ul>
<li>foo: `

func handleRequest(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	foo := queryParams.Get("foo")
	switch r.Method {
		case http.MethodGet:
    // tmpl, err := template.ParseFiles("template/index.html")
    // if err != nil {
    //   http.Error(w, err.Error(), http.StatusInternalServerError)
    //   return
    // }
	// tmpl.Execute(w, nil)
			w.Header().Set("Content-Type", "text/html")
			escaped := html.EscapeString(foo)
    		response := htmlStart + greeting + queryHtml + escaped +"</li>\n</ul>"
			w.Write([]byte(response))
  		case http.MethodPost:
    		body, _ := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			w.Header().Set("Content-Type", "text/plain")
			escaped := html.EscapeString(string(body))
			response := "<!DOCTYPE html><html>" + escaped + "</html>"
			w.Write([]byte(response))
  		default:
    		fmt.Fprintf(w, "Unsupported method: %s", r.Method)
  	}
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "This is the protected handler")
}

func authenticated(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if ok {
		// Calculate SHA-256 hashes for the provided and expected
		// usernames and passwords.
		usernameHash := sha256.Sum256([]byte(username))
		passwordHash := sha256.Sum256([]byte(password))
		expectedUsernameHash := sha256.Sum256([]byte(os.Getenv("AUTH_USERNAME")))
		expectedPasswordHash := sha256.Sum256([]byte(os.Getenv("AUTH_PASSWORD")))

		// Use the subtle.ConstantTimeCompare() function to check if 
		// the provided username and password hashes equal the  
		// expected username and password hashes. ConstantTimeCompare
		// will return 1 if the values are equal, or 0 otherwise. 
		// Importantly, we should to do the work to evaluate both the 
		// username and password before checking the return values to 
		// avoid leaking information.
		usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
		passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

		// If the username and password are correct, then call
		// the next handler in the chain. Make sure to return 
		// afterwards, so that none of the code below is run.
		if usernameMatch && passwordMatch {
			response := htmlStart + "Hello " + username
			w.Write([]byte(response))
			return
		}
	}

	// If the Authentication header is not present, is invalid, or the
	// username or password is wrong, then set a WWW-Authenticate 
	// header to inform the client that we expect them to use basic
	// authentication and send a 401 Unauthorized response.
	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}
func main() {
	http.HandleFunc("/", handleRequest)
	http.HandleFunc("/authenticated", authenticated)
	http.ListenAndServe(":8080", nil)
}