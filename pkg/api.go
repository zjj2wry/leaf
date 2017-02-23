package pkg

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"leaf/log"
	"sync/atomic"

	"net/http"
	"net/url"
	"strings"
	"sync"
)

var (
	ERR_NO_API    = errors.New("no api")
	ERR_NO_OBJECT = errors.New("no api object")
	ERR_NO_MATCH  = errors.New("no prefix match")
)

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
			return ERR_NO_OBJECT
		}

		var line string
		//read a line util not null
		for {
			if !sc.Scan() {
				return ERR_NO_API
			}
			line = strings.TrimSpace(sc.Text())
			log.Infof("read line:%s", line)
			if len(line) != 0 {
				break
			}
		}
		//format check
		tokens := strings.SplitN(line, " ", 2)
		if len(tokens) < 2 {
			return fmt.Errorf("bad line: %s", line)
		}

		if !isHTTPMethod(tokens[0]) {
			return fmt.Errorf("bad method: %s", tokens[0])
		}
		api.Method = tokens[0]

		if _, err = url.ParseRequestURI(tokens[1]); err != nil {
			log.Fatalf("bad Url: %s", tokens[1])
		}
		api.URL = tokens[1]

		//get new line ,if null or prefix is httpmethod return
		line = strings.TrimSpace(sc.Peek())

		//null or http method return
		if line == "" || isHTTPMethod(line) {
			return nil
		}
		for sc.Scan() {
			if line = strings.TrimSpace(sc.Text()); line == "" {
				break
			}
			log.Infof("read line:%s", line)

			//read ret,if prefix is ret,break
			if value, err := getRealValue(line, "ret"); err == nil {
				api.ret = value[0]
				break
			}

			//read api body
			if value, err := getRealValue(line, "body"); err == nil {
				if strings.HasPrefix(value[0], "@") {
					if api.Body, err = ioutil.ReadFile(value[0][1:]); err != nil {
						log.Fatal(err)
					}
				} else {
					api.Body = []byte(value[0])
				}
			}

			//read api header; now just suport key1 = value1,value2 key2 = value2
			if value, err := getRealValue(line, "header"); err == nil {
				api.Header = make(map[string][]string)

				for i := range value {
					if value[i] = strings.TrimSpace(value[i]); value[i] == "" {
						return fmt.Errorf("bad header: %s", line)
					}
					if !strings.ContainsAny(value[i], "=") {
						log.Fatal("header must use =")
					}
					kv := strings.Split(value[i], "=")
					sl := strings.Split(kv[1], ",")

					api.Header[kv[0]] = sl
				}
			}

			if err = sc.Err(); err != nil {
				return ERR_NO_API
			}
		}
		return nil
	}
}

func NewApiList(src io.Reader) (ApiGeter, error) {
	var (
		ap   = NewApi(src)
		apis []Api
		api  Api
	)
	for {
		if err := ap(&api); err == ERR_NO_API {
			break
		} else if err != nil {
			log.Fatal(err)
			return nil, err
		}
		apis = append(apis, api)
	}
	if len(apis) == 0 {
		return nil, ERR_NO_API
	}
	return getApiMultiple("sequence", apis...), nil
}

//now just suport sequence
//TODO:
//get api random
//get api weight
func getApiMultiple(mode string, apis ...Api) ApiGeter {
	i := int64(-1)

	switch mode {
	case "random":
		return nil
	case "weight":
		return nil
	default:
		return func(api *Api) error {
			*api = apis[atomic.AddInt64(&i, 1)%int64(len(apis))]
			return nil
		}
	}
}

//check first word and return other value
func getRealValue(line, pre string) ([]string, error) {
	if !strings.HasPrefix(line, pre) {
		return nil, ERR_NO_MATCH
	}

	tokens := strings.Split(line, " ")
	if len(tokens) < 2 {
		return nil, fmt.Errorf("bad format: %s", line)
	}

	return tokens[1:], nil
}

//check method
func isHTTPMethod(method string) bool {
	var array2 = [...]string{"GET", "POST", "PUT", "DELETE"}
	for _, value := range array2 {
		if strings.ToUpper(method) == value {
			return true
		}
	}
	return false
}
