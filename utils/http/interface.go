package http

import (
	"io/ioutil"
	"net/http"
	"net/url"
)

type AuthClient struct {
	client *http.Client
}

func NewAuthClient() *AuthClient {
	cli := &AuthClient{
		client: &http.Client{},
	}

	jar := &myjar{}
	jar.jar = make(map[string][]*http.Cookie)
	cli.client.Jar = jar

	return cli
}

func (this *AuthClient) Auth(urlStr string, values map[string]string) error {
	form := make(url.Values)
	for k, v := range values {
		form.Set(k, v)
	}
	resp, err := this.client.PostForm(urlStr, form)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (this *AuthClient) Get(urlStr string) ([]byte, error) {
	resp, err := this.client.Get(urlStr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

type myjar struct {
	jar map[string][]*http.Cookie
}

func (this *myjar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	this.jar[u.Host] = cookies
}

func (this *myjar) Cookies(u *url.URL) []*http.Cookie {
	return this.jar[u.Host]
}
