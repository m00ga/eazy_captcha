package eazycaptcha

import(
	"fmt"
)

type reCapError struct{
	string
}

func (rcErr *reCapError) Error() string{
	return rcErr.string
}

func (rcErr *reCapError) String() string{
	return rcErr.Error()
}

type reCapResponse struct{
	response string
	err error
}

func (rc *reCapResponse) getData() (string, error){
	return rc.response, rc.err
}

//type reCapRequest struct{
	//URL string
	//postParams, getParams map[string] string
//}

type CapType int

const(
	RC2 CapType = iota
	RC3
)

type ReCapSolver struct{
	URL, Key string
	Type CapType
	Settings map[string] string
}

func (rcs *ReCapSolver) Solve(ch chan CapResponse){
	var typeUsuall string

	switch rcs.Type{
		case RC2:
			typeUsuall += fmt.Sprintf("key=%s&method=userrecaptcha&softguru=104431&", rcs.Key)
		case RC3:
			typeUsuall += fmt.Sprintf("key=%s&method=userrecaptcha&version=v3&softguru=104431&", rcs.Key)
	}

	for key, value := range rcs.Settings{
		if key == "key"{ 
			continue
		}
		typeUsuall += fmt.Sprintf("%s=%s&", key, value)
	}

	getParams := fmt.Sprintf("key=%s&action=get", rcs.Key)

	resp, err := rudeSolve(&rudeRequest{rcs.URL, typeUsuall, getParams})

	ch <- &reCapResponse{resp, err}
}