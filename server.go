package heartbeat

import (
	"time"
	"log"
	"sync"
	"net"
)

type Service struct {
	name	string
	conn	net.Conn
	mu	sync.Mutex
	timeout	time.Duration
	fails	int
	threshold	int
}

func (s *Service) ResetFailures() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.fails = 0
}

func (s *Service) ChangeThreshold(thres time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.threshold = thres
}

func (s *Service) Receive() {
	defer s.Close()
	timeout := time.After(s.timeout)
	for {
		select {
			// performing fine
			case <-timeout:
				s.fails = s.fails + 1
				if s.fails > s.threshold {
					return
				}
		}
	}
}

func (s *Service) Close() {
	s.conn.Close()
}

func CreateServer(name string, addr net.Addr, t time.Duration, limit int) (s *Service) {
	s.name = name
	// resolve address to connection
	//conn, _ := net.Dial("tcp",addr)
	c, _ := net.Listen("tcp",addr)
	conn, _ := c.Accept()
	s.conn = conn
	s.timeout = t
	s.threshold = limit
}
