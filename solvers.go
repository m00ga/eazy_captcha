package eazycaptcha

import (
	//"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
	"fmt"
)

// *** Rude Solving block BEGIN ***

func (rs rudeSolver) solve() (string, error) {
	var response string
	client := http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(rs.request.url + "/in.php?" + rs.request.postParams)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var rcresp map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&rcresp)

	if rcresp["status"] != 1.0 {
		return "", errors.New(rcresp["request"].(string))
	}
	rs.request.getParams += fmt.Sprintf("&id=%s", rcresp["request"])

	for {
		resp, err := client.Get(rs.request.url + "/res.php?" + rs.request.getParams)
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
			postParams: fmt.Sprintf(
				"key=%s&method=userrecaptcha&googlekey=%s&pageurl=%s&softguru=104431&json=1",
				rc.Key, rc.Sitekey, rc.Pageurl,
			),
			getParams: fmt.Sprintf(
				"key=%s&action=get&json=1",
				rc.Key, 
			),
		}
		response, err := rudeSolver{&request}.solve()
		if err != nil {
			ch <- capResp{response: "", err: err}
		}
		ch <- capResp{response, err}
	}

	ch <- capResp{response: "", err: errors.New("please initialize a struct")}
}

//SetTarget is func for setting a site
func (rc *ReCap2) SetTarget(sitekey, pageurl string) {
	rc.Sitekey = sitekey
	rc.Pageurl = pageurl
}

// *** ReCaptcha2 block END ***

// *** ReCaptcha3 block START ***

//ReCap3 is struct for solving ReCaptcha3
type ReCap3 struct {
	URL      string
	Key      string
	Sitekey string
	Pageurl string
	MinScore float32
}


//Solve is Solvable realization
func (rc3 *ReCap3) Solve(ch chan capResp) {
	if rc3.Sitekey != "" && rc3.Pageurl != "" && rc3.MinScore != 0.0{
		if url, key := rc3.URL, rc3.Key; url != "" && key != "" {
			
			request := recapRequest{
				url: url,
				postParams: fmt.Sprintf(
					"key=%s&method=userrecaptcha&googlekey=%s&pageurl=%s&version=v3&min_score=%f&softguru=104431&json=1",
					key, rc3.Sitekey, rc3.Pageurl, rc3.MinScore,
				),
				getParams: fmt.Sprintf(
					"key=%s&action=get&json=1",
					key,
				),
			}

			response, err := rudeSolver{&request}.solve()

			if err != nil {
				ch <- capResp{response: "", err: err}
			}
			ch <- capResp{response, err}
		} else {
			ch <- capResp{response: "", err: errors.New("please first initialize struct")}
		}
	} else {
		ch <- capResp{response: "", err: errors.New("please use SetTarget first")}
	}
}



//SetTarget is func for setting a site
func (rc3 *ReCap3) SetTarget(pageurl, sitekey string, minscore float32) {
	rc3.Pageurl = pageurl
	rc3.Sitekey = sitekey
	rc3.MinScore = minscore
}

func isNil(v interface{}) bool {
	return v == nil
}

// *** ReCaptcha3 block END
