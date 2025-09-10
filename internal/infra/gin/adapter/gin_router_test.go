package adapter

import (
	"bytes"
	"net/http/httptest"
	"testing"

	http "scheduling/internal/infra/gin"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewRouter(t *testing.T) {
	router := NewRouter()
	assert.NotNil(t, router)
	
	ginRouter, ok := router.(*ginRouter)
	assert.True(t, ok)
	assert.NotNil(t, ginRouter.group)
	assert.NotNil(t, ginRouter.root)
}

func TestGinRouter_GET(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := NewRouter()
	
	handlerCalled := false
	testHandler := func(ctx http.Context) error {
		handlerCalled = true
		return ctx.JSON(200, map[string]string{"method": "GET"})
	}
	
	router.GET("/test", testHandler)
	
	ginRouter := router.(*ginRouter)
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	ginRouter.root.ServeHTTP(w, req)
	
	assert.True(t, handlerCalled)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "GET")
}

func TestGinRouter_POST(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := NewRouter()
	
	handlerCalled := false
	testHandler := func(ctx http.Context) error {
		handlerCalled = true
		return ctx.JSON(200, map[string]string{"method": "POST"})
	}
	
	router.POST("/test", testHandler)
	
	ginRouter := router.(*ginRouter)
	req := httptest.NewRequest("POST", "/test", bytes.NewBuffer([]byte("{}")))
	w := httptest.NewRecorder()
	ginRouter.root.ServeHTTP(w, req)
	
	assert.True(t, handlerCalled)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "POST")
}

func TestGinRouter_PUT(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := NewRouter()
	
	handlerCalled := false
	testHandler := func(ctx http.Context) error {
		handlerCalled = true
		return ctx.JSON(200, map[string]string{"method": "PUT"})
	}
	
	router.PUT("/test", testHandler)
	
	ginRouter := router.(*ginRouter)
	req := httptest.NewRequest("PUT", "/test", bytes.NewBuffer([]byte("{}")))
	w := httptest.NewRecorder()
	ginRouter.root.ServeHTTP(w, req)
	
	assert.True(t, handlerCalled)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "PUT")
}

func TestGinRouter_DELETE(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := NewRouter()
	
	handlerCalled := false
	testHandler := func(ctx http.Context) error {
		handlerCalled = true
		return ctx.JSON(200, map[string]string{"method": "DELETE"})
	}
	
	router.DELETE("/test", testHandler)
	
	ginRouter := router.(*ginRouter)
	req := httptest.NewRequest("DELETE", "/test", nil)
	w := httptest.NewRecorder()
	ginRouter.root.ServeHTTP(w, req)
	
	assert.True(t, handlerCalled)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "DELETE")
}

func TestGinRouter_Use(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := NewRouter()
	
	middlewareCalled := false
	testMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
		return func(ctx http.Context) error {
			middlewareCalled = true
			ctx.Header("X-Middleware", "applied")
			return next(ctx)
		}
	}
	
	handlerCalled := false
	testHandler := func(ctx http.Context) error {
		handlerCalled = true
		return ctx.JSON(200, map[string]string{"message": "ok"})
	}
	
	router.Use(testMiddleware)
	router.GET("/test", testHandler)
	
	ginRouter := router.(*ginRouter)
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	ginRouter.root.ServeHTTP(w, req)
	
	assert.True(t, middlewareCalled)
	assert.True(t, handlerCalled)
	assert.Equal(t, "applied", w.Header().Get("X-Middleware"))
}

func TestGinRouter_Group(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := NewRouter()
	
	group := router.Group("/api/v1")
	assert.NotNil(t, group)
	
	ginGroup, ok := group.(*ginRouter)
	assert.True(t, ok)
	assert.NotNil(t, ginGroup.group)
	assert.NotNil(t, ginGroup.root)
	
	handlerCalled := false
	testHandler := func(ctx http.Context) error {
		handlerCalled = true
		return ctx.JSON(200, map[string]string{"group": "test"})
	}
	
	group.GET("/test", testHandler)
	
	ginRouter := router.(*ginRouter)
	req := httptest.NewRequest("GET", "/api/v1/test", nil)
	w := httptest.NewRecorder()
	ginRouter.root.ServeHTTP(w, req)
	
	assert.True(t, handlerCalled)
	assert.Equal(t, 200, w.Code)
}

func TestGinRouter_Group_Nested(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := NewRouter()
	
	apiGroup := router.Group("/api")
	v1Group := apiGroup.Group("/v1")
	
	assert.NotNil(t, v1Group)
	
	handlerCalled := false
	testHandler := func(ctx http.Context) error {
		handlerCalled = true
		return ctx.JSON(200, map[string]string{"nested": "group"})
	}
	
	v1Group.GET("/test", testHandler)
	
	ginRouter := router.(*ginRouter)
	req := httptest.NewRequest("GET", "/api/v1/test", nil)
	w := httptest.NewRecorder()
	ginRouter.root.ServeHTTP(w, req)
	
	assert.True(t, handlerCalled)
	assert.Equal(t, 200, w.Code)
}

func TestGinRouter_Run(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := NewRouter()
	
	ginRouter, ok := router.(*ginRouter)
	assert.True(t, ok)
	
	assert.NotNil(t, ginRouter.root)
	
	go func() {
		err := router.Run(":0")
		assert.NoError(t, err)
	}()
}

func TestGinRouter_Run_NoRoot(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := &ginRouter{root: nil}
	
	err := router.Run(":8080")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no root engine to run")
}

func TestGinRouter_MultipleMiddlewares(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := NewRouter()
	
	middleware1Called := false
	middleware1 := func(next http.HandlerFunc) http.HandlerFunc {
		return func(ctx http.Context) error {
			middleware1Called = true
			ctx.Header("X-Middleware-1", "called")
			return next(ctx)
		}
	}
	
	middleware2Called := false
	middleware2 := func(next http.HandlerFunc) http.HandlerFunc {
		return func(ctx http.Context) error {
			middleware2Called = true
			ctx.Header("X-Middleware-2", "called")
			return next(ctx)
		}
	}
	
	handlerCalled := false
	testHandler := func(ctx http.Context) error {
		handlerCalled = true
		return ctx.JSON(200, map[string]string{"message": "ok"})
	}
	
	router.Use(middleware1, middleware2)
	router.GET("/test", testHandler)
	
	ginRouter := router.(*ginRouter)
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	ginRouter.root.ServeHTTP(w, req)
	
	assert.True(t, middleware1Called)
	assert.True(t, middleware2Called)
	assert.True(t, handlerCalled)
	assert.Equal(t, "called", w.Header().Get("X-Middleware-1"))
	assert.Equal(t, "called", w.Header().Get("X-Middleware-2"))
}