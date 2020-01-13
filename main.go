package main

import (
	"bullshit/generator"
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/rs/cors"
)

type Data struct {
	Topic  string
	MinLen int
}

func bullshitHandler(w http.ResponseWriter, r *http.Request) {
	var d Data
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	ret := generator.Generate(d.Topic, d.MinLen)
	w.Write([]byte(ret))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/bullshit", bullshitHandler)
	handler := cors.Default().Handler(mux)
	err := http.ListenAndServe("0.0.0.0:80", handler)
	if err != nil {
		fmt.Println(err)
	}
}
