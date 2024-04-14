package tests

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPickUpPointController_LoggingMiddleware(t *testing.T) {
	Test1Method := "GET"
	Test1URI := "/test"
	req := httptest.NewRequest(Test1Method, Test1URI, nil)
	w := httptest.NewRecorder()

	type args struct {
		handler http.Handler
	}
	tests := []struct {
		name  string
		args  args
		want1 string
		want2 string
	}{
		{
			name:  "smoke test",
			args:  args{dummyHandler()},
			want1: Test1Method,
			want2: Test1URI,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pickUpPointController.LoggingMiddleware(tt.args.handler).ServeHTTP(w, req)

			lm := getKafkaMessage(t)

			assert.Equal(t, lm.Method, tt.want1)
			assert.Equal(t, lm.URI, tt.want2)
		})
	}
}
