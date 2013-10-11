package octokit

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Client struct {
	httpClient *http.Client
	BaseURL    string
	Login      string
	Password   string
	Token      string
}

func (c *Client) WithLogin(login, password string) *Client {
	c.Login = login
	c.Password = password
	return c
}

func (c *Client) WithToken(token string) *Client {
	c.Token = token
	return c
}

func (c *Client) Requester(url *url.URL) *Requester {
	return &Requester{client: c, URL: url}
}

func (c *Client) Get(url *url.URL, headers Headers) (resp *Response, err error) {
	resp, err = c.Request("GET", url, headers, nil)
	return
}

func (c *Client) Patch(url *url.URL, headers Headers, params interface{}) (resp *Response, err error) {
	reader, err := jsonMarshalToReader(params)
	if err != nil {
		return
	}

	resp, err = c.Request("PATCH", url, headers, reader)
	return
}

func (c *Client) Request(method string, url *url.URL, headers Headers, content io.Reader) (resp *Response, err error) {
	if url == nil {
		url, _ = url.Parse(c.BaseURL)
	}

	request, err := http.NewRequest(method, url.String(), content)
	if err != nil {
		return
	}

	c.setDefaultHeaders(request)

	if headers != nil {
		for h, v := range headers {
			request.Header.Set(h, v)
		}
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return
	}

	err = checkResponse(response)
	if err != nil {
		return
	}

	resp = &Response{Response: response}
	return
}

// Deprecated
func (c *Client) get(path string, headers Headers) ([]byte, error) {
	return c.request("GET", path, headers, nil)
}

func (c *Client) post(path string, headers Headers, content io.Reader) ([]byte, error) {
	return c.request("POST", path, headers, content)
}

func (c *Client) jsonGet(path string, options *Options, v interface{}) error {
	var headers Headers
	if options != nil {
		headers = options.Headers
	}

	body, err := c.get(path, headers)
	if err != nil {
		return err
	}

	return jsonUnmarshal(body, v)
}

func (c *Client) jsonPost(path string, options *Options, v interface{}) error {
	var headers Headers
	if options != nil {
		headers = options.Headers
	}

	var buffer *bytes.Buffer
	if options != nil && options.Params != nil {
		b, err := jsonMarshal(options.Params)
		if err != nil {
			return err
		}

		buffer = bytes.NewBuffer(b)
	}

	// *bytes.Buffer(nil) != nil
	// see http://golang.org/doc/faq#nil_error
	var content io.Reader
	if buffer == nil {
		content = nil
	} else {
		content = buffer
	}

	body, err := c.post(path, headers, content)
	if err != nil {
		return err
	}

	return jsonUnmarshal(body, v)
}

func (c *Client) request(method, path string, headers Headers, content io.Reader) ([]byte, error) {
	url, err := c.buildURL(path)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(method, url.String(), content)
	if err != nil {
		return nil, err
	}

	c.setDefaultHeaders(request)

	if headers != nil {
		for h, v := range headers {
			request.Header.Set(h, v)
		}
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 400 && response.StatusCode < 600 {
		return nil, checkResponse(response)
	}

	return body, nil
}

// Deprecated

func (c *Client) buildURL(pathOrURL string) (*url.URL, error) {
	u, e := url.ParseRequestURI(pathOrURL)
	if e != nil {
		u, e = url.Parse(c.BaseURL)
		if e != nil {
			return nil, e
		}

		return u.Parse(pathOrURL)
	}

	return u, nil
}

func (c *Client) setDefaultHeaders(request *http.Request) {
	request.Header.Set("Accept", MediaType)
	request.Header.Set("User-Agent", UserAgent)
	request.Header.Set("Content-Type", DefaultContentType)
	if c.isBasicAuth() {
		request.Header.Set("Authorization", fmt.Sprintf("Basic %s", hashAuth(c.Login, c.Password)))
	} else if c.isTokenAuth() {
		request.Header.Set("Authorization", fmt.Sprintf("token %s", c.Token))
	}
}

func (c *Client) isBasicAuth() bool {
	return c.Login != "" && c.Password != ""
}

func (c *Client) isTokenAuth() bool {
	return c.Token != ""
}

func jsonMarshalToReader(v interface{}) (r io.Reader, err error) {
	if v != nil {
		b, e := jsonMarshal(v)
		if err != nil {
			err = e
			return
		}

		r = bytes.NewBuffer(b)
	}

	return
}
