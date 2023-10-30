# Nonce middleware for Gin

This Gin middleware generates a nonce for each request.

[![Build Status](https://github.com/joeig/gin-nonce/workflows/Tests/badge.svg)](https://github.com/joeig/gin-nonce/actions)
[![Test coverage](https://img.shields.io/badge/coverage-100%25-success)](https://github.com/joeig/gin-nonce/tree/main/.github/testcoverage.yml)
[![Go Report Card](https://goreportcard.com/badge/go.eigsys.de/gin-nonce)](https://goreportcard.com/report/go.eigsys.de/gin-nonce)
[![PkgGoDev](https://pkg.go.dev/badge/go.eigsys.de/gin-nonce)](https://pkg.go.dev/go.eigsys.de/gin-nonce)

## Usage

```go
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.eigsys.de/gin-nonce"
)

func main() {
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
```

## Documentation

See [Go reference](https://pkg.go.dev/go.eigsys.de/gin-nonce).
