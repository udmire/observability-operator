package http_test

import (
	"html/template"
	"math/rand"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-kit/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/udmire/observability-operator/pkg/utils/http"
	"gopkg.in/yaml.v3"
)

func TestRenderHTTPResponse(t *testing.T) {
	type testStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	tests := []struct {
		name                string
		headers             map[string]string
		tmpl                string
		expectedOutput      string
		expectedContentType string
		value               testStruct
	}{
		{
			name: "Test Renders json",
			headers: map[string]string{
				"Accept": "application/json",
			},
			tmpl:                "<html></html>",
			expectedOutput:      `{"name":"testName","value":42}`,
			expectedContentType: "application/json",
			value: testStruct{
				Name:  "testName",
				Value: 42,
			},
		},
		{
			name:                "Test Renders html",
			headers:             map[string]string{},
			tmpl:                "<html>{{ .Name }}</html>",
			expectedOutput:      "<html>testName</html>",
			expectedContentType: "text/html; charset=utf-8",
			value: testStruct{
				Name:  "testName",
				Value: 42,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl := template.Must(template.New("webpage").Parse(tt.tmpl))
			writer := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/", nil)

			for k, v := range tt.headers {
				request.Header.Add(k, v)
			}

			http.RenderHTTPResponse(writer, tt.value, tmpl, request)

			assert.Equal(t, tt.expectedContentType, writer.Header().Get("Content-Type"))
			assert.Equal(t, 200, writer.Code)
			assert.Equal(t, tt.expectedOutput, writer.Body.String())
		})
	}
}

func TestWriteTextResponse(t *testing.T) {
	w := httptest.NewRecorder()

	http.WriteTextResponse(w, "hello world")

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "hello world", w.Body.String())
	assert.Equal(t, "text/plain", w.Header().Get("Content-Type"))
}

func TestStreamWriteYAMLResponse(t *testing.T) {
	type testStruct struct {
		Name  string `yaml:"name"`
		Value int    `yaml:"value"`
	}
	tt := struct {
		name                string
		headers             map[string]string
		expectedOutput      string
		expectedContentType string
		value               map[string]*testStruct
	}{
		name: "Test Stream Render YAML",
		headers: map[string]string{
			"Content-Type": "application/yaml",
		},
		expectedContentType: "application/yaml",
		value:               make(map[string]*testStruct),
	}

	// Generate some data to serialize.
	for i := 0; i < rand.Intn(100)+1; i++ {
		ts := testStruct{
			Name:  "testName" + strconv.Itoa(i),
			Value: i,
		}
		tt.value[ts.Name] = &ts
	}
	d, err := yaml.Marshal(tt.value)
	require.NoError(t, err)
	tt.expectedOutput = string(d)
	w := httptest.NewRecorder()

	done := make(chan struct{})
	iter := make(chan interface{})
	go func() {
		http.StreamWriteYAMLResponse(w, iter, log.NewNopLogger())
		close(done)
	}()
	for k, v := range tt.value {
		iter <- map[string]*testStruct{k: v}
	}
	close(iter)
	<-done
	assert.Equal(t, tt.expectedContentType, w.Header().Get("Content-Type"))
	assert.Equal(t, 200, w.Code)
	assert.YAMLEq(t, tt.expectedOutput, w.Body.String())
}
