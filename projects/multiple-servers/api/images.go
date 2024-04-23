package api

import (
	"fmt"
	"context"
	"os"
	"github.com/jackc/pgx/v5"
	"net/http"
	"html"
	"strconv"
	"encoding/json"
	"strings"
)

type Image struct {
	Title   string
	Alt_Text string
	URL     string
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
		images = append(images, Image{Title: title, URL: url, Alt_Text: altText})
	}
	return images, nil
}

func addImage(title, url, altText string, conn *pgx.Conn) ([]Image, error) {
	var image []Image
    query := "INSERT INTO public.images(title, url, alt_text) VALUES($1, $2, $3)  RETURNING id, title, url, alt_text"
    row := conn.QueryRow(context.Background(), query, title, url, altText)
	var id int
	var insertedTitle, insertedURL, insertedAltText string
	err := row.Scan(&id, &insertedTitle, &insertedURL, &insertedAltText)
	if err != nil {
        // Handle error
        fmt.Println("Error querying database:", err)
        return nil, err
    }
	image = append(image, Image{Title: insertedTitle, URL: insertedURL, Alt_Text: insertedAltText})
	fmt.Printf(insertedTitle, insertedURL, insertedAltText)
	return image, nil
}

func returnImage(amount string, w http.ResponseWriter, images []Image) []byte{
	var jsonArray []byte
	if amount != ""{
		indent := html.EscapeString(amount)
		w.Header().Set("Content-Type", "application/json")
		i, err := strconv.Atoi(indent)
		if err != nil {
			fmt.Println("Error:", err)
			return nil
		}
		blankSpace := strings.Repeat(" ", i) 
		jsonArray, _ = json.MarshalIndent(images," ",blankSpace)
	} else{
		jsonArray, _ = json.Marshal(images)
	}
	return jsonArray
}

