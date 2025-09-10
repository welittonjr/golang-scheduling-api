package adapter

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	http "scheduling/internal/infra/gin"

	"github.com/gin-gonic/gin"
)

func TestGinContext_Param(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name     string
		path     string
		url      string
		param    string
		expected string
	}{
		{
			name:     "valid param",
			path:     "/test/:id",
			url:      "/test/123",
			param:    "id",
			expected: "123",
		},
		{
			name:     "non-existent param",
			path:     "/test/:id",
			url:      "/test/123",
			param:    "name",
			expected: "",
		},
		{
			name:     "multiple params",
			path:     "/test/:id/:name",
			url:      "/test/456/john",
			param:    "id",
			expected: "456",
		},
		{
			name:     "get second param",
			path:     "/test/:id/:name",
			url:      "/test/456/john",
			param:    "name",
			expected: "john",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			
			router.GET(tt.path, func(c *gin.Context) {
				ginCtx := &GinContext{ctx: c}
				result := ginCtx.Param(tt.param)
				c.String(200, result)
			})
			
			req := httptest.NewRequest("GET", tt.url, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			
			if w.Code == 404 {
				t.Errorf("route not found for URL %s", tt.url)
				return
			}
			
			if w.Body.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, w.Body.String())
			}
		})
	}
}

func TestGinContext_Query(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name     string
		url      string
		key      string
		expected string
	}{
		{
			name:     "valid query param",
			url:      "/test?name=john",
			key:      "name",
			expected: "john",
		},
		{
			name:     "empty query param",
			url:      "/test?name=",
			key:      "name",
			expected: "",
		},
		{
			name:     "missing query param",
			url:      "/test",
			key:      "name",
			expected: "",
		},
		{
			name:     "multiple query params",
			url:      "/test?name=john&age=30",
			key:      "age",
			expected: "30",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			
			router.GET("/test", func(c *gin.Context) {
				ginCtx := &GinContext{ctx: c}
				result := ginCtx.Query(tt.key)
				c.String(200, result)
			})
			
			req := httptest.NewRequest("GET", tt.url, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			
			if w.Body.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, w.Body.String())
			}
		})
	}
}

func TestGinContext_Bind(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	type TestData struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	
	tests := []struct {
		name        string
		jsonData    string
		expectError bool
		expected    TestData
	}{
		{
			name:        "valid json",
			jsonData:    `{"name":"John","age":30}`,
			expectError: false,
			expected:    TestData{Name: "John", Age: 30},
		},
		{
			name:        "empty json",
			jsonData:    `{}`,
			expectError: false,
			expected:    TestData{Name: "", Age: 0},
		},
		{
			name:        "invalid json",
			jsonData:    `{"name":"John","age":}`,
			expectError: true,
			expected:    TestData{},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			
			router.POST("/test", func(c *gin.Context) {
				ginCtx := &GinContext{ctx: c}
				var data TestData
				err := ginCtx.Bind(&data)
				
				if tt.expectError {
					if err == nil {
						t.Error("expected error but got none")
					}
					c.JSON(400, gin.H{"error": "binding failed"})
					return
				}
				
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					c.JSON(400, gin.H{"error": err.Error()})
					return
				}
				
				if data.Name != tt.expected.Name || data.Age != tt.expected.Age {
					t.Errorf("expected %+v, got %+v", tt.expected, data)
				}
				
				c.JSON(200, data)
			})
			
			req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(tt.jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			
			if tt.expectError {
				if w.Code != 400 {
					t.Errorf("expected status 400, got %d", w.Code)
				}
			} else {
				if w.Code != 200 {
					t.Errorf("expected status 200, got %d", w.Code)
				}
			}
		})
	}
}

