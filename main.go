package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/karim-w/kyoto/spotify"
	"github.com/karim-w/kyoto/tui"
	"github.com/karim-w/kyoto/utils"
)

var authCodeChan = make(chan string, 10)

func main() {
	authCodeChan = make(chan string, 10)
	go startServer()
	urlBuilder := strings.Builder{}
	_, err := urlBuilder.WriteString("https://accounts.spotify.com/authorize?")
	assert(err)
	_, err = urlBuilder.WriteString("response_type=code&client_id=")
	assert(err)
	clientId := os.Getenv("SPOTIFY_CLIENT_ID")
	if clientId == "" {
		panic("SPOTIFY_CLIENT_ID Missing")
	}
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	if clientSecret == "" {
		panic("SPOTIFY_CLIENT_SECRET Missing")
	}

	_, err = urlBuilder.WriteString(clientId)
	assert(err)
	_, err = urlBuilder.WriteString(
		"&scope=user-read-currently-playing&redirect_uri=http://localhost:8080/callback&state=",
	)
	assert(err)
	_, err = urlBuilder.WriteString(utils.GenerateRandomString(16))
	assert(err)
	cmdExec := exec.Command("open", urlBuilder.String())
	go cmdExec.Run()
	code := <-authCodeChan
	client := spotify.InitOrDieClient(code, clientId, clientSecret)
	tui.Start(client)
}

func startServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		code := query.Get("code")
		authCodeChan <- code
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	http.ListenAndServe(":8080", nil)
}

func assert(err error) {
	if err != nil {
		log.Println(err.Error())
		panic(err)
	}
}
