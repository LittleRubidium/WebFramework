package restTemplate

import (
	"net/http"
	"net/url"
)

func PostForObject(URL string, data url.Values) (err error) {
	client := http.DefaultClient
	_, err = client.PostForm(URL, data)
	return
}

func DeleteForObject(URL string) (err error) {
	req, _ := http.NewRequest(http.MethodDelete, URL, nil)
	_, err = http.DefaultClient.Do(req)
	return
}