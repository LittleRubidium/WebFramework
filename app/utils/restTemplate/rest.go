package restTemplate

import (
	"net/http"
	"net/url"
)

func PostForObject(URL string, data url.Values) error {
	client := http.DefaultClient
	_, err := client.PostForm(URL,data)
	if err != nil {
		return err
	}
	return nil
}