package spotify

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/karim-w/kyoto/utils"
)

type state struct {
	Token  TokenResponse `json:"token"`
	Expiry time.Time     `json:"expiry"`
	Cid    string        `json:"cid"`
	Sec    string        `json:"sec"`
}

func StartAuth() {
	authCodeChan := make(chan string, 10)
	reader := bufio.NewReader(os.Stdin)
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			query := r.URL.Query()
			code := query.Get("code")
			authCodeChan <- code
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		})
		http.ListenAndServe(":9499", nil)
	}()
	urlBuilder := strings.Builder{}
	_, err := urlBuilder.WriteString("https://accounts.spotify.com/authorize?")
	assert(err)
	_, err = urlBuilder.WriteString("response_type=code&client_id=")
	assert(err)
	fmt.Println("Spotify Client Id: (only saved on ur device)")
	clientId, err := reader.ReadString('\n')
	assert(err)
	if clientId == "" {
		panic("need The spotify client Id")
	}
	clientId = strings.TrimSpace(clientId)
	fmt.Println("Spotify Client secret: (only saved on ur device)")
	clientSecret, err := reader.ReadString('\n')
	assert(err)
	if clientSecret == "" {
		panic("need The spotify client secret")
	}
	clientSecret = strings.TrimSpace(clientSecret)
	_, err = urlBuilder.WriteString(clientId)
	assert(err)
	_, err = urlBuilder.WriteString(
		"&scope=user-read-currently-playing&redirect_uri=http://localhost:9499/callback&state=",
	)
	assert(err)
	_, err = urlBuilder.WriteString(utils.GenerateRandomString(16))
	assert(err)
	cmdExec := exec.Command("open", urlBuilder.String())
	go cmdExec.Run()
	code := <-authCodeChan
	token, expiry := InitOrDieClient(code, clientId, clientSecret).GetToken()
	s := state{
		*token,
		expiry,
		clientId,
		clientSecret,
	}
	byts, err := json.Marshal(s)
	assert(err)
	// convert to base64
	auth := base64.RawStdEncoding.EncodeToString(byts)
	// write to a .kyoto file
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	var file *os.File
	_, err = os.Stat(usr.HomeDir + "/.kyoto")

	if err != nil {
		err = nil
		file, err = create(usr.HomeDir + "/.kyoto")
		assert(err)
	} else {
		file, err = os.OpenFile(usr.HomeDir+"/.kyoto", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	}
	_, err = file.WriteString(auth)
	assert(err)
	os.Exit(0)
}

func SourceAuth() *Client {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	byts, err := os.ReadFile(usr.HomeDir + "/.kyoto")
	assert(err)
	byts, err = base64.RawStdEncoding.DecodeString(string(byts))
	assert(err)
	s := state{}
	err = json.Unmarshal(byts, &s)
	assert(err)
	cl := Client{
		token:        s.Token,
		expiresAt:    s.Expiry,
		clientId:     s.Cid,
		clientSecret: s.Sec,
	}
	return &cl
}

func create(p string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(p), 0770); err != nil {
		return nil, err
	}
	return os.Create(p)
}

func assert(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
