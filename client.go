/*
	client instantiation
	beat setup in a goroutine
	beats sent to a remote server every timeout
	failures can be tolerated to allow for shaky connections
	Close is called on multiple failures over the threshold
	should be threadsafe
*/

package heartbeat

import (
	"log"
	"net"
	"sync"
	"time"
)

type Client struct {
	name    string
	conn    net.Conn
	ch      chan string
	mu      sync.Mutex
	timeout time.Duration
}

// change the timeout of the function
func (c *Client) ChangeTimeout(t time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.timeout = t
}

// send a heartbeat to the conn
func (c *Client) Beat() {
	defer c.Close()
	out := time.Tick(c.timeout)
	for {
		select {
		case <-out:
			c.conn.Write([]byte("ok"))
		}
	}
}

// close down the channels and also the net.conn
func (c *Client) Close() {
	close(c.ch)
	c.conn.Close()
	log.Printf("Connection to %s closed", c.name)
}

func CreateClient(name string, addr string, time time.Duration) (c *Client) {
	c.name = name
	// do a dns lookup for the addr and grab connection net.Dial("tcp",addr)
	conn, _ := net.Dial("tcp", addr)
	c.conn = conn
	c.ch = make(chan string)
	c.timeout = time
	return
}
