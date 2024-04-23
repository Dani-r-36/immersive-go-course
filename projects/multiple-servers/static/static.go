package static

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
)

type Config struct {
	Dir  string
	Port int
}

func Run(config Config) {
   log.Println("Hello!")
   http.HandleFunc("/",func(w http.ResponseWriter, r *http.Request){
	path := filepath.Join(config.Dir, r.URL.EscapedPath())
	http.ServeFile(w,r,path)
   })
   port := fmt.Sprintf(":%d", config.Port)
   http.ListenAndServe(port, nil)
}