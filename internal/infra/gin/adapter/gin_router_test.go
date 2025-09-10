package adapter

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	http "scheduling/internal/infra/gin"

	"github.com/gin-gonic/gin"
)

func TestNewRouter(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"create new router"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := NewRouter()
			if router == nil {
				t.Error("expected non-nil router")
			}
			
			ginRouter, ok := router.(*ginRouter)
			if !ok {
				t.Error("expected *ginRouter type")
			}
			
			if ginRouter.group == nil {
				t.Error("expected non-nil group")
			}
			
			if ginRouter.root == nil {
				t.Error("expected non-nil root")
			}
		})
	}
}

func TestGinRouter_HTTPMethods(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name   string
		method string
		path   string
		url    string
		body   string
	}{
		{"GET request", "GET", "/test", "/test", ""},
		{"POST request", "POST", "/test", "/test", `{"data":"test"}`},
		{"PUT request", "PUT", "/test", "/test", `{"data":"update"}`},
		{"DELETE request", "DELETE", "/test", "/test", ""},
		{"GET with param", "GET", "/test/:id", "/test/123", ""},
		{"POST with param", "POST", "/user/:id", "/user/456", `{"name":"john"}`},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := NewRouter()
			handlerCalled := false
			
			testHandler := func(ctx http.Context) error {
				handlerCalled = true
				return ctx.JSON(200, map[string]string{"method": tt.method})
			}
			
			switch tt.method {
			case "GET":
				router.GET(tt.path, testHandler)
			case "POST":
				router.POST(tt.path, testHandler)
			case "PUT":
				router.PUT(tt.path, testHandler)
			case "DELETE":
				router.DELETE(tt.path, testHandler)
			}
			
			ginRouter := router.(*ginRouter)
			req := httptest.NewRequest(tt.method, tt.url, strings.NewReader(tt.body))
			if tt.body != "" {
				req.Header.Set("Content-Type", "application/json")
			}
			w := httptest.NewRecorder()
			ginRouter.root.ServeHTTP(w, req)
			
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
			
			if response["method"] != tt.method {
				t.Errorf("expected method %s, got %s", tt.method, response["method"])
			}
		})
	}
}

func TestGinRouter_Use(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name            string
		middlewareCount int
		expectedHeaders map[string]string
	}{
		{
			name:            "single middleware",
			middlewareCount: 1,
			expectedHeaders: map[string]string{"X-Middleware-1": "applied"},
		},
		{
			name:            "multiple middlewares",
			middlewareCount: 3,
			expectedHeaders: map[string]string{
				"X-Middleware-1": "applied",
				"X-Middleware-2": "applied", 
				"X-Middleware-3": "applied",
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := NewRouter()
			middlewaresCalled := make([]bool, tt.middlewareCount)
			
			var middlewares []http.MiddlewareFunc
			for i := 0; i < tt.middlewareCount; i++ {
				index := i
				headerKey := "X-Middleware-" + string(rune('1'+i))
				middleware := func(next http.HandlerFunc) http.HandlerFunc {
					return func(ctx http.Context) error {
						middlewaresCalled[index] = true
						ctx.Header(headerKey, "applied")
						return next(ctx)
					}
				}
				middlewares = append(middlewares, middleware)
			}
			
			handlerCalled := false
			testHandler := func(ctx http.Context) error {
				handlerCalled = true
				return ctx.JSON(200, map[string]string{"message": "ok"})
			}
			
			router.Use(middlewares...)
			router.GET("/test", testHandler)
			
			ginRouter := router.(*ginRouter)
			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()
			ginRouter.root.ServeHTTP(w, req)
			
			if !handlerCalled {
				t.Error("handler was not called")
			}
			
			for i, called := range middlewaresCalled {
				if !called {
					t.Errorf("middleware %d was not called", i+1)
				}
			}
			
			for headerKey, expectedValue := range tt.expectedHeaders {
				actualValue := w.Header().Get(headerKey)
				if actualValue != expectedValue {
					t.Errorf("expected header %s: %s, got %s", headerKey, expectedValue, actualValue)
				}
			}
		})
	}
}

func TestGinRouter_Group(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name        string
		groupPath   string
		handlerPath string
		fullURL     string
	}{
		{
			name:        "simple group",
			groupPath:   "/api",
			handlerPath: "/test",
			fullURL:     "/api/test",
		},
		{
			name:        "versioned api group",
			groupPath:   "/api/v1",
			handlerPath: "/users",
			fullURL:     "/api/v1/users",
		},
		{
			name:        "admin group",
			groupPath:   "/admin",
			handlerPath: "/dashboard",
			fullURL:     "/admin/dashboard",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := NewRouter()
			group := router.Group(tt.groupPath)
			
			if group == nil {
				t.Error("expected non-nil group")
			}
			
			ginGroup, ok := group.(*ginRouter)
			if !ok {
				t.Error("expected *ginRouter type")
			}
			
			if ginGroup.group == nil {
				t.Error("expected non-nil group.group")
			}
			
			if ginGroup.root == nil {
				t.Error("expected non-nil group.root")
			}
			
			handlerCalled := false
			testHandler := func(ctx http.Context) error {
				handlerCalled = true
				return ctx.JSON(200, map[string]string{"group": "test"})
			}
			
			group.GET(tt.handlerPath, testHandler)
			
			ginRouter := router.(*ginRouter)
			req := httptest.NewRequest("GET", tt.fullURL, nil)
			w := httptest.NewRecorder()
			ginRouter.root.ServeHTTP(w, req)
			
			if !handlerCalled {
				t.Error("handler was not called")
			}
			
			if w.Code != 200 {
				t.Errorf("expected status 200, got %d", w.Code)
			}
		})
	}
}

