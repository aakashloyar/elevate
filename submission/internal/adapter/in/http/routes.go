package http

import (
	"net/http"
	"strings"
)

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	//health check route
	mux.HandleFunc("/submissions/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	//create submission route
	mux.HandleFunc("/submissions", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.CreateSubmission(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/submissions/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/submissions/")
		parts := strings.Split(path, "/")
		submissionID := parts[0]
		if submissionID == "" {
			http.Error(w, "missing submission id", http.StatusBadRequest)
			return
		}

		if len(parts) == 1 {
			switch r.Method {
			//get submission by id route
			case http.MethodGet:
				h.GetSubmissionByID(w, r, submissionID)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		if len(parts) == 2 && parts[1] == "save" {
			switch r.Method {
			//save answer route
			case http.MethodPost:
				h.SaveAnswer(w, r, submissionID)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		if len(parts) == 2 && parts[1] == "status" {
			switch r.Method {
			//get submission status route
			case http.MethodGet:
				h.GetSubmissionStatus(w, r, submissionID)
			//update submission status route
			case http.MethodPatch:
				h.UpdateSubmissionStatus(w, r, submissionID)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		if len(parts) == 2 && parts[1] == "start" {
			switch r.Method {
			//start submission route
			case http.MethodPost:
				h.StartSubmission(w, r, submissionID)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		if len(parts) == 2 && parts[1] == "submit" {
			switch r.Method {
			//submit submission route
			case http.MethodPost:
				h.SubmitSubmission(w, r, submissionID)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		if len(parts) == 3 && parts[1] == "save" && parts[2] == "batch" {
			switch r.Method {
			//save answer batch route
			case http.MethodPost:
				h.SaveAnswerBatch(w, r, submissionID)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		http.NotFound(w, r)
	})
}

//1. create submission route
//2. get submission by id route
//3. save answer route
//4. get submission status route
//5. start submission route
//6. submit submission route
//7. save answer batch route
