package gadget

import (
	"errors"
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/merliot/dean"
	"github.com/merliot/device/prime"
	"golang.org/x/net/websocket"
)

var (
	id         = "gadget01"
	model      = "gadget"
	name       = "gadget-01"
	user       = "user"
	passwd     = "passwd"
	host       = "0.0.0.0"
	portdevice = "8050"
	portprime  = "8051"
)

func check(t *testing.T, err error) {
	if err != nil {
		t.Fatalf(err.Error())
	}
}

/*
func TestHomePage(t *testing.T) {
	url := fmt.Sprintf("http://%s:%s/device/%s/", host, portprime, id)
	println(url)
	req, err := http.NewRequest("GET", url, nil)
	check(t, err)
	req.SetBasicAuth(user, passwd)

	client := http.Client{}
	resp, err := client.Do(req)
	check(t, err)

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
*/

func newConfig(url, user, passwd string) (*websocket.Config, error) {
	origin := "http://localhost/"

	// Configure the websocket
	config, err := websocket.NewConfig(url, origin)
	if err != nil {
		return nil, err
	}

	if user != "" {
		// Set the basic auth header for the request
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.SetBasicAuth(user, passwd)
		config.Header = req.Header
	}

	return config, nil
}

func newDevice(port string) *dean.Server {
	gadget := New(id, model, name)
	server := dean.NewServer(gadget, user, passwd, port)
	go server.Run()
	// Wait a bit for server to spin up
	time.Sleep(time.Second)
	return server
}

func newPrime(port string) *prime.Prime {
	gadget := New(id, model, name)
	prime := prime.NewPrime("p1", "prime", "p1", port, user, passwd, gadget).(*prime.Prime)
	go prime.Serve()
	// Wait a bit for prime to spin up
	time.Sleep(time.Second)
	return prime
}

func newUrl(t *testing.T, port string, trunk bool) *url.URL {
	surl := "ws://" + host + ":" + port + "/ws/"
	if trunk {
		surl += "?trunk"
	}
	url, err := url.Parse(surl)
	check(t, err)
	return url
}

type client struct {
	Bottles int
	*testing.T
	*websocket.Conn
	in chan *dean.Packet
}

func (c *client) run() {
	var data []byte

	for {
		var packet = &dean.Packet{}

		err := websocket.Message.Receive(c.Conn, &data)
		if errors.Is(err, net.ErrClosed) {
			break
		}
		check(c.T, err)

		err = packet.SetMessage(data)
		check(c.T, err)

		c.in <- packet
	}
}

func (c *client) send(pkt *dean.Packet) error {
	data, err := pkt.Message()
	check(c.T, err)
	return websocket.Message.Send(c.Conn, string(data))
}

func newClient(t *testing.T, url *url.URL) *client {
	client := &client{
		T:  t,
		in: make(chan *dean.Packet),
	}
	cfg, err := newConfig(url.String(), user, passwd)
	check(t, err)
	client.Conn, err = websocket.DialConfig(cfg)
	check(t, err)
	go client.run()
	return client
}

// TestPrime
func TestPrime(t *testing.T) {

	var pkt = &dean.Packet{}

	prime := newPrime(portprime)
	device := newDevice(portdevice)

	url := newUrl(t, portdevice, true)
	c1 := newClient(t, url)
	c2 := newClient(t, url)

	url = newUrl(t, portprime, true)
	c3 := newClient(t, url)
	c4 := newClient(t, url)

	url = newUrl(t, portprime, false)
	device.Dial(url, 1)

	time.Sleep(time.Second)

	pkt.SetPath("get/state")
	check(t, c1.send(pkt))
	check(t, c2.send(pkt))
	pkt.PushTag("gadget01")
	check(t, c3.send(pkt))
	check(t, c4.send(pkt))

	pkt = <-c1.in
	pkt = <-c2.in
	pkt = <-c3.in
	pkt = <-c4.in

	c4.Close()
	c3.Close()
	c2.Close()
	c1.Close()

	device.Close()
	prime.Close()
}

// TestDevice serves a gadget device server to two ws clients.
func TestDevice(t *testing.T) {

	var pkt = &dean.Packet{}

	device := newDevice(portdevice)

	url := newUrl(t, portdevice, true)
	c1 := newClient(t, url)
	c2 := newClient(t, url)

	// Run a little state machine, step 1
	err := c1.send(pkt.SetPath("get/state"))
	check(t, err)

loop:
	for {
		select {
		case pkt = <-c1.in:
			switch pkt.Path {
			case "state":
				// step 2
				pkt.Unmarshal(c1)
				err = c2.send(pkt.ClearPayload().SetPath("get/state"))
				check(t, err)
			case "tookone":
				// step 4, we'll keep taking bottles until there are none
				c1.Bottles--
				if c1.Bottles == 0 && c2.Bottles == 0 {
					break loop
				}
				err = c1.send(pkt.ClearPayload().SetPath("takeone"))
				check(t, err)
			}
		case pkt = <-c2.in:
			switch pkt.Path {
			case "state":
				// step 3
				pkt.Unmarshal(c2)
				err = c1.send(pkt.ClearPayload().SetPath("takeone"))
				check(t, err)
			case "tookone":
				c2.Bottles--
				if c1.Bottles == 0 && c2.Bottles == 0 {
					break loop
				}
			}
		}
	}

	// Check if our accounting of how many bottles remain matches the
	// device's accounting

	check(t, c1.send(pkt.SetPath("get/state")))

	pkt = <-c1.in

	var c = &client{}
	pkt.Unmarshal(c)

	if c.Bottles != c1.Bottles {
		t.Fatalf("Expected %d bottles; got %d bottles", c1.Bottles, c.Bottles)
	}

	c2.Close()
	c1.Close()
	device.Close()
}

/*
func TestTagged(t *testing.T) {
	gadget := New(id, model, name)
	runner := dean.NewRunner(gadget, user, passwd)
	url, _ := url.Parse("ws://" + host + ":" + portprime + "/ws/?ping-period=4")
	ws := runner.Dial(url, 1)
	runner.Run()
	ws.Close()
}

func TestMain(m *testing.M) {

	// New gadget
	gadget := New(id, model, name).(*Gadget)

	// Start a new prime web server, adopting gadget
	prime := prime.NewPrime("p1", "prime", "p1", portprime, user, passwd, gadget).(*prime.Prime)

	// Run prime
	go prime.Serve()

	// Wait a bit for prime to spin up
	time.Sleep(time.Second)

	// Run the tests
	m.Run()

	// Cleanup
	prime.Close()
}
*/
