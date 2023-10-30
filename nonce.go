// Package nonce is a Gin middleware which generates a nonce for each request.
package nonce

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"net/http"
)

const defaultKey = "go.eigsys.de/gin-nonce:nonce"

type randReadFunc func(b []byte) (n int, err error)

// Handler defines a nonce handler.
type Handler struct {
	key      string
	randRead randReadFunc
}

// New creates a new Handler instance.
func New() *Handler {
	return &Handler{
		key:      defaultKey,
		randRead: rand.Read,
	}
}

// WithKey overrides the default Gin context key.
func (h *Handler) WithKey(key string) *Handler {
	h.key = key
	return h
}

// WithRandRead overrides the default random reader ([rand.Read]).
func (h *Handler) WithRandRead(randRead randReadFunc) *Handler {
	h.randRead = randRead
	return h
}

// Middleware creates a new Gin middleware which generates a nonce for each request.
func (h *Handler) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		rawID := make([]byte, 16)
		if _, err := h.randRead(rawID); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		c.Set(h.key, base64.RawStdEncoding.EncodeToString(rawID))
	}
}

// GetKey returns the Gin context key.
func (h *Handler) GetKey() string {
	return h.key
}

// GetNonce returns the nonce for the current Gin context.
func (h *Handler) GetNonce(c *gin.Context) (string, bool) {
	nonce, ok := c.Get(h.key)
	if nonce == nil {
		return "", false
	}

	return nonce.(string), ok
}
