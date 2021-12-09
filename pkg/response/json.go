package response

import (
	"io"
)

type Response interface {
	PrettyPrint(io.Writer, interface{}, ...Options)
}

type Resp struct{}

func (r *Resp) PrettyPrint(w io.Writer, o interface{}, i ...Options) {
	option := NewOptions()
	for key := range i {
		i[key](option)
	}

	//for key, value := range option.Headers {
	//	w.Header().Set(key, value)
	//}
	//
	//w.WriteHeader(option.Code)
	//_ = json.NewEncoder(w).Encode(o)
}

func NewPrint() Response {
	return &Resp{}
}
