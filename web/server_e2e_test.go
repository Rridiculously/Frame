//go:build e2e

package web

import (
	"fmt"
	"net/http"
	"testing"
)

func TestServer_e2e(t *testing.T) {
	h := &HTTPServer{}

	h.addRoute(http.MethodGet, "/user", func(ctx Context) {
		fmt.Println("handler")
	})
	handler1 := func(ctx Context) {
		fmt.Println("handler1")
	}
	handler2 := func(ctx Context) {
		fmt.Println("handler2")
	}
	h.addRoute(http.MethodGet, "/user", func(ctx Context) {
		handler2(ctx)
		handler1(ctx)
	})
	//h.AddRoute1(http.MethodGet, "/user", func(ctx *Context) {
	//
	//}, func(ctx *Context) {
	//
	//})
	////用法一
	//http.ListenAndServe(":8080", h)

	//用法二
	h.Start(":8080")
}
