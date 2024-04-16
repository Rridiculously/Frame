package web

import (
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
	h := &HTTPServer{}

	h.AddRoute(http.MethodGet, "/user", func(ctx *Context) {

	})
	h.Get("/user", func(ctx *Context) {})
	//h.AddRoute1(http.MethodGet, "/user", func(ctx *Context) {
	//
	//}, func(ctx *Context) {
	//
	//})
	//用法一
	http.ListenAndServe(":8080", h)

	//用法二
	h.Start(":8080")
}
