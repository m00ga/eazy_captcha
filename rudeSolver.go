package eazycaptcha

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	timeout    = 10
	maxTimeout = 120
)

type rudeRequest struct {
	URL, postParams, getParams string
}

func rudeSolve(req *rudeRequest) (string, error) {
	var timeToExit int
	client := http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(req.URL + "/in.php?" + req.postParams)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	status, err := getResult(resp)
	if err.Error() != "ready" {
		return "", err
	}
	id := status[1]
	req.getParams += fmt.Sprintf("&id=%s", id)

	for {
		resp, err = client.Get(req.URL + "/res.php?" + req.getParams)
		if err != nil {
			return "", err
		}
		//defer respGet.Body.Close()

		status, err = getResult(resp)
		if err.Error() == "wait" {
			if timeToExit >= maxTimeout {
				return "", &reCapError{"reached timeout"}
			}
			time.Sleep(timeout * time.Second)
			timeToExit += timeout
			continue
		} else if err.Error() == "ready" {
			return status[1], nil
		} else {
			return "", err
		}
	}
}

func getResult(resp *http.Response) ([]string, *reCapError) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []string{}, &reCapError{err.Error()}
	}
	res := string(body)
	if res == "CAPCHA_NOT_READY" {
		return []string{}, &reCapError{"wait"}
	}
	status := strings.Split(res, "|")
	if len(status) != 2 {
		return []string{}, &reCapError{res}
	}
	return status, &reCapError{"ready"}
}
