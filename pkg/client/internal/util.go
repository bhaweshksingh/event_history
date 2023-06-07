package internal

import (
	"net/url"
)

func BuildURL(baseURL, path string, queryParams map[string]string) (*url.URL, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	u.Path = path

	if len(queryParams) == 0 {
		return u, nil
	}

	q := u.Query()

	for key, value := range queryParams {
		q.Set(key, value)
	}

	u.RawQuery = q.Encode()
	return u, nil
}
