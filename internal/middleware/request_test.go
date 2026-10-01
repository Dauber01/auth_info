package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"

	"auth_info/internal/pkg/trace"
)

func TestTraceAccessAndFailureEnvelope(t *testing.T) {
	gin.SetMode(gin.TestMode)
	core, logs := observer.New(zap.DebugLevel)
	log := zap.New(core)
	r := gin.New()
	r.Use(TraceID(), AccessLog(log), Recovery(log), ErrorHandler(log))
	r.GET("/failure", func(c *gin.Context) {
		if trace.ID(c.Request.Context()) != "request-42" {
			t.Error("missing context ID")
		}
		_ = c.Error(errors.New("fixture error"))
	})
	req := httptest.NewRequest("GET", "/failure?token=never-log-this", strings.NewReader("never-log-body"))
	req.Header.Set(trace.Header, "request-42")
	req.Header.Set("Authorization", "never-log-authorization")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if w.Code != 500 || body["code"] != float64(500) || w.Header().Get(trace.Header) != "request-42" {
		t.Fatal(w.Body.String())
	}
	if len(logs.FilterMessage("request failed").All()) != 1 {
		t.Fatal("duplicate error logging")
	}
	for _, entry := range logs.All() {
		data, _ := json.Marshal(entry.ContextMap())
		s := string(data)
		if strings.Contains(s, "never-log") || entry.ContextMap()["trace_id"] != "request-42" {
			t.Fatal(s)
		}
	}
}

func TestTimeoutPanicAndAlreadyWrittenResponses(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(TraceID(), Recovery(zap.NewNop()), ErrorHandler(zap.NewNop()))
	r.GET("/timeout", RequestTimeout(time.Millisecond), func(c *gin.Context) { <-c.Request.Context().Done() })
	r.GET("/panic", func(*gin.Context) { panic("fixture") })
	r.GET("/written", RequestTimeout(time.Millisecond), func(c *gin.Context) {
		c.String(202, "accepted")
		<-c.Request.Context().Done()
	})
	for path, status := range map[string]int{"/timeout": 504, "/panic": 500, "/written": 202} {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest("GET", path, nil))
		if w.Code != status || w.Header().Get(trace.Header) == "" {
			t.Fatalf("%s: %d %s", path, w.Code, w.Body.String())
		}
		if path == "/written" && w.Body.String() != "accepted" {
			t.Fatal("response overwritten")
		}
	}
}

func TestCanceledErrorMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(ErrorHandler(zap.NewNop()))
	r.GET("/", func(c *gin.Context) { _ = c.Error(context.Canceled) })
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))
	if w.Code != 408 {
		t.Fatal(w.Code)
	}
}