func TestGinRouter_Group_Nested(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name         string
		firstGroup   string
		secondGroup  string
		handlerPath  string
		fullURL      string
	}{
		{
			name:         "api version nested",
			firstGroup:   "/api",
			secondGroup:  "/v1",
			handlerPath:  "/users",
			fullURL:      "/api/v1/users",
		},
		{
			name:         "admin section nested",
			firstGroup:   "/admin",
			secondGroup:  "/settings",
			handlerPath:  "/config",
			fullURL:      "/admin/settings/config",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := NewRouter()
			firstGroup := router.Group(tt.firstGroup)
			secondGroup := firstGroup.Group(tt.secondGroup)
			
			if secondGroup == nil {
				t.Error("expected non-nil nested group")
			}
			
			handlerCalled := false
			testHandler := func(ctx http.Context) error {
				handlerCalled = true
				return ctx.JSON(200, map[string]string{"nested": "group"})
			}
			
			secondGroup.GET(tt.handlerPath, testHandler)
			
			ginRouter := router.(*ginRouter)
			req := httptest.NewRequest("GET", tt.fullURL, nil)
			w := httptest.NewRecorder()
			ginRouter.root.ServeHTTP(w, req)
			
			if !handlerCalled {
				t.Error("handler was not called")
			}
			
			if w.Code != 200 {
				t.Errorf("expected status 200, got %d", w.Code)
			}
		})
	}
}

func TestGinRouter_Run(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name        string
		hasRoot     bool
		expectError bool
	}{
		{
			name:        "valid router with root",
			hasRoot:     true,
			expectError: false,
		},
		{
			name:        "router without root",
			hasRoot:     false,
			expectError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var router *ginRouter
			
			if tt.hasRoot {
				r := NewRouter()
				router = r.(*ginRouter)
			} else {
				router = &ginRouter{root: nil}
			}
			
			if tt.expectError {
				err := router.Run(":8080")
				if err == nil {
					t.Error("expected error but got none")
				}
				
				if !strings.Contains(err.Error(), "no root engine to run") {
					t.Errorf("expected 'no root engine to run' error, got %v", err)
				}
			} else {
				if router.root == nil {
					t.Error("expected non-nil root for valid router")
				}
			}
		})
	}
}

func TestGinRouter_Group_WithMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name        string
		groupPath   string
		handlerPath string
		fullURL     string
		headerKey   string
		headerValue string
	}{
		{
			name:        "group with auth middleware",
			groupPath:   "/api",
			handlerPath: "/protected",
			fullURL:     "/api/protected",
			headerKey:   "X-Auth",
			headerValue: "required",
		},
		{
			name:        "admin group with logging",
			groupPath:   "/admin",
			handlerPath: "/logs",
			fullURL:     "/admin/logs",
			headerKey:   "X-Log",
			headerValue: "admin-access",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := NewRouter()
			
			middlewareCalled := false
			groupMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
				return func(ctx http.Context) error {
					middlewareCalled = true
					ctx.Header(tt.headerKey, tt.headerValue)
					return next(ctx)
				}
			}
			
			group := router.Group(tt.groupPath)
			group.Use(groupMiddleware)
			
			handlerCalled := false
			testHandler := func(ctx http.Context) error {
				handlerCalled = true
				return ctx.JSON(200, map[string]string{"message": "success"})
			}
			
			group.GET(tt.handlerPath, testHandler)
			
			ginRouter := router.(*ginRouter)
			req := httptest.NewRequest("GET", tt.fullURL, nil)
			w := httptest.NewRecorder()
			ginRouter.root.ServeHTTP(w, req)
			
			if !middlewareCalled {
				t.Error("group middleware was not called")
			}
			
			if !handlerCalled {
				t.Error("handler was not called")
			}
			
			headerValue := w.Header().Get(tt.headerKey)
			if headerValue != tt.headerValue {
				t.Errorf("expected header %s: %s, got %s", tt.headerKey, tt.headerValue, headerValue)
			}
		})
	}
}

func TestGinRouter_ComplexRouting(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name           string
		setupRoutes    func(router http.Router)
		requestMethod  string
		requestURL     string
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "multiple routes same path different methods",
			setupRoutes: func(router http.Router) {
				router.GET("/resource", func(ctx http.Context) error {
					return ctx.JSON(200, map[string]string{"action": "get"})
				})
				router.POST("/resource", func(ctx http.Context) error {
					return ctx.JSON(201, map[string]string{"action": "create"})
				})
			},
			requestMethod:  "GET",
			requestURL:     "/resource",
			expectedStatus: 200,
			expectedBody:   `{"action":"get"}`,
		},
		{
			name: "grouped routes with params",
			setupRoutes: func(router http.Router) {
				api := router.Group("/api/v1")
				api.GET("/users/:id", func(ctx http.Context) error {
					id := ctx.Param("id")
					return ctx.JSON(200, map[string]string{"userId": id})
				})
			},
			requestMethod:  "GET",
			requestURL:     "/api/v1/users/123",
			expectedStatus: 200,
			expectedBody:   `{"userId":"123"}`,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := NewRouter()
			tt.setupRoutes(router)
			
			ginRouter := router.(*ginRouter)
			req := httptest.NewRequest(tt.requestMethod, tt.requestURL, nil)
			w := httptest.NewRecorder()
			ginRouter.root.ServeHTTP(w, req)
			
			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
			
			actualBody := strings.TrimSpace(w.Body.String())
			if actualBody != tt.expectedBody {
				t.Errorf("expected body %s, got %s", tt.expectedBody, actualBody)
			}
		})
	}
}