package http

import (
	"net/http"
	"strings"
)

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("/assessments/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/assessments", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.ListAssessments(w, r)
		case http.MethodPost:
			h.CreateAssessment(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/assessments/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.Trim(strings.TrimPrefix(r.URL.Path, "/assessments/"), "/")
		parts := strings.Split(path, "/")
		assessmentID := parts[0]
		if assessmentID == "" {
			http.Error(w, "missing assessment id", http.StatusBadRequest)
			return
		}
		if len(parts) == 1 {
			switch r.Method {
			case http.MethodGet:
				h.GetAssessmentByID(w, r, assessmentID)
			case http.MethodDelete:
				h.DeleteAssessment(w, r, assessmentID)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		if len(parts) == 2 && parts[1] == "problem" {
			switch r.Method {
			case http.MethodPost:
				h.CreateAssessmentProblem(w, r, assessmentID)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		if len(parts) == 2 && parts[1] == "problems" {
			if r.Method != http.MethodGet {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			h.GetAssessmentProblems(w, r, assessmentID)
			return
		}

		if len(parts) == 2 && parts[1] == "marking-scheme" {
			switch r.Method {
			case http.MethodPost:
				h.CreateAssessmentMarkingScheme(w, r, assessmentID)
			case http.MethodGet:
				h.GetAssessmentMarkingScheme(w, r, assessmentID)
			case http.MethodPut:
				h.PutAssessmentMarkingScheme(w, r, assessmentID)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		http.NotFound(w, r)
	})
}
