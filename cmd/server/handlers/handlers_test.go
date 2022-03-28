package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/gopherlearning/go-advanced-devops/cmd/server/storage"
	"github.com/gopherlearning/go-advanced-devops/internal/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_Update(t *testing.T) {
	type fields struct {
		s repositories.Repository
	}
	type want struct {
		statusCode int
		body       string
		value1     interface{}
		value2     interface{}
	}
	tests := []struct {
		name     string
		fields   fields
		content  string
		request1 string
		request2 string
		method   string
		want     want
	}{
		{
			name:     "Сохранение gauge",
			fields:   fields{s: storage.NewStorage()},
			content:  "text/plain",
			method:   http.MethodPost,
			request1: "/update/gauge/RandomValue/123.456",
			request2: "/update/gauge/RandomValue/123.457",
			want: want{
				statusCode: http.StatusOK,
				body:       "",
				value1:     float64(123.456),
				value2:     float64(123.457),
			},
		},
		{
			name:     "Сохранение counter",
			fields:   fields{s: storage.NewStorage()},
			content:  "text/plain",
			method:   http.MethodPost,
			request1: "/update/counter/PollCount/2",
			request2: "/update/counter/PollCount/3",
			want: want{
				statusCode: http.StatusOK,
				body:       "",
				value1:     []int{2},
				value2:     []int{2, 3},
			},
		},
		{
			name:     "Неправильный Content-Type",
			fields:   fields{s: storage.NewStorage()},
			content:  "",
			method:   http.MethodPost,
			request1: "/update/counter/PollCount/2",
			request2: "/update/counter/PollCount/3",
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "Only text/plain content are allowed!\n",
				value1:     nil,
				value2:     nil,
			},
		},
		{
			name:     "Неправильный http метод",
			fields:   fields{s: storage.NewStorage()},
			content:  "",
			method:   http.MethodGet,
			request1: "/update/counter/PollCount/2",
			request2: "/update/counter/PollCount/3",
			want: want{
				statusCode: http.StatusMethodNotAllowed,
				body:       "Only POST requests are allowed!\n",
				value1:     nil,
				value2:     nil,
			},
		},
		{
			name:     "Неподдерживаемый тип метрики",
			fields:   fields{s: storage.NewStorage()},
			content:  "text/plain",
			method:   http.MethodPost,
			request1: "/update/sumator/PollCount/2",
			request2: "/update/sumator/PollCount/3",
			want: want{
				statusCode: http.StatusAccepted,
				body:       "wrong metric format",
				value1:     nil,
				value2:     nil,
			},
		},
		{
			name:     "Сохранение неправильного counter",
			fields:   fields{s: storage.NewStorage()},
			content:  "text/plain",
			method:   http.MethodPost,
			request1: "/update/counter/PollCount/2.2",
			request2: "/update/counter/PollCount/3.1",
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "wrong metric type",
				value1:     []int{2},
				value2:     []int{2, 3},
			},
		},
	}
	var rMetricURL = regexp.MustCompile(`^.*\/(\w+)\/(\w+)\/([0-9\.]+)$`)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &Handler{
				s: tt.fields.s,
			}
			request := httptest.NewRequest(tt.method, tt.request1, nil)
			request.Header.Add("Content-Type", tt.content)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(handler.Update)
			h.ServeHTTP(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			body, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.body, string(body))

			if result.StatusCode == http.StatusOK {
				match := rMetricURL.FindStringSubmatch(tt.request1)
				assert.Equal(t, len(match), 4)
				assert.Contains(t, tt.fields.s.List()["192.0.2.1"], fmt.Sprintf("%s - %s - %v", match[1], match[2], tt.want.value1))
				fmt.Println(tt.name, tt.fields.s.List())
			}

			request = httptest.NewRequest(tt.method, tt.request2, nil)
			request.Header.Add("Content-Type", tt.content)
			w = httptest.NewRecorder()
			h = http.HandlerFunc(handler.Update)
			h.ServeHTTP(w, request)
			result = w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			body, err = ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.body, string(body))

			if result.StatusCode == http.StatusOK {
				match := rMetricURL.FindStringSubmatch(tt.request1)
				assert.Equal(t, len(match), 4)
				assert.Contains(t, tt.fields.s.List()["192.0.2.1"], fmt.Sprintf("%s - %s - %v", match[1], match[2], tt.want.value2))
				fmt.Println(tt.name, tt.fields.s.List())
			}
		})
	}
}
