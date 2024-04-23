package main

import (
	// "html/template"
	"context"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"os"
	"strconv"
	"github.com/jackc/pgx/v5"
	// "golang.org/x/time/rate"
	// "io"
	"strings"
)
type Image struct {
	Title   string
	AltText string
	URL     string
}

//this overrides default marshal beahviour converting keys to snakecase
func (i Image) MarshalJSON() ([]byte, error) {
	type Alias Image
	return json.Marshal(&struct {
		Title    string `json:"title"`
		Alt_Text string `json:"alt_text"`
		URL      string `json:"url"`
		*Alias
	}{
		Title:    i.Title,
		Alt_Text: i.AltText,
		URL:      i.URL,
		Alias:    (*Alias)(&i),
	})
}

func databaseConnection(database_url string) (*pgx.Conn, error){
	conn, err := pgx.Connect(context.Background(), database_url)
	if err != nil {
		return nil, err
	}
	return conn, nil
	}

func getImage(query string, conn *pgx.Conn) ([]Image, error){
	var images []Image
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
        // Handle error
        fmt.Println("Error querying database:", err)
        return nil, err
    }
    defer rows.Close()
    for rows.Next(){
		var title, url, altText string
		err = rows.Scan(&title, &url, &altText)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Couldn't scan the row")
		}
		images = append(images, Image{Title: title, URL: url, AltText: altText})
	}
	return images, nil
}

func addImage(title, url, altText string, conn *pgx.Conn) error {
    query := "INSERT INTO public.images(title, url, alt_text) VALUES($1, $2, $3)"
    _, err := conn.Exec(context.Background(), query, title, url, altText)
	if err != nil {
        // Handle error
        fmt.Println("Error querying database:", err)
        return err
    }
	return nil
}


func main() {

	database_url := os.Getenv("DATABASE_URL")
	if database_url == ""{
		fmt.Fprintf(os.Stderr, "not connected to database\n")
		os.Exit(1)
	}
	conn, err := databaseConnection(database_url)
	if err != nil {
			fmt.Fprintln(os.Stderr, "Error connecting to database")
	}

	http.HandleFunc("/images.json",func(w http.ResponseWriter, r *http.Request){
		var jsonArray []byte
		queryParams := r.URL.Query()
		amount := queryParams.Get("indent")
		switch r.Method {
			case http.MethodGet:
				fmt.Println("get method")
				images, err := getImage("SELECT title, url, alt_text FROM public.images", conn)
				if err != nil {
					fmt.Fprintln(os.Stderr, "Error collecting data")
					os.Exit(1)
				}
				if amount != ""{
					indent := html.EscapeString(amount)
					w.Header().Set("Content-Type", "application/json")
					i, err := strconv.Atoi(indent)
					if err != nil {
						fmt.Println("Error:", err)
						return
					}
					blankSpace := strings.Repeat(" ", i) 
					jsonArray, _ = json.MarshalIndent(images," ",blankSpace)
				} else{
					jsonArray, _ = json.Marshal(images)
				}
				w.Write(jsonArray)
			case http.MethodPost:
				fmt.Println("post method")
				var image Image
				if err := json.NewDecoder(r.Body).Decode(&image); err != nil {
					http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
					fmt.Fprintln(os.Stderr, "Failed to decode json")
					return
				}

				if err := addImage(image.Title, image.URL, image.AltText, conn); err != nil {
					http.Error(w, "Failed to insert data", http.StatusBadRequest)
					fmt.Fprintln(os.Stderr, "Failed to insert image")
					return
				}
				fmt.Fprintln(os.Stderr,"Sucessfully added image")
			default:
				fmt.Fprintf(w, "Unsupported method: %s", r.Method)
			}
		})
	
	http.ListenAndServe(":8080", nil)
	defer conn.Close(context.Background())
}