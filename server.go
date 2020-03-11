package scale

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/devplayg/hippo"
	"github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

var (
	log *logrus.Logger
)

const (
	NewLine  byte = 0x0a
	BuffSize      = 32
)

type Server struct {
	hippo.Launcher // DO NOT REMOVE; Launcher links servers and engine each other.
	sampleData     []byte
	sampleDataLen  int
	idx            int
	initialized    bool
}

func NewServer(sampleData string) *Server {
	b := []byte(sampleData)
	return &Server{
		sampleData:    b,
		sampleDataLen: len(b),
		idx:           rand.Intn(len(b)),
	}
}

//
//func initData(b []byte) {
//	b = nil
//}
//
//func initString(b []string) {
//	b = nil
//}
//
//func initMap(m map[string]bool) {
//	m = map[string]bool{
//		"a":  true,
//	}
//
//}

func (s *Server) Start() error {
	log = s.Log

	log.Debugf("sample data length=%d", len(s.sampleData))

	buff := make([]byte, BuffSize)
	data := make([]byte, 0)
	var ui string

	for {
		n, err := s.read(buff)
		if err != nil {
			log.Error(err)
			continue
		}
		//spew.Dump(buff[0:n])

		log.Tracef("read %d", n)

		if !s.initialized {
			idx := s.findInitialPoint(buff[0:n])
			if idx < 0 { // could not find initial point
				continue
			}
			log.WithFields(logrus.Fields{
				"data": string(buff[0:n]),
			}).Debugf("found initial point=%d", idx)

			// found initial point
			data = make([]byte, 0)
			if n > idx+1 {
				data = append(data, buff[idx+1:]...)
			}
			s.initialized = true
			continue
		}

		for i := 0; i < n; i++ {
			data = append(data, buff[i])
			if buff[i] == NewLine {
				//fmt.Printf("%v\n", data)
				//spew.Dump(data)

				if ui != string(data) {
					spew.Dump(data)
					ui = string(data)
				}
				//spew.du
				data = make([]byte, 0)
				continue
			}
		}

		select {
		case <-s.Ctx.Done(): // for gracefully shutdown
			s.Log.Debug("server canceled; no longer works")
			return nil
		case <-time.After(100 * time.Millisecond):
		}
	}
}

func (s *Server) findInitialPoint(buff []byte) int {
	for i := 0; i < len(buff); i++ {
		if buff[i] == NewLine {
			return i
		}
	}
	return -1
}

func (s *Server) read(buff []byte) (int, error) {
	count := rand.Intn(len(s.sampleData)) + 1

	for i := 0; i < count; i++ {
		p := (i + s.idx) % s.sampleDataLen

		buff[i] = s.sampleData[p]
		if s.sampleData[p] == NewLine {
			count = i + 1
			s.idx = (s.idx + count) % s.sampleDataLen
			return count, nil
		}
	}

	s.idx = (s.idx + count) % s.sampleDataLen
	return count, nil
}

func (s *Server) read2() []byte {
	count := rand.Intn(len(s.sampleData)) + 1
	b := make([]byte, count)
	for i := 0; i < count; i++ {
		p := (i + s.idx) % s.sampleDataLen

		b[i] = s.sampleData[p]
		if s.sampleData[p] == NewLine {
			count = i + 1
			s.idx = (s.idx + count) % s.sampleDataLen
			return b[0:count]
		}
	}

	s.idx = (s.idx + count) % s.sampleDataLen
	return b
}

func (s *Server) Stop() error {
	s.Log.Debug("server has been stopped")
	return nil
}
