package core

import (
	"liuyi-test/expression"
	"liuyi-test/httpclient"
	"liuyi-test/id"
	"liuyi-test/logger"
	"liuyi-test/utils"
)

type Pipeline struct {
	TraceContext *TraceContext
	head *Node
	tail *Node
}

func NewPipeline(traceConfig *TraceConfig,config *Config) *Pipeline {

	if traceConfig == nil {
		utils.PanicError(utils.NewError1(traceConfigNil))
	}

	if traceConfig.Id == "" {
		traceConfig.Id = id.StrID()
	}

	ctx := &TraceContext{
		traceConfig: 	*traceConfig,
		data:        		make(map[string]interface{}),
		http: 				config.Http,
		socket: 			config.Socket,
	}


	p := &Pipeline{
		TraceContext: ctx,
		head:         		nil,
		tail: 					nil,
	}

	p.addNodeAll(traceConfig,ctx)
	return p
}





func (p *Pipeline) addLast(node *Node){

	if p.tail == nil{
		p.tail = node
		p.head = node
	}else{
		p.tail.next = node
		node.pre = p.tail
		p.tail = node
	}
}

func getHeaderOrParamFromCtx(ctx map[string]interface{},script string)string{
	v,e:=expression.Param(ctx,script)
	utils.PanicError(e)
	return v.(string)
}

func(p *Pipeline) addNodeAll(traceConfig *TraceConfig,traceContext * TraceContext){
	for _,spanConfig := range traceConfig.Spans{
		if spanConfig.Id == "" {
			spanConfig.Id = id.StrID()
		}
		node := &Node{
			next:    nil,
			pre:     nil,
			name:	spanConfig.Name,
			spanContext: &SpanContext{
				spanConfig: spanConfig,
				params:     make(map[string]interface{}),
				headers:    make(map[string]string),
				url:        "",
				result:     "",
				log: logger.NewLogger(traceConfig.Id,spanConfig.Id),
			},
			traceContext: traceContext,
			before: func(traceContext *TraceContext, spanContext *SpanContext) {
				if spanContext.spanConfig.Kind == HTTP_KIND{

					if traceContext.http.Headers != nil && len(traceContext.http.Headers)>0{
						for _,head := range traceContext.http.Headers{
							spanContext.headers[head.Name] = getHeaderOrParamFromCtx(traceContext.data,head.Value)
						}
					}

					if spanContext.spanConfig.HttpEntry.Headers != nil && len(spanContext.spanConfig.HttpEntry.Headers)>0{
						for _,head := range spanContext.spanConfig.HttpEntry.Headers{
							spanContext.headers[head.Name] = getHeaderOrParamFromCtx(traceContext.data,head.Value)
						}
					}

					if traceContext.http.Domain != "" {
						spanContext.url = traceContext.http.Domain+spanContext.spanConfig.HttpEntry.Url
					}

					if spanContext.spanConfig.HttpEntry.Params != nil && len(spanContext.spanConfig.HttpEntry.Params)>0{
						for _,p := range spanContext.spanConfig.HttpEntry.Params{
							spanContext.params[p.Name] = getHeaderOrParamFromCtx(traceContext.data,p.Value)
						}
					}
				}
			},
			invoker: func(traceContext *TraceContext, spanContext *SpanContext) {
				if spanContext.params != nil && len(spanContext.params)>0{
					if spanContext.headers!=nil && len(spanContext.headers)>0{
						if spanContext.spanConfig.HttpEntry.Method == GET{
							c,r,e := httpclient.GetWithParamHeader(spanContext.params,spanContext.headers,spanContext.url)
							utils.PanicError(e)
							spanContext.status = c
							spanContext.result = r
						}else{
							c,r,e := httpclient.PostWithParamHeader(spanContext.params,spanContext.headers,spanContext.url)
							utils.PanicError(e)
							spanContext.status = c
							spanContext.result = r
						}
					}else{
						if spanContext.spanConfig.HttpEntry.Method == GET {
							c,r, e := httpclient.GetWithParam(spanContext.params, spanContext.url)
							utils.PanicError(e)
							spanContext.status = c
							spanContext.result = r
						}else{
							c,r, e := httpclient.PostWithParam(spanContext.params, spanContext.url)
							utils.PanicError(e)
							spanContext.status = c
							spanContext.result = r
						}
					}
				}else{
					if spanContext.spanConfig.HttpEntry.Method == GET {
						c,r,e:= httpclient.Get(spanContext.url)
						utils.PanicError(e)
						spanContext.status = c
						spanContext.result = r
					}else{
						utils.PanicError(utils.NewError1("Post must include parameters"))
					}
				}

				fields := make(map[string] interface{})
				fields["headers"] = spanContext.headers
				fields["params"] = spanContext.params
				fields["url"] = spanContext.url
				fields["resp"] = spanContext.result
				spanContext.log.WithFields(fields).Info("success")
			},
			after: func(traceContext *TraceContext, spanContext *SpanContext) {

				f1 := spanContext.status == 200
				f2,e:=expression.Assert(traceContext.data,spanContext.result,spanContext.spanConfig.HttpEntry.State)
				utils.PanicError(e)

				if f1 && f2{
					if spanContext.spanConfig.HttpEntry.Success != nil && len(spanContext.spanConfig.HttpEntry.Success)>0 {
						for _,v:= range spanContext.spanConfig.HttpEntry.Success{
							expression.Assign(traceContext.data,spanContext.result,v)
						}
					}
				}else{
					if spanContext.spanConfig.HttpEntry.Fail != "" && BREAK==spanContext.spanConfig.HttpEntry.Fail{
						utils.PanicError(utils.NewError2(breakIfInvokeFail,spanContext.spanConfig.Name))
					}
				}
			},
		}

		p.addLast(node)
	}
}

func (p *Pipeline) Start(){
	if p.head!=nil {
		p.head.invoke()
	}
}



type Node struct {
	next *Node
	pre *Node
	name string
	spanContext *SpanContext
	traceContext * TraceContext

	before func(traceContext *TraceContext,spanContext *SpanContext)
	invoker func(traceContext *TraceContext,spanContext *SpanContext)
	after func(traceContext *TraceContext,spanContext *SpanContext)

}

func(node *Node) invoke(){
	node.before(node.traceContext,node.spanContext)
	node.invoker(node.traceContext,node.spanContext)
	node.after(node.traceContext,node.spanContext)

	if node.next != nil{
		node.next.invoke()
	}
}

const (
	breakIfInvokeFail string = "break, [%s] failed"
	traceConfigNil string = "trace config is nil"
)





