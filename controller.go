package scale

import "net/http"

type Controller struct {
	addr string
	hub  *Hub
}

func NewController(addr string) *Controller {
	return &Controller{
		addr: addr,
	}
}

func (c *Controller) Start() error {
	c.hub = newHub()
	go c.hub.run()

	http.HandleFunc("/", c.serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(c.hub, w, r)
	})
	go func() {
		err := http.ListenAndServe(c.addr, nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()

	return nil
}

func (c *Controller) serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}
