package pkg

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"errors"
)

var ERR_NO_API = errors.New("no api")

type Api struct {
	Method string
	URL    string
	Body   []byte
	Header http.Header
	ret    string
}

func (api *Api) Request() (*http.Request, error) {
	req, err := http.NewRequest(api.Method, api.URL, bytes.NewReader(api.Body))
	if err != nil {
		return nil, err
	}
	for k, vs := range api.Header {
		req.Header[k] = make([]string, len(vs))
		copy(req.Header[k], vs)
	}
	return req, nil
}

type ApiGeter func(*Api) error

func NewApi(src io.Reader) ApiGeter {
	var mu sync.Mutex
	sc := peekingScanner{src: bufio.NewScanner(src)}
	return func(api *Api) (err error) {
		mu.Lock()
		defer mu.Unlock()

		if api == nil {
			return errors.New("no api object")
		}

		var line string
		//read a line util it not null
		for {
			if !sc.Scan() {
				return ERR_NO_API
			}
			line = strings.TrimSpace(sc.Text())
			if len(line) != 0 {
				break
			}
		}

		//format check
		tokens := strings.SplitN(line, " ", 2)
		if len(tokens) < 2 {
			return fmt.Errorf("bad target: %s", line)
		}
		if !isHTTPMethod(line) {
			return fmt.Errorf("bad method: %s", tokens[0])
		}
		api.Method = tokens[0]
		if _, err = url.ParseRequestURI(tokens[1]); err != nil {
			return fmt.Errorf("bad URL: %s", tokens[1])
		}
		api.URL = tokens[1]
		line = strings.TrimSpace(sc.Peek())

		//null or http method return
		if line == "" || isHTTPMethod(line) {
			return nil
		}
		for sc.Scan() {
			if line = strings.TrimSpace(sc.Text()); line == "" {
				break
			}
			//read ret
			if value, err := getRealValue(line, "ret"); err == nil {
				api.ret = value[0]
				break
			}
			//api body
			if value, err := getRealValue(line, "ret"); err == nil {
				if api.Body, err = ioutil.ReadFile(value[0]); err != nil {
					return fmt.Errorf("bad body: %s", err)
				}
			} else if value, err := getRealValue(line, "json"); err == nil {
				api.Body = []byte(value[0])
			}
			//api header
			if value, err := getRealValue(line, "header"); err == nil {
				for i := range value {
					if value[i] = strings.TrimSpace(value[i]); value[i] == "" {
						return fmt.Errorf("bad header: %s", line)
					}
					kv := strings.Split(value[i], ":")

					api.Header.Set(kv[0], kv[1])
				}
			}

			if err = sc.Err(); err != nil {
				return ERR_NO_API
			}
		}
		return nil
	}
}

func NewApiList(src io.Reader) ([]Api, error) {
	var (
		ap   = NewApi(src)
		apis []Api
		api  Api
	)
	for {
		if err := ap(&api); err == ERR_NO_API {
			break
		} else if err != nil {
			return nil, err
		}
		apis = append(apis, api)
	}
	if len(apis) == 0 {
		return nil, ERR_NO_API
	}
	return apis, nil
}

//check first word and return other value
func getRealValue(line, pre string) ([]string, error) {
	if !strings.HasPrefix(line, pre) {
		return nil, errors.New("no head match")
	}
	tokens := strings.Split(line, " ")
	if len(tokens) < 2 {
		return nil, fmt.Errorf("bad format: %s", line)
	}
	return tokens[1:], nil
}

func isHTTPMethod(method string) bool {
	return true
}

//just use peek instead of Scan()+Text()
type peekingScanner struct {
	src    *bufio.Scanner
	peeked string
}

func (s *peekingScanner) Err() error {
	return s.src.Err()
}

func (s *peekingScanner) Peek() string {
	if !s.src.Scan() {
		return ""
	}
	s.peeked = s.src.Text()
	return s.peeked
}

func (s *peekingScanner) Scan() bool {
	if s.peeked == "" {
		return s.src.Scan()
	}
	return true
}

func (s *peekingScanner) Text() string {
	if s.peeked == "" {
		return s.src.Text()
	}
	t := s.peeked
	s.peeked = ""
	return t
}
