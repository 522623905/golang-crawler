package worker

import (
	"errors"
	"fmt"
	"log"
	"u2pppw/crawler/crawler-distribute/config"
	"u2pppw/crawler/crawler-queue/engine"
	"u2pppw/crawler/crawler-queue/zhenai/parser"
)

//注意封装成可用rpc传递的结构
//主要是序列化和反序列化的操作

type SerializedParser struct {
	Name string      //函数名
	Args interface{} //参数
}

type Request struct {
	Url    string
	Parser SerializedParser
}

type ParseResult struct {
	Items    []engine.Item
	Requests []Request
}

//将engine.Request序列化成可供rpc调用的Request
func SerializeRequest(r engine.Request) Request {
	name, args := r.Parser.Serialize()
	return Request{
		Url: r.Url,
		Parser: SerializedParser{
			Name: name,
			Args: args,
		},
	}

}

//将engine.ParseResult序列化成可供rpc调用的ParseResult
func SerializeResult(
	r engine.ParseResult) ParseResult {
	result := ParseResult{
		Items: r.Items,
	}

	for _, req := range r.Requests {
		result.Requests = append(result.Requests, SerializeRequest(req))
	}

	return result
}

//将Request反序列化成engine.Request
func DeserializeRequest(r Request) (engine.Request, error) {
	parser, err := deserializeParser(r.Parser)
	if err != nil {
		return engine.Request{}, err
	}
	return engine.Request{
		Url:    r.Url,
		Parser: parser,
	}, nil
}

//将ParseResult反序列化成engine.ParseResult
func DeserializeResult(r ParseResult) engine.ParseResult {
	result := engine.ParseResult{
		Items: r.Items,
	}

	for _, req := range r.Requests {
		engineReq, err := DeserializeRequest(req)
		if err != nil {
			log.Printf("error deserializing request:%v", err)
			continue
		}
		result.Requests = append(result.Requests,
			engineReq)
	}
	return result
}

//将SerializedParser反序列化成engine.Parser接口
func deserializeParser(p SerializedParser) (engine.Parser, error) {
	switch p.Name {
	case config.ParseCityList:
		return engine.NewFuncParser(
			parser.ParseCityList,
			config.ParseCityList), nil

	case config.ParseCity:
		return engine.NewFuncParser(
			parser.ParseCity,
			config.ParseCity), nil

	case config.ParseProfile:
		if userName, ok := p.Args.(string); ok {
			return parser.NewProfileParser(
				userName), nil
		} else {
			return nil, fmt.Errorf("invalid arg:%v", p.Args)
		}

	case config.NilParser:
		return engine.NilParser{}, nil

	default:
		fmt.Println(p.Name)
		return nil, errors.New(
			"unknown parser name")
	}
}
