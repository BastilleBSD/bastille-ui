package web

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func callBastilleAPI(path string, params map[string]string) (string, error) {

	node := getActiveNode()
	if node == nil {
		return "", fmt.Errorf("no node selected")
	}

	scheme := "http"
	if node.Port == "443" {
		scheme = "https"
	}
	
	rawurl := fmt.Sprintf("%s://%s:%s%s", scheme, node.Host, node.Port, path)
	u, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}

	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	req, _ := http.NewRequest("POST", u.String(), nil)
	req.Header.Set("Authorization-ID", node.KeyID)
	req.Header.Set("Authorization", "Bearer "+node.Key)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return string(body), fmt.Errorf("API error: %s", resp.Status)
	}
	return string(body), nil
}

func callBastilleAPILive(path string, params map[string]string) (string, error) {

	node := getActiveNode()
	if node == nil {
		return "", fmt.Errorf("no node selected")
	}

	scheme := "http"
	if node.Port == "443" {
		scheme = "https"
	}
	
	rawurl := fmt.Sprintf("%s://%s:%s%s", scheme, node.Host, node.Port, path)
	u, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}

	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization-ID", node.KeyID)
	req.Header.Set("Authorization", "Bearer "+node.Key)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("API error: %s", resp.Status)
	}

	// Get the port from a header instead of body
	ttydpath := resp.Header.Get("X-TTYD-Url")
	if ttydpath == "" {
		return "", fmt.Errorf("API did not return ttyd info")
	}

	return fmt.Sprintf("%s://%s:%s%s", scheme, node.Host, node.Port, ttydpath), nil
}
