package scale

import (
	"bytes"
	"github.com/bugst/go-serial"
	"github.com/devplayg/hippo"
	"github.com/sirupsen/logrus"
	"math/rand"
)

var (
	log           *logrus.Logger
	sampleData    = []byte("ABCDEFGHIJKLMNOP\r\n")
	sampleDataLen = len(sampleData)
)

const (
	NewLine  byte = 0x0a
	BuffSize      = 32
)

type Server struct {
	hippo.Launcher // DO NOT REMOVE; Launcher links servers and engine each other.
	idx            int
	initialized    bool
	controller     *Controller
	port           serial.Port
	lastValue      []byte
}

func NewServer() *Server {
	return &Server{
		idx: rand.Intn(sampleDataLen),
	}
}

func (s *Server) init() error {
	mode := &serial.Mode{
		BaudRate: 9600,
		DataBits: 8,
		StopBits: serial.TwoStopBits,
		Parity:   serial.NoParity,
	}

	port, err := serial.Open("COM3", mode)
	if err != nil {
		return err
	}
	s.port = port
	return nil
}

func (s *Server) Start() error {
	log = s.Log
	s.controller = NewController(":8000")
	s.controller.Start()

	if err := s.init(); err != nil {
		return err
	}

	buff := make([]byte, BuffSize)
	data := make([]byte, 0)
	//var ui []byte

	for {
		n, err := s.port.Read(buff)
		if err != nil {
			return err
		}
		if n == 0 {
			log.Debug("EOF")
			break
		}

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
				if !bytes.Equal(s.lastValue, data) {
					s.controller.hub.broadcast <- append([]byte("Measure: "), data...)

					// Deep copy
					s.lastValue = make([]byte, len(data))
					copy(s.lastValue, data)
				}
				data = make([]byte, 0)
				continue
			}
		}
	}

	return nil
}

//for {
//	n, err := s.read(buff)
//	if err != nil {
//		log.Error(err)
//		continue
//	}
//	//spew.Dump(buff[0:n])
//
//	log.Tracef("read %d", n)
//
//	if !s.initialized {
//		idx := s.findInitialPoint(buff[0:n])
//		if idx < 0 { // could not find initial point
//			continue
//		}
//		log.WithFields(logrus.Fields{
//			"data": string(buff[0:n]),
//		}).Debugf("found initial point=%d", idx)
//
//		// found initial point
//		data = make([]byte, 0)
//		if n > idx+1 {
//			data = append(data, buff[idx+1:]...)
//		}
//		s.initialized = true
//		continue
//	}
//
//	for i := 0; i < n; i++ {
//		data = append(data, buff[i])
//		if buff[i] == NewLine {
//			if !bytes.Equal(ui, data) {
//				s.controller.hub.broadcast <- append([]byte("Measure: "), data...)
//
//				// Deep copy
//				ui = make([]byte, len(data))
//				copy(ui, data)
//			}
//			data = make([]byte, 0)
//			continue
//		}
//	}
//
//	select {
//	case <-s.Ctx.Done(): // for gracefully shutdown
//		s.Log.Debug("server canceled; no longer works")
//		return nil
//	case <-time.After(100 * time.Millisecond):
//	}
//}
//}

//func (s *Server) Start() error {
//	log = s.Log
//	s.controller = NewController(":8000")
//	s.controller.Start()
//
//	buff := make([]byte, BuffSize)
//	data := make([]byte, 0)
//	var ui []byte
//
//	//go func() {
//	//	for {
//	//		if ui != nil {
//	//			ui[rand.Intn(5)] = byte(randomByte())
//	//			time.Sleep(5 * time.Second)
//	//		}
//	//	}
//	//}()
//
//	for {
//		n, err := s.read(buff)
//		if err != nil {
//			log.Error(err)
//			continue
//		}
//		//spew.Dump(buff[0:n])
//
//		log.Tracef("read %d", n)
//
//		if !s.initialized {
//			idx := s.findInitialPoint(buff[0:n])
//			if idx < 0 { // could not find initial point
//				continue
//			}
//			log.WithFields(logrus.Fields{
//				"data": string(buff[0:n]),
//			}).Debugf("found initial point=%d", idx)
//
//			// found initial point
//			data = make([]byte, 0)
//			if n > idx+1 {
//				data = append(data, buff[idx+1:]...)
//			}
//			s.initialized = true
//			continue
//		}
//
//		for i := 0; i < n; i++ {
//			data = append(data, buff[i])
//			if buff[i] == NewLine {
//				//fmt.Printf("%v\n", data)
//				//spew.Dump(data)
//
//				if !bytes.Equal(ui, data) {
//					//spew.Dump(ui)
//					//spew.Dump(data)
//					s.controller.hub.broadcast <- append([]byte("Measure: "), data...)
//
//					// Deep copy
//					ui = make([]byte, len(data))
//					copy(ui, data)
//				}
//				//if ui != string(data) {
//				//if ui != string(data) {
//				//}
//				//spew.du
//				data = make([]byte, 0)
//				continue
//			}
//		}
//
//		select {
//		case <-s.Ctx.Done(): // for gracefully shutdown
//			s.Log.Debug("server canceled; no longer works")
//			return nil
//		case <-time.After(100 * time.Millisecond):
//		}
//	}
//}

func (s *Server) findInitialPoint(buff []byte) int {
	for i := 0; i < len(buff); i++ {
		if buff[i] == NewLine {
			return i
		}
	}
	return -1
}

// generator
func (s *Server) read_for_test(buff []byte) (int, error) {
	count := rand.Intn(len(sampleData)) + 1

	for i := 0; i < count; i++ {
		p := (i + s.idx) % sampleDataLen

		buff[i] = sampleData[p]
		if sampleData[p] == NewLine {
			count = i + 1
			s.idx = (s.idx + count) % sampleDataLen
			return count, nil
		}
	}

	s.idx = (s.idx + count) % sampleDataLen
	return count, nil
}

func (s *Server) read2() []byte {
	count := rand.Intn(len(sampleData)) + 1
	b := make([]byte, count)
	for i := 0; i < count; i++ {
		p := (i + s.idx) % sampleDataLen

		b[i] = sampleData[p]
		if sampleData[p] == NewLine {
			count = i + 1
			s.idx = (s.idx + count) % sampleDataLen
			return b[0:count]
		}
	}

	s.idx = (s.idx + count) % sampleDataLen
	return b
}

func (s *Server) Stop() error {
	s.Log.Debug("server has been stopped")
	return nil
}
