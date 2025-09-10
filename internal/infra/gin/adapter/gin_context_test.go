package adapter

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGinContext_Param(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	router.GET("/test/:id", func(c *gin.Context) {
		ginCtx := &GinContext{ctx: c}
		result := ginCtx.Param("id")
		c.String(200, result)
	})
	
	req := httptest.NewRequest("GET", "/test/123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, "123", w.Body.String())
}

func TestGinContext_Query(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	router.GET("/test", func(c *gin.Context) {
		ginCtx := &GinContext{ctx: c}
		result := ginCtx.Query("name")
		c.String(200, result)
	})
	
	req := httptest.NewRequest("GET", "/test?name=john", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, "john", w.Body.String())
}

func TestGinContext_Bind(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	type TestData struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	
	router.POST("/test", func(c *gin.Context) {
		ginCtx := &GinContext{ctx: c}
		var data TestData
		err := ginCtx.Bind(&data)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, data)
	})
	
	testData := TestData{Name: "John", Age: 30}
	jsonData, _ := json.Marshal(testData)
	
	req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, 200, w.Code)
	
	var response TestData
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "John", response.Name)
	assert.Equal(t, 30, response.Age)
}

func TestGinContext_JSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	router.GET("/test", func(c *gin.Context) {
		ginCtx := &GinContext{ctx: c}
		data := map[string]string{"message": "hello"}
		err := ginCtx.JSON(200, data)
		assert.NoError(t, err)
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Body.String(), "hello")
}

func TestGinContext_Status(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	router.GET("/test", func(c *gin.Context) {
		ginCtx := &GinContext{ctx: c}
		ginCtx.Status(201)
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, 201, w.Code)
}

func TestGinContext_Header(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	router.GET("/test", func(c *gin.Context) {
		ginCtx := &GinContext{ctx: c}
		ginCtx.Header("X-Custom-Header", "custom-value")
		ginCtx.Status(200)
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, "custom-value", w.Header().Get("X-Custom-Header"))
}

func TestWrapHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	handlerCalled := false
	testHandler := func(ctx Context) error {
		handlerCalled = true
		param := ctx.Param("id")
		ctx.JSON(200, map[string]string{"id": param})
		return nil
	}
	
	router := gin.New()
	router.GET("/test/:id", wrapHandler(testHandler))
	
	req := httptest.NewRequest("GET", "/test/123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.True(t, handlerCalled)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "123")
}

func TestWrapMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	middlewareCalled := false
	testMiddleware := func(next HandlerFunc) HandlerFunc {
		return func(ctx Context) error {
			middlewareCalled = true
			ctx.Header("X-Middleware", "called")
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
	
	assert.True(t, middlewareCalled)
	assert.Equal(t, "called", w.Header().Get("X-Middleware"))
}