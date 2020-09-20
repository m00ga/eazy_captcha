 //Realization a strategy pattern for captcha solving

package eazycaptcha

import (
	"errors"
	"sync"
)

//Solvable is abstract realization for solver structs
type Solvable interface {
	Solve(ch chan capResp)
}

type rudeSolver struct {
	request *recapRequest
}

type recapRequest struct {
	url        string
	postParams map[string]interface{}
	getParams  map[string]interface{}
}

//Solver is struct
type Solver struct {
	alghoritm Solvable
	mu sync.Mutex
	counter   int
	tasks     map[int]chan capResp
}

type capResp struct {
	response string
	err      error
}

//Alghoritm is func for change solving alghoritm
func (s *Solver) Alghoritm(alg Solvable) {
	s.alghoritm = alg
}

//Solve is func for solve captcha with selected alghoritm
func (s *Solver) Solve() int {
	s.mu.Lock()
	if s.tasks == nil {
		s.tasks = make(map[int]chan capResp)
	}
	curr := s.counter
	ch := make(chan capResp)
	s.tasks[curr] = ch
	s.counter++
	s.mu.Unlock()
	go s.alghoritm.Solve(ch)
	return curr
}

//Get is func for get a result by id
func (s *Solver) Get(id int) (string, error) {
	if ch := s.tasks[id]; ch != nil {
		defer delete(s.tasks, id)
		defer s.decrement()
		if val, opened := <-ch; opened {
			return val.response, val.err
		}

		return "", errors.New("channel closed")
	}

	return "", errors.New("please use a valid id")
}

func (s *Solver) decrement() {
	s.mu.Lock()
	s.counter--
	s.mu.Unlock()
}

