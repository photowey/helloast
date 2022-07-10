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
	ResponseFiled    Response         `autowired:"controller.Response" json:"responseFiled"`
	ResponseFiledPtr *Response        `autowired:"controller.Response" json:"ResponseFiledPtr"`
	Ctx              request.Context  `autowired:"request.Context" json:"ctx"`
	hctx             *request.Context `autowired:"request.Context" json:"hctx"`
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

// Supports2 a predicate func of http request handler
func (h PingHandler) Supports2(hctx *request.Context) bool {
	return false
}

// Supports3 a predicate func of http request handler
func (h PingHandler) Supports3(resp Response) bool {
	return false
}

// Supports4 a predicate func of http request handler
func (h PingHandler) Supports4(resp *Response) bool {
	return false
}

// Supports5 a predicate func of http request handler
func (h PingHandler) Supports5(
	resp Response,
	hctx *request.Context,
) bool {
	return false
}

// Supports6 a predicate func of http request handler
func (h PingHandler) Supports6(resp *Response) (*Response, bool) {
	return resp, false
}

// Supports7 a predicate func of http request handler
func (h PingHandler) Supports7(resp *Response, request *request.Context) (*request.Context, bool) {
	return request, false
}

// Handle a func of http handler
// @GetMapping("/ping")
func (h PingHandler) Handle() string {
	return Pong
}

// Run run...
func (h PingHandler) Run() (time string, err error) {
	return Pong, nil
}
