package main

import (
	"bullshit/generator"
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/rs/cors"
	"os"
	"strconv"
	"time"
)

type Data struct {
	Topic  string
	MinLen int
}

func log(topic string, minLen int) {
	f, err := os.OpenFile("bullshit_log.txt", os.O_APPEND|os.OWRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	data := time.Now().Format("2006-01-02 15:04:05") + " Topic: " + topic + ", Length: " + strconv.Itoa(minLen)
	if _, err = f.WriteString(data); err != nil {
		fmt.Println(err)
	}
}

func bullshitHandler(w http.ResponseWriter, r *http.Request) {
	var d Data
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	log(d.Topic, d.MinLen)
	ret := generator.Generate(d.Topic, d.MinLen)
	w.Write([]byte(ret))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/bullshit", bullshitHandler)
	handler := cors.Default().Handler(mux)
	var addr string
	if len(os.Args) < 2 {
		addr = "0.0.0.0:10000"
	} else {
		addr = "0.0.0.0:" + os.Args[1]
	}
	err := http.ListenAndServe(addr, handler)
	if err != nil {
		fmt.Println(err)
	}
}
