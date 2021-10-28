package logger

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestLogWrite(t *testing.T) {
	type args struct {
		level   Level
		message string
		code    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Log info",
			args: args{
				level:   Info,
				message: "Some info",
				code:    "test_info",
			},
			want: `time="` + time.Now().Format("02-01-2006 15:04:05") + `" level=info msg="[test_info] Some info"`,
		},
		{
			name: "Log warning",
			args: args{
				level:   Warning,
				message: "Some warning",
				code:    "test_warning",
			},
			want: `time="` + time.Now().Format("02-01-2006 15:04:05") + `" level=warning msg="[test_warning] Some warning"`,
		},
		{
			name: "Log error",
			args: args{
				level:   Error,
				message: "Some error",
				code:    "test_error",
			},
			want: `time="` + time.Now().Format("02-01-2006 15:04:05") + `" level=error msg="[test_error] Some error"`,
		},
		{
			name: "Log unknown",
			args: args{
				level:   "Unknown",
				message: "Some unknown",
				code:    "test_unknown",
			},
			want: ``,
		},
	}

	for _, test := range tests {
		b := *bytes.NewBuffer([]byte{})
		logrus.SetOutput(&b)
		t.Run(test.name, func(t *testing.T) {
			log := NewLogger()
			log.Write(test.args.level, test.args.message, test.args.code)
		})

		require.Equal(t, test.want, strings.TrimRight(b.String(), "\n"))
	}
}
