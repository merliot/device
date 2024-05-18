package gadget

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/merliot/dean"
	"github.com/merliot/device/prime"
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
	url := fmt.Sprintf("http://%s:%s/device/%s/", host, port, id)
	println(url)
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
	url, _ := url.Parse("ws://" + host + ":" + port + "/ws/?ping-period=4")
	ws := runner.Dial(url, 1)
	runner.Run()
	ws.Close()
}

func TestMain(m *testing.M) {

	// New gadget
	gadget := New(id, model, name).(*Gadget)

	// Start a new prime web server, adopting gadget
	prime := prime.NewPrime("p1", "prime", "p1", port, user, passwd, gadget).(*prime.Prime)

	// Run prime
	go prime.Serve()

	// Wait a bit for prime to spin up
	time.Sleep(time.Second)

	// Run the tests
	m.Run()

	// Cleanup
	prime.Close()
}
