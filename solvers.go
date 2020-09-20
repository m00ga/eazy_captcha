package eazycaptcha

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

// *** Rude Solving block BEGIN ***

func (rs rudeSolver) solve() (string, error) {
	var response string
	client := http.Client{Timeout: 5 * time.Second}
	req, err := json.Marshal(&rs.request.postParams)
	if err != nil {
		return "", err
	}

	resp, err := client.Post(rs.request.url+"/in.php", "application/json", bytes.NewBuffer(req))

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var rcresp map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&rcresp)

	if rcresp["status"] != 1.0 {
		return "", errors.New(rcresp["request"].(string))
	}
	rs.request.getParams["id"] = rcresp["request"]

	get, err := json.Marshal(&rs.request.getParams)
	if err != nil {
		return "", err
	}

	for {
		resp, err := client.Post(rs.request.url+"/res.php", "application/json", bytes.NewBuffer(get))
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		var rcresp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&rcresp)

		if status, request := rcresp["status"], rcresp["request"].(string); status != 1.0 && request != "CAPCHA_NOT_READY" {
			return "", errors.New(request)
		} else if status == 1.0 {
			response = request
			break
		} else {
			time.Sleep(10 * time.Second)
		}
	}

	return response, nil
}

// *** Rude Solving block END ***

// *** ReCaptcha2 block START ***

//ReCap2 is struct for solving ReCaptcha2
type ReCap2 struct {
	Url     string
	Key     string
	Sitekey string
	Pageurl string
}

//Solve is Solvable realization
func (rc *ReCap2) Solve(ch chan capResp) {
	if rc.Sitekey == "" && rc.Pageurl == "" {
		ch <- capResp{response: "", err: errors.New("please SetTarget first")}
	}
	if rc.Key != "" && rc.Url != "" {
		var response string
		request := recapRequest{
			url: rc.Url,
			postParams: map[string]interface{}{
				"key":       rc.Key,
				"method":    "userrecaptcha",
				"googlekey": rc.Sitekey,
				"pageurl":   rc.Pageurl,
				"softguru":  "104431",
				"json":      1,
			},
			getParams: map[string]interface{}{
				"key":    rc.Key,
				"action": "get",
				"id":     0,
				"json":   1,
			},
		}
		response, err := rudeSolver{&request}.solve()
		if err != nil {
			ch <- capResp{response: "", err: err}
		}
		ch <- capResp{response, err}
	}

	ch <- capResp{response: "", err: errors.New("please initialize a struct")}
}

//SetTarget is func
func (rc *ReCap2) SetTarget(sitekey, pageurl string) {
	rc.Sitekey = sitekey
	rc.Pageurl = pageurl
}

// *** ReCaptcha2 block END ***
