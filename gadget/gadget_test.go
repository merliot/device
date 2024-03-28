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

func TestHomePage(t *testing.T) {
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

	got := strings.TrimSpace(string(body))
	want := id + " " + model + " " + name
	if got != want {
		t.Fatalf("Response '%s' not equal to what we wanted '%s'", got, want)
	}
}

func TestWebSocket(t *testing.T) {
	gadget := New(id, model, name)
	runner := dean.NewRunner(gadget, user, passwd)
	url := "ws://" + host + ":" + port + "/ws?ping-period=4"
	runner.Dial(url)
	runner.Run()
}

func TestMain(m *testing.M) {

	// Start the gadget as an http web server
	gadget := New(id, model, name).(*Gadget)
	server := dean.NewServer(gadget, user, passwd, port)
	go server.Run()

	// Wait a bit for http server to spin up
	time.Sleep(time.Second)

	// Run the tests
	m.Run()

	// Shut down the web server
	gadget.quit <- true
}
