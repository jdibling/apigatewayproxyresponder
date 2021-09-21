package apigatewayproxyresponder

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type Responder struct {
	headers http.Header
}

func defaultHeaders() http.Header {
	return http.Header{
		"Content-Type": 					{"application/json"},
		"Access-Control-Allow-Origin":      {"*"},
		"Access-Control-Allow-Credentials": {"true"},
		"Access-Control-Allow-Methods":     {"OPTIONS,GET,POST"},
	}
}

func New(headers http.Header) (Responder, error) {
	// set headers
	resp := Responder{headers: defaultHeaders()}
	for k, v := range headers {
		resp.headers.Set(k, v[0])
		for _, vv := range v[1:] {
			resp.headers.Add(k, vv)
		}
	}

	return resp, nil
}

func (r *Responder) Respond(code int, body interface{}) (events.APIGatewayProxyResponse, error) {
	retBody := ""

	switch v := body.(type) {
	case nil:
		retBody = ""
	case string:
		retBody = fmt.Sprintf("%s", v)
	case error:
		retBody = v.Error()
	default:
		bytes, err := json.Marshal(body)
		if err != nil {
			return events.APIGatewayProxyResponse{}, fmt.Errorf("marshaling body; %w", err)
		}
		retBody = string(bytes)
	}

	return events.APIGatewayProxyResponse{
		Body:       retBody,
		MultiValueHeaders: r.headers,
		StatusCode: code,
	}, nil
}

// MakeHeaders returns the headers for any response.
/*
Every succesful response must retrurn at least these headers:
- Content-Type 											(default: "application/json")
- Access-Control-Allow-Origin				(default: "*")
- Access-Control-Allow-Credentials	(default: "true")

For a normal success response, simply call xcr.RespHeaders():

return events.APIGatewayProxyResponse{
	Body:       string(bytes),
		Headers:    xcr.RespHeaders(),
		StatusCode: http.StatusCreated}, nil

To return plain text instead of json, override the Content-Type header:

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Hello from Stream. %s", time.Now().String()),
		Headers:    MakeHeaders(map[string]string{"Content-Type": "text"}),
		StatusCode: 200}, nil
*/
func MakeHeaders(overrides ...map[string]string) (ret map[string]string) {

	// DefaultHeaders are the defaults for successful http responses
	var DefaultRespHeaders map[string]string = map[string]string{
		"Content-Type":                     "application/json",
		"Access-Control-Allow-Origin":      "*",
		"Access-Control-Allow-Credentials": "true",
		"Access-Control-Allow-Methods":     "OPTIONS,GET,POST",
	}


	ret = DefaultRespHeaders
	for _, override := range overrides {
		for key, value := range override {
			ret[key] = value
		}
	}
	return ret
}


