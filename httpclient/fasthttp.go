package httpclient

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"liuyi-test/utils"
)

type PostKind string

var currentPostKind PostKind = PostKindFormUrlencoded

func (p *PostKind) ParamHandler() func(params map[string]interface{}) (string,error){
	return func(params map[string]interface{}) (s string, err error) {

		if *p == PostKindFormUrlencoded {
			var p string
			if params==nil || len(params)==0{
				return p,utils.NewError1("empty params")
			}

			for k,v:=range params{
				if p == ""{
					p = k+"="+fmt.Sprintf("%v",v)
				}else{
					p = "&"+k+"="+fmt.Sprintf("%v",v)
				}
			}

			return p,nil
		}else{
			return "",utils.NewError1("not implement yet")
		}
	}
}

func (p *PostKind) JsonHandler() func(json string)(string,error){
	return func(json string) (string, error) {
		return "",utils.NewError1("not implement yet")
	}
}

const(
	PostKindFormUrlencoded PostKind = "application/x-www-form-urlencoded"
	Post_Kind_From_Json    PostKind = "application/json"
)

func paramForGet(url string, params map[string]interface{}) string {
	if params != nil && len(params) > 0 {
		var p string

		for k, v := range params {
			if utils.StringIsBlank(p) {
				p = "?" + k + "=" + fmt.Sprintf("%v", v)
			} else {
				p = p + "&" + k + "=" + fmt.Sprintf("%v", v)
			}
		}

		return url + p
	} else {
		return url
	}
}

func header(request *fasthttp.Request, headers map[string]string) {

	if headers == nil || len(headers) == 0 {
		return
	}

	for k, v := range headers {
		request.Header.Add(k, v)
	}
}

func Get(url string) (int, string, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	req.SetRequestURI(url)
	req.Header.SetMethod("GET")
	return execute(req, resp)
}

func GetWithParam(params map[string]interface{}, url string) (int, string, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(paramForGet(url, params))
	req.Header.SetMethod("GET")
	return execute(req, resp)
}

func GetWithParamHeader(params map[string]interface{}, headers map[string]string, url string) (int, string, error) {

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(paramForGet(url, params))
	req.Header.SetMethod("GET")
	header(req, headers)

	return execute(req, resp)
}

func PostWithParam(params map[string]interface{}, url string) (int, string, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(url)
	line,e:=currentPostKind.ParamHandler()(params)
	if e!=nil {
		return 0,"",e
	}
	req.SetBody([]byte(line))
	req.Header.SetMethod("POST")

	return execute(req, resp)
}

func PostWithParamHeader(params map[string]interface{}, headers map[string]string, url string) (int, string, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(url)

	line,e:=currentPostKind.ParamHandler()(params)
	if e!=nil {
		return 0,"",e
	}
	req.SetBody([]byte(line))
	req.Header.SetMethod("POST")
	header(req, headers)

	return execute(req, resp)
}


func execute(req *fasthttp.Request, resp *fasthttp.Response) (int, string, error) {
	err := fasthttp.Do(req, resp)
	if err != nil {
		return 0, "", err
	}
	return resp.StatusCode(), string(resp.Body()), nil
}
