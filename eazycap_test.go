package eazycaptcha

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

const (
	URL = "cap service URL"
	API = "cap service API key"
)

func TestReCap3(t *testing.T) {
	client := &http.Client{}
	solver := new(Solver)
	rc3 := ReCap3{
		URL:      URL,
		Key:      API,
		Sitekey:  "6LdyC2cUAAAAACGuDKpXeDorzUDWXmdqeg-xy696",
		Pageurl:  "https://recaptcha-demo.appspot.com/recaptcha-v3-request-scores.php",
		MinScore: 0.7,
	}

	solver.Alghoritm(&rc3)

	id := solver.Solve()
	resp, err := solver.Get(id)
	if err != nil {
		t.Fatalf("FATAL ERROR: details %v", err)
	}

	req, err := json.Marshal(map[string]string{
		"action": "examples/v3scores",
		"token":  resp,
	})
	if err != nil {
		t.Fatalf("FATAL ERROR: details %v", err)
	}

	res, err := client.Post("https://recaptcha-demo.appspot.com/recaptcha-v3-verify.php", "application/json", bytes.NewReader(req))
	if err != nil {
		t.Fatalf("FATAL ERROR: details %v", err)
	}
	defer res.Body.Close()

	var resBody map[string]interface{}

	err = json.NewDecoder(res.Body).Decode(&resBody)
	if err != nil {
		t.Fatalf("FATAL ERROR: details %v", err)
	}

	fmt.Println(resp)

	fmt.Scanln()
}

func TestReCap2(t *testing.T) {
	//client := &http.Client{}
	solver := new(Solver)
	rc2 := ReCap2{
		URL,
		API,
		"6LfW6wATAAAAAHLqO2pb8bDBahxlMxNdo9g947u9",
		"https://recaptcha-demo.appspot.com/recaptcha-v2-checkbox.php",
	}

	solver.Alghoritm(&rc2)
	id := solver.Solve()
	resp, err := solver.Get(id)
	if err != nil {
		t.Fatalf("FATAL ERROR: details %v", err)
	}

	fmt.Println(resp)
}
