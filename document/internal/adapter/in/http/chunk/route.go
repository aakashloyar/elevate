package http 

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("/mark-chunk", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost: 
		{
			h.MarkChunk(w,r)
		}
		default: {
			http.Error(w,"Method Not Allowed",http.StatusMethodNotAllowed)
		}
		}
	})
}
