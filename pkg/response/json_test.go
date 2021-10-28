package response

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestResp_PrettyPrint(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		o interface{}
		i []Options
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Pretty string",
			args: args{
				w: httptest.NewRecorder(),
				o: "some text",
			},
			want: `"some text"`,
		},
		{
			name: "Pretty struct",
			args: args{
				w: httptest.NewRecorder(),
				o: struct {
					Name string
					Age  uint8
				}{
					Name: "Alex",
					Age:  100,
				},
			},
			want: `{"Name":"Alex","Age":100}`,
		},
		{
			name: "Pretty map",
			args: args{
				w: httptest.NewRecorder(),
				o: map[string]string{
					"Name": "Alex",
				},
			},
			want: `{"Name":"Alex"}`,
		},
		{
			name: "Pretty struct with status code",
			args: args{
				w: httptest.NewRecorder(),
				o: struct {
					Name string
					Age  uint8
				}{
					Name: "Alex",
					Age:  100,
				},
				i: []Options{
					WithCode(http.StatusBadRequest),
					WithHeaders(map[string]string{"header": "value"}),
				},
			},
			want: `{"Name":"Alex","Age":100}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := NewPrint()
			w := httptest.NewRecorder()
			p.PrettyPrint(w, test.args.o, test.args.i...)
			expect := w.Body.String()
			expect = strings.TrimRight(expect, "\n")

			if expect != test.want {
				t.Errorf("Failed " + test.name + "\nExpected: " + test.want + "\nGot: " + expect)
			}

			if test.args.i != nil && w.Result().Header.Get("header") != "value" {
				t.Errorf("Failed " + test.name + "\nInvalid header")
			}

			if test.args.i != nil && w.Result().StatusCode != http.StatusBadRequest {
				t.Errorf("Failed " + test.name + "\nInvalid status code")
			}
		})
	}
}
