package main

import (
	"encoding/json"
	"fmt"
	"github.com/StillFantastic/bullshit/generator"
	"github.com/rs/cors"
	"golang.org/x/time/rate"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type Data struct {
	Topic  string
	MinLen int
}

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var visitors = make(map[string]*visitor)
var mu sync.Mutex

func init() {
	go cleanupVisitors()
}

func logRequest(topic string, minLen int) {
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

func getVisitor(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	v, exists := visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(1, 1000)
		visitors[ip] = &visitor{limiter, time.Now()}
		return limiter
	}

	v.lastSeen = time.Now()
	return v.limiter
}

func cleanupVisitors() {
	for {
		time.Sleep(time.Minute)

		mu.Lock()
		for ip, v := range visitors {
			if time.Since(v.lastSeen) > 20*time.Minute {
				delete(visitors, ip)
			}
		}
		mu.Unlock()
	}
}

func limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}

		limiter := getVisitor(ip)
		if limiter.Allow() == false {
			http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
func bullshitHandler(w http.ResponseWriter, r *http.Request) {
	var d Data
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
	}

	// logRequest(d.Topic, d.MinLen)
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
	err := http.ListenAndServe(addr, limit(handler))
	if err != nil {
		fmt.Println(err)
	}
}
