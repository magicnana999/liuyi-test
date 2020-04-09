package httpclient

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGet(t *testing.T) {
	status,resp,e:= Get("http://127.0.0.1:12345/liuyi/foo")
	fmt.Println(status,resp,e)
	assert.Nil(t,e,e)
	assert.True(t,status==200,status)
	assert.True(t,jsoniter.Get([]byte(resp),"code").ToInt() == 200,resp)
}

func TestGetWithParam(t *testing.T) {
	params:=make(map[string]interface{})
	url := `http://127.0.0.1:12345/liuyi/foo`
	status,resp,e := GetWithParam(params,url)
	fmt.Println(status,resp,e)
	assert.Nil(t,e,e)
	assert.True(t,status==200,status)
	assert.True(t,jsoniter.Get([]byte(resp),"code").ToInt() == 200,resp)
}

func TestGetWithParamHeader(t *testing.T) {
	params:=make(map[string]interface{})
	headers:=make(map[string]string)
	url := `http://127.0.0.1:12345/liuyi/foo`
	status,resp,e := GetWithParamHeader(params,headers,url)
	fmt.Println(status,resp,e)
	assert.Nil(t,e,e)
	assert.True(t,status==200,status)
	assert.True(t,jsoniter.Get([]byte(resp),"code").ToInt() == 200,resp)
}

func TestPostWthParam(t *testing.T) {
	params:=make(map[string]interface{})
	params["phone"] = "110"
	headers:=make(map[string]string)
	headers["X-Liuyi-App-Key"] = "88937d1372b277f8080a3accb9d6174b"
	headers["Authorization"]="bearer xkskdlfrjjslwrjjfgjkdsfkjwer"
	url := `http://127.0.0.1:12345/liuyi/foo`
	status,resp,e := PostWithParam(params,url)
	fmt.Println(status,resp,e)
	assert.Nil(t,e,e)
	assert.True(t,status==200,status)
	assert.True(t,jsoniter.Get([]byte(resp),"code").ToInt() == 200,resp)
}

func TestPostWithParamHeader(t *testing.T) {
	params:=make(map[string]interface{})
	params["phone"] = "110"
	headers:=make(map[string]string)
	headers["X-Liuyi-App-Key"] = "88937d1372b277f8080a3accb9d6174b"
	headers["Authorization"]="bearer xkskdlfrjjslwrjjfgjkdsfkjwer"
	url := `http://127.0.0.1:12345/liuyi/foo`
	status,resp,e := PostWithParamHeader(params,headers,url)
	fmt.Println(status,resp,e)
	assert.Nil(t,e,e)
	assert.True(t,status==200,status)
	assert.True(t,jsoniter.Get([]byte(resp),"code").ToInt() == 200,resp)
}
