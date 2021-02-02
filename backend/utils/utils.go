package utils

import (
	"io/ioutil"
	"net/http"
	"time"
)

const USER_AGENT = "Mozilla/5.0 (Windows NT 10.0; rv:68.0) Gecko/20100101 Firefox/68.0"

func Request(URL string) (string, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// log.Println(fmt.Sprintf("Requesting %s", URL))
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", USER_AGENT)

	response, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", nil
	}

	return string(body), nil
}
