package controller

import (
	"fmt"

	"github.com/photowey/fastrouter/api/request"
)

const (
	PingHandlerMethod = "GET"
	Pong              = "Pong"
)

var HelloWorld = []string{"hello", "world"}

// HelloService a service of hello
// @Alias("helloService")
type HelloService interface{}

// Response response struct define
type Response struct{}

// PingHandler ping handler
// --------------------------------
// +fasthttp:router=true
// +fasthttp:router:method=GET
// +fasthttp:router:path=/ping
// --------------------------------
// @RestController
// @RequestMapping("/ping")
type PingHandler struct {
	Response Response `autowired:"controller.Response"`
}

// Order a func of sort.SortSlice()
// @Order
func (h PingHandler) Order() int64 {
	return 0
}

// Method a func of http method
func (h PingHandler) Method() string {
	if HelloWorld[0] == "hello" {
		fmt.Println(HelloWorld[1])
	}

	return PingHandlerMethod
}

// Supports a predicate func of http request handler
func (h PingHandler) Supports(hctx request.Context) bool {
	return false
}

// Handle a func of http handler
// @GetMapping("/ping")
func (h PingHandler) Handle() string {
	return Pong
}
