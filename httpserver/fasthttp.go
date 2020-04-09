package httpserver

import (
	"github.com/buaazp/fasthttprouter"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"log"
	"strings"
)

type ApiResult struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

const (
	PHONE   string = "110"
	CODE    string = "12345"
	TOKEN   string = "xkskdlfrjjslwrjjfgjkdsfkjwer"
	APP_KEY string = "88937d1372b277f8080a3accb9d6174b"
)

var (
	strContentType     = []byte("Content-Type")
	strApplicationJSON = []byte("application/json")
)

func toJson(obj ApiResult) string {
	s, _ := jsoniter.MarshalToString(obj)
	return s
}

func intercept(ctx *fasthttp.RequestCtx)bool {
	appKey := string(ctx.Request.Header.Peek("X-Liuyi-App-Key"))

	if appKey == "" {
		res401(ctx, "appKey is not found")
		return false
	}

	if appKey!=APP_KEY{
		res401(ctx,"invalid appKey")
		return false
	}
	return true
}

func authorize(ctx *fasthttp.RequestCtx)bool {

	if !intercept(ctx){
		return false
	}

	authorization := string(ctx.Request.Header.Peek("Authorization"))
	if authorization == "" {
		res401(ctx, "unauthorized")
		return false
	}

	array := strings.Split(authorization, "bearer ")
	if array == nil || len(array) != 2 {
		res401(ctx, "unauthorized")
		return false
	}

	token := array[1]
	if token == "" {
		res401(ctx, "unauthorized")
		return false
	}

	if token != TOKEN {
		res401(ctx, "unauthorized")
		return false
	}

	return true

}

func res401(ctx *fasthttp.RequestCtx, msg string) {
	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
	ctx.Response.SetStatusCode(401)
	ctx.Response.SetBodyString(toJson(ApiResult{401, nil, msg}))
}

func res400(ctx *fasthttp.RequestCtx, msg string) {
	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
	ctx.Response.SetStatusCode(400)
	ctx.Response.SetBodyString(toJson(ApiResult{400, nil, msg}))
}
func res200(ctx *fasthttp.RequestCtx, obj interface{}) {
	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
	ctx.Response.SetStatusCode(200)
	ctx.Response.SetBodyString(toJson(ApiResult{200, obj, ""}))
}

func pushRegCode(ctx *fasthttp.RequestCtx) {
	if intercept(ctx){
		res200(ctx, nil)
	}
}

func showMeCode(ctx *fasthttp.RequestCtx) {
	if intercept(ctx){
		res200(ctx, CODE)
	}
}

func login(ctx *fasthttp.RequestCtx) {
	if intercept(ctx) {
		phone := string(ctx.QueryArgs().Peek("phone"))
		code := string(ctx.QueryArgs().Peek("code"))

		if phone == PHONE && code == CODE {
			res200(ctx, TOKEN)
		} else {
			res400(ctx, "invalid code")
		}
	}
}

func foo(ctx *fasthttp.RequestCtx) {
	if authorize(ctx) {
		result := make(map[string] string)
		result["name"]="张三丰"
		result["category"] = "武当"
		result["title"]="掌门"
		res200(ctx,result)
	}
}

func Server() {
	// 创建路由
	router := fasthttprouter.New()
	// 不同的路由执行不同的处理函数
	router.GET("/liuyi/push_code", pushRegCode)
	router.GET("/liuyi/show_code", showMeCode)
	router.GET("/liuyi/login", login)
	router.POST("/liuyi/foo", foo)

	// 启动web服务器，监听 0.0.0.0:12345
	log.Fatal(fasthttp.ListenAndServe(":12345", router.Handler))
}
