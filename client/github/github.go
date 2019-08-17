package clientGithub

import (
	"bytes"
	"io"
	"net/http"
)

const githubUrl = "https://api.github.com"
const defaultReference = "master"

func GetFileContent(repository string, path string, reference string, token string) (string, error) {
	headers := make(map[string]string)

	if reference == "" {
		reference = defaultReference
	}

	url := githubUrl + "/repos/" + repository + "/contents/" + path + "?ref=" + reference

	headers["Accept"] = "application/vnd.github.v3.raw"
	headers["Accept"] = "application/vnd.github.v3.raw"
	if token != "" {
		headers["Authorization"] = "token " + token
	}

	response, err := Request("GET", url, nil, headers)
	if err != nil {
		return "", err
	}

	buffer := new(bytes.Buffer)
	_, err = buffer.ReadFrom(response.Body)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func Request(method string, url string, body io.Reader, headers map[string]string) (*http.Response, error) {
	var response *http.Response
	client := http.Client{}

	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return response, err
	}

	for key, value := range headers {
		request.Header.Add(key, value)
	}

	if request.Header.Get("User-Agent") == "" {
		request.Header.Add("User-Agent", "go-client")
	}

	return client.Do(request)
}
