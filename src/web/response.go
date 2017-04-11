package web

import (
	"encoding/json"
	"io"
	"net/http"
)

type Response struct {
	OK      bool        `json:"ok"`
	Message string      `json:"message"`
	Status  int         `json:"status"`
	Data    interface{} `json:"data,omitempty"`
}

func (r Response) EncodeTo(w io.Writer) error {
	httpWriter, ok := w.(http.ResponseWriter)
	if ok {
		if r.OK {
			if r.Status != 0 {
				httpWriter.WriteHeader(r.Status)
			} else {
				r.Status = http.StatusOK
				httpWriter.WriteHeader(http.StatusOK)
			}
		}
	}
	if r.Message == "" && r.Status != 0 {
		r.Message = http.StatusText(r.Status)
	}
	return json.NewEncoder(w).Encode(r)
}

func ErrorResponse(w http.ResponseWriter, err error, status int) {
	r := Response{
		Status:  status,
		OK:      false,
		Message: err.Error(),
	}
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(r); err != nil {
		http.Error(w, err.Error(), status)
	}
}

// OK is shortcut for web.Response{OK: true, Data: data}.EncodeTo(w)
func OK(w io.Writer, data interface{}) {
	Response{
		OK:   true,
		Data: data,
	}.EncodeTo(w)
}
