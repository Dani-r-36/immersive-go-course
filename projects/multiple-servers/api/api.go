package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
)

func DatabaseConnection(database_url string, port int) {
	if database_url == ""{
		fmt.Fprintf(os.Stderr, "No database url given\n")
		os.Exit(1)
	}
	conn, err := pgx.Connect(context.Background(), database_url)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error connecting to database")
		os.Exit(1)
	}
	http.HandleFunc("/api/images.json",func(w http.ResponseWriter, r *http.Request){
		var jsonArray []byte
		queryParams := r.URL.Query()
		amount := queryParams.Get("indent")
		switch r.Method {
			case http.MethodGet:
				fmt.Println("get method from api")
				images, err := getImage("SELECT title, url, alt_text FROM public.images", conn)
				if err != nil {
					http.Error(w, "Error collecting data", http.StatusInternalServerError)
					fmt.Fprintln(os.Stderr, "Error collecting data")
					//erros handled in main.go
					os.Exit(1)
				}
				jsonArray = returnImage(amount, w, images)
				w.Header().Add("Content-Type", "text/json")
				w.Header().Add("Access-Control-Allow-Origin", "*")
				w.Header().Add("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
    			w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
				// Send it back!
				w.Write(jsonArray)
			case http.MethodPost:
				fmt.Println("post method")
				var image Image
				if err := json.NewDecoder(r.Body).Decode(&image); err != nil {
					http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
					fmt.Fprintln(os.Stderr, "Failed to decode json")
					return
				}
				returnedImage, err := addImage(image.Title, image.URL, image.Alt_Text, conn)
				if err != nil {
					//handle in main.go
					http.Error(w, "Failed to insert data", http.StatusInternalServerError)
					return
				}
				fmt.Fprintln(os.Stderr,"Sucessfully added image")
				jsonArray = returnImage(amount, w, returnedImage)
				w.Header().Add("Content-Type", "text/json")
				w.Header().Add("Access-Control-Allow-Origin", "*")
				w.Header().Add("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
    			w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
				// Send it back!
				w.Write(jsonArray)
			default:
				fmt.Fprintf(w, "Unsupported method: %s", r.Method)
			}
		})
	portString := fmt.Sprintf(":%d", port)
   	http.ListenAndServe(portString, nil)
	defer conn.Close(context.Background())
}