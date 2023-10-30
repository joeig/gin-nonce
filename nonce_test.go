package nonce_test

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	nonce "go.eigsys.de/gin-nonce"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewNonceHandler(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		resp := httptest.NewRecorder()
		ginContext, _ := gin.CreateTestContext(resp)
		handler := nonce.New().WithRandRead(func(b []byte) (int, error) {
			b[0] = 'b'
			b[1] = 'a'
			b[2] = 'r'

			return 3, nil
		})

		handler.Middleware()(ginContext)

		if handler.GetKey() != "go.eigsys.de/gin-nonce:nonce" {
			t.Error("wrong key")
		}
		if value, ok := ginContext.Get("go.eigsys.de/gin-nonce:nonce"); !ok || value.(string) != "YmFyAAAAAAAAAAAAAAAAAA" {
			t.Error("wrong nonce")
		}
	})

	t.Run("rand-read-error", func(t *testing.T) {
		resp := httptest.NewRecorder()
		ginContext, _ := gin.CreateTestContext(resp)
		handler := nonce.New().WithRandRead(func(b []byte) (n int, err error) {
			return 0, errors.New("oops")
		})

		handler.Middleware()(ginContext)

		if resp.Code != http.StatusInternalServerError {
			t.Error("wrong status code")
		}
	})
}

func TestHandler_WithKey(t *testing.T) {
	handler := nonce.New().WithKey("foo")

	if handler.GetKey() != "foo" {
		t.Error("wrong key")
	}
}

func TestHandler_GetKey(t *testing.T) {
	handler := nonce.New().WithKey("foo")

	if handler.GetKey() != "foo" {
		t.Error("wrong key")
	}
}

func TestHandler_GetNonce(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		resp := httptest.NewRecorder()
		ginContext, _ := gin.CreateTestContext(resp)
		ginContext.Set("go.eigsys.de/gin-nonce:nonce", "bar")
		handler := nonce.New()

		value, ok := handler.GetNonce(ginContext)

		if value != "bar" {
			t.Error("wrong value")
		}
		if !ok {
			t.Error("not ok")
		}
	})

	t.Run("error", func(t *testing.T) {
		resp := httptest.NewRecorder()
		ginContext, _ := gin.CreateTestContext(resp)
		handler := nonce.New()

		_, ok := handler.GetNonce(ginContext)

		if ok {
			t.Error("unexpectedly ok")
		}
	})
}

func ExampleNew() {
	router := gin.Default()

	handler := nonce.New()
	router.Use(handler.Middleware())

	router.GET("/", func(ginCtx *gin.Context) {
		currentNonce, _ := handler.GetNonce(ginCtx)
		ginCtx.Header("Content-Security-Policy", fmt.Sprintf("style-src 'nonce-%s';", currentNonce))

		ginCtx.String(http.StatusOK, "Hello, Gopher!")
	})

	_ = router.Run()
}
