package httpfuzz

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
)

// Client is a modified net/http Client that can natively handle our request and response types
type Client struct {
	*http.Client
}

// Do wraps Go's net/http client with our Request and Response types.
func (c *Client) Do(req *Request) (*Response, error) {
	resp, err := c.Client.Do(req.Request)
	return &Response{Response: resp}, err
}

// Request is a *http.Request that allows cloning its body.
type Request struct {
	*http.Request
}

// CloneBody makes a copy of a request, including its body, while leaving the original body intact.
func (r *Request) CloneBody(ctx context.Context) (*Request, error) {
	req := &Request{Request: r.Request.Clone(ctx)}
	if req.Body == nil {
		return req, nil
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	// Put back the original body
	r.Request.Body = ioutil.NopCloser(bytes.NewReader(body))

	// Clone the request body
	req.Request.Body = ioutil.NopCloser(bytes.NewReader(body))
	return req, nil
}

// Response is a *http.Response that allows cloning its body.
type Response struct {
	*http.Response
}

// CloneBody makes a copy of a response, including its body, while leaving the original body intact.
func (r *Response) CloneBody() (*Response, error) {
	newResponse := new(http.Response)
	newResponse.Header = r.Response.Header.Clone()

	body, err := ioutil.ReadAll(r.Response.Body)
	if err != nil {
		return nil, err
	}

	// Put back the original body
	r.Response.Body = ioutil.NopCloser(bytes.NewReader(body))

	// Clone the request body
	newResponse.Body = ioutil.NopCloser(bytes.NewReader(body))
	newResponse.Trailer = r.Response.Trailer.Clone()
	newResponse.ContentLength = r.Response.ContentLength
	newResponse.Uncompressed = r.Response.Uncompressed
	newResponse.Request = r.Response.Request
	newResponse.TLS = r.Response.TLS
	newResponse.Status = r.Response.Status
	newResponse.StatusCode = r.Response.StatusCode
	newResponse.Proto = r.Response.Proto
	newResponse.ProtoMajor = r.Response.ProtoMajor
	newResponse.ProtoMinor = r.Response.ProtoMinor
	newResponse.Close = r.Response.Close
	copy(newResponse.TransferEncoding, r.Response.TransferEncoding)
	return &Response{Response: newResponse}, nil
}
