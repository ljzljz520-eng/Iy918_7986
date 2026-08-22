package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"

	"guardpanel.local/guardpanel/internal/service"
)

type Server struct {
	service *service.Service
	mux     *http.ServeMux
}

func New(svc *service.Service) *Server {
	server := &Server{service: svc, mux: http.NewServeMux()}
	server.routes()
	return server
}

func (s *Server) routes() {
	s.mux.HandleFunc("/health", s.health)
	s.mux.HandleFunc("/records", s.records)
	s.mux.HandleFunc("/records/", s.record)
}
func (s *Server) Handler() http.Handler { return s.mux }

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) records(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		machine := r.URL.Query().Get("machine")
		status := r.URL.Query().Get("status")
		query := r.URL.Query().Get("q")
		records, err := s.service.Search(machine, status, query)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, records)
	case http.MethodPost:
		var input struct {
			ID        string   `json:"id"`
			MachineID string   `json:"machine_id"`
			Title     string   `json:"title"`
			Owner     string   `json:"owner"`
			Tags      []string `json:"tags"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			writeError(w, err)
			return
		}
		record, err := s.service.CreateRecord(input.ID, input.MachineID, input.Title, input.Owner, input.Tags)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, record)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) record(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/records/")
	if id == "" {
		writeError(w, http.ErrMissingFile)
		return
	}
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	action := r.URL.Query().Get("action")
	var err error
	switch action {
	case "review":
		_, err = s.service.SubmitReview(id, 0)
	case "approve":
		err = s.service.Approve(id)
	case "archive":
		err = s.service.Archive(id)
	case "publish":
		err = s.service.Publish(id)
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err != nil {
		writeError(w, err)
		return
	}
	record, err := s.service.Search("", "", "")
	if err != nil {
		writeError(w, err)
		return
	}
	for _, item := range record {
		if item.ID == id {
			writeJSON(w, http.StatusOK, item)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func writeError(w http.ResponseWriter, err error) {
	writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
}
