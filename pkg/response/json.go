package response

import (
	"encoding/json"
	"net/http"
)

type Response interface {
	PrettyPrint(http.ResponseWriter, interface{}, ...Options)
}

type Resp struct{}

func (r *Resp) PrettyPrint(w http.ResponseWriter, o interface{}, i ...Options) {
	option := &Opt{}
	for key, _ := range i {
		i[key](option)
	}

	for key, value := range option.Headers {
		w.Header().Set(key, value)
	}
	w.WriteHeader(option.Code)

	_ = json.NewEncoder(w).Encode(o)
}

func NewPrint() Response {
	return &Resp{}
}