func TestGinContext_JSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name   string
		status int
		data   interface{}
	}{
		{
			name:   "string data",
			status: 200,
			data:   map[string]string{"message": "hello"},
		},
		{
			name:   "number data",
			status: 201,
			data:   map[string]int{"count": 42},
		},
		{
			name:   "error status",
			status: 500,
			data:   map[string]string{"error": "internal error"},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			
			router.GET("/test", func(c *gin.Context) {
				ginCtx := &GinContext{ctx: c}
				err := ginCtx.JSON(tt.status, tt.data)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			})
			
			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			
			if w.Code != tt.status {
				t.Errorf("expected status %d, got %d", tt.status, w.Code)
			}
			
			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json; charset=utf-8" {
				t.Errorf("expected Content-Type application/json; charset=utf-8, got %s", contentType)
			}
		})
	}
}

func TestGinContext_Status(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name   string
		status int
	}{
		{"ok", 200},
		{"created", 201},
		{"bad request", 400},
		{"not found", 404},
		{"internal error", 500},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			
			router.GET("/test", func(c *gin.Context) {
				ginCtx := &GinContext{ctx: c}
				ginCtx.Status(tt.status)
			})
			
			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			
			if w.Code != tt.status {
				t.Errorf("expected status %d, got %d", tt.status, w.Code)
			}
		})
	}
}

func TestGinContext_Header(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name  string
		key   string
		value string
	}{
		{"custom header", "X-Custom-Header", "custom-value"},
		{"content type", "Content-Type", "text/plain"},
		{"authorization", "Authorization", "Bearer token123"},
		{"empty value", "X-Empty", ""},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			
			router.GET("/test", func(c *gin.Context) {
				ginCtx := &GinContext{ctx: c}
				ginCtx.Header(tt.key, tt.value)
				ginCtx.Status(200)
			})
			
			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			
			headerValue := w.Header().Get(tt.key)
			if headerValue != tt.value {
				t.Errorf("expected header %s: %s, got %s", tt.key, tt.value, headerValue)
			}
		})
	}
}

func TestWrapHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name     string
		path     string
		url      string
		paramKey string
		expected string
	}{
		{
			name:     "simple handler",
			path:     "/test/:id",
			url:      "/test/123",
			paramKey: "id",
			expected: "123",
		},
		{
			name:     "handler with different param",
			path:     "/user/:userId",
			url:      "/user/456",
			paramKey: "userId",
			expected: "456",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerCalled := false
			testHandler := func(ctx http.Context) error {
				handlerCalled = true
				param := ctx.Param(tt.paramKey)
				return ctx.JSON(200, map[string]string{"id": param})
			}
			
			router := gin.New()
			router.GET(tt.path, wrapHandler(testHandler))
			
			req := httptest.NewRequest("GET", tt.url, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			
			if !handlerCalled {
				t.Error("handler was not called")
			}
			
			if w.Code != 200 {
				t.Errorf("expected status 200, got %d", w.Code)
			}
			
			var response map[string]string
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Errorf("failed to unmarshal response: %v", err)
			}
			
			if response["id"] != tt.expected {
				t.Errorf("expected id %s, got %s", tt.expected, response["id"])
			}
		})
	}
}

func TestWrapMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name        string
		headerKey   string
		headerValue string
	}{
		{"trace middleware", "X-Trace-ID", "trace123"},
		{"auth middleware", "X-Auth", "authenticated"},
		{"custom middleware", "X-Custom", "middleware-applied"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middlewareCalled := false
			testMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
				return func(ctx http.Context) error {
					middlewareCalled = true
					ctx.Header(tt.headerKey, tt.headerValue)
					return next(ctx)
				}
			}
			
			router := gin.New()
			router.Use(wrapMiddleware(testMiddleware))
			router.GET("/test", func(c *gin.Context) {
				c.Status(200)
			})
			
			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			
			if !middlewareCalled {
				t.Error("middleware was not called")
			}
			
			headerValue := w.Header().Get(tt.headerKey)
			if headerValue != tt.headerValue {
				t.Errorf("expected header %s: %s, got %s", tt.headerKey, tt.headerValue, headerValue)
			}
		})
	}
}