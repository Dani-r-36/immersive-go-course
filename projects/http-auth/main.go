package main

import (
	// "html/template"
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"html"
	"net/http"
	"os"
	"golang.org/x/time/rate"
	"io"
	"strings"
)

var htmlStart string = `
<!DOCTYPE html>
<html>`
var greeting string = `<em>Hello, world</em>`
var queryHtml string = `
<p>Query parameters:
<ul>
<li>foo: `

    // tmpl, err := template.ParseFiles("template/index.html")
    // if err != nil {
    //   http.Error(w, err.Error(), http.StatusInternalServerError)
    //   return
    // }
	// tmpl.Execute(w, nil)



func rateLimit(limiter *rate.Limiter, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if limiter.Allow() {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		}
	})
}

func main() {
	http.HandleFunc("/",func(w http.ResponseWriter, r *http.Request){
		queryParams := r.URL.Query()
		foo := queryParams.Get("foo")
		switch r.Method {
			case http.MethodGet:
				w.Header().Set("Content-Type", "text/html")
				escaped := html.EscapeString(foo)
				response := htmlStart + greeting + queryHtml + escaped +"</li>\n</ul>"
				w.Write([]byte(response))
			case http.MethodPost:
				body := new(strings.Builder)
			if _, err := io.Copy(body, r.Body); err != nil {
				fmt.Printf("Error getting body: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal server error"))
				return
			}
			response := htmlStart + greeting +body.String()
			w.Write([]byte(html.EscapeString(response)))
			default:
				fmt.Fprintf(w, "Unsupported method: %s", r.Method)
			}
		})
		
	http.HandleFunc("/200", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("200"))
	})
	http.HandleFunc("/authenticated", func(w http.ResponseWriter, r *http.Request) {
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
			// expected username and password hashes.
			usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
			passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)
			if usernameMatch && passwordMatch {
				response := htmlStart + "Hello " + username
				w.Write([]byte(response))
				return
			}
		}

		// If the Authentication header is not present, is invalid, or the
		// username or password is wrong,
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		})

	limiter := rate.NewLimiter(100, 30)
	http.HandleFunc("/limited", rateLimit(limiter, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		w.Write([]byte("<!DOCTYPE html>\n<html>\nHello world!"))
	}))

	http.ListenAndServe(":8080", nil)
}