package gadget

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/merliot/dean"
)

var (
	id     = "id"
	model  = "gadget"
	name   = "name"
	user   = "user"
	passwd = "passwd"
	host   = "0.0.0.0"
	port   = "8050"
)

func testHomePage(t *testing.T) {
	url := fmt.Sprintf("http://%s:%s", host, port)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth(user, passwd)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		t.Errorf("Get %s failed: %s", url, err)
	}

	contents := strings.TrimSpace(string(body))
	println(contents)
}

func TestBasic(t *testing.T) {
	gadget := New(id, model, name).(*Gadget)
	server := dean.NewServer(gadget, user, passwd, port)
	go server.Run()

	// Wait a bit for http server to spin up
	time.Sleep(time.Second)

	t.Run("TestHomePage", testHomePage)

	gadget.quit <- true
}
