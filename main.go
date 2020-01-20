package main

import (
	"bullshit/generator"
	"encoding/json"
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/cors"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Data struct {
	Topic  string
	MinLen int
}

func log(topic string, minLen int) {
	f, err := os.OpenFile("bullshit_log.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	data := time.Now().Format("2006-01-02 15:04:05") + " Topic: " + topic + ", Length: " + strconv.Itoa(minLen) + "\n"
	if _, err = f.WriteString(data); err != nil {
		fmt.Println(err)
	}
}


func ipIsBanned(ip string) bool {
	bannedIp := os.Getenv("BANNED_IP")
	ipList := strings.Split(bannedIp, ",")
	for _, v := range ipList {
		if strings.Count(ip, v) > 0 {
			return true
		}
	}
	return false
}

func bullshitHandler(w http.ResponseWriter, r *http.Request) {
	var d Data
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if ipIsBanned(r.RemoteAddr) {
		return
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
