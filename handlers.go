package sm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"

	// Routing
	"github.com/julienschmidt/httprouter"
)

// IndexHandler ...
func (s *Server) IndexHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		metrics.CounterVec("server", "requests").WithLabelValues("GET", "/").Inc()
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "SEV Manager %s", FullVersion())
	}
}

// SearchHandler ...
func (s *Server) SearchHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var (
			err    error
			events []*Event
		)

		metrics.CounterVec("server", "requests").WithLabelValues("GET", "/search").Inc()

		qs := r.URL.Query()

		if id := ParseId(p.ByName("id")); id > 0 {
			events, err = db.Find(id)
			if err != nil {
				http.Error(w, "Not Found", http.StatusNotFound)
				return
			}
		} else if q := qs.Get("q"); q != "" {
			events, err = db.Search(q)
			if err != nil {
				http.Error(w, "Not Found", http.StatusNotFound)
				return
			}
		} else {
			events, err = db.All()
			if err != nil {
				http.Error(w, "Not Found", http.StatusNotFound)
				return
			}
		}

		out, err := json.Marshal(events)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
	}
}

// CreateHandler ...
func (s *Server) CreateHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		metrics.CounterVec("server", "requests").WithLabelValues("POST", "/create").Inc()

		qs := r.URL.Query()

		title := qs.Get("title")
		if title == "" {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		level := SafeParseInt(qs.Get("level"), DefaultSEVLevel)

		event, err := NewEvent(title, level)
		if err != nil {
			log.Errorf("error creating new event: %s", err)
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}

		u, err := url.Parse(fmt.Sprintf("/search/%d", event.ID))
		if err != nil {
			http.Error(w, "Internal Error", http.StatusInternalServerError)
		}
		http.Redirect(w, r, r.URL.ResolveReference(u).String(), http.StatusFound)
	}
}

// CloseHandler ...
func (s *Server) CloseHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		metrics.CounterVec("server", "requests").WithLabelValues("POST", "/close").Inc()

		id := ParseId(p.ByName("id"))

		if id <= 0 {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		event, err := db.Get(id)
		if err != nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		err = event.Close()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}
