package httpclient;

import (
	"time"
	"log"
	"fmt"
	"github.com/valyala/fasthttp"
);

func sendReq(url, method string,
			 body []byte, 
			 header *fasthttp.RequestHeader, 
			 timeout time.Duration) ([]byte, string, int, error) {
	request := fasthttp.AcquireRequest();
	defer fasthttp.ReleaseRequest(request);

	response := fasthttp.AcquireResponse();
	defer fasthttp.ReleaseResponse(response);

	request.SetRequestURI(url);
	request.Header.SetMethod(method);
	request.SetBody(body);

	log.Printf("httpclient :: making a %s request to %s", method, url);

	err := fasthttp.DoTimeout(request, response, timeout);

	log.Printf("request: %+v\n", request);
	log.Printf("response: %+v\n", response);

	if(err != nil) {
		log.Printf("httpclient :: error during %s request to %s :: %v", method, url, err);
		//TODO(Dylan): instead of just returning 500, determine a reasonable status to return based of the type of error.
        return []byte(""), "", 500, err;	
	}

    var bodyBytes []byte = []byte("");

	bodyBytes = response.Body();
	contentType := string(response.Header.ContentType());
	statusCode := response.Header.StatusCode();

	return bodyBytes, contentType, statusCode, err;
}

func GetBytes(url string, timeoutMS int) ([]byte, string, int, error) {
	durationString := fmt.Sprintf("%dms",timeoutMS);
	timeout, parseErr := time.ParseDuration(durationString);
	if(parseErr != nil) {
		return []byte(""), "", 500, parseErr;
	}
	header := &fasthttp.RequestHeader{};
    bodyBytes, contentType, statusCode, err := sendReq(url, "GET", []byte(""), header, timeout);
    return bodyBytes, contentType, statusCode, err;
}