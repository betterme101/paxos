package paxos

import (
	"testing"
)

func start(acceptorIds []int, learnerIds []int) ([]*Acceptor, []*Learner) {
	acceptors := make([]*Acceptor, 0)
	for _, aid := range acceptorIds {
		a := newAcceptor(aid, learnerIds)
		// not :=
		acceptors = append(acceptors, a)
	}

	learners := make([]*Learner, 0)
	for _, lid := range learnerIds {
		l := newLearner(lid, acceptorIds)
		learners = append(learners, l)
	}

	return acceptors, learners
}

func cleanup(acceptors []*Acceptor, learners []*Learner) {
	for _, a := range acceptors {
		a.close()
	}

	for _, l := range learners {
		l.close()
	}
}

func TestSingleProposer(t *testing.T) {
	acceptorIds := []int {1001, 1002, 1003}

	learnerIds := []int {2001}
	acceptors, learners := start(acceptorIds, learnerIds)

	defer cleanup(acceptors, learners)

	p := &Proposer {
		id: 1,
		acceptors: acceptorIds,
	}

	value := p.propose("hello world")
	if value != "hello world" {
		t.Errorf("learnValue = %s, expected %s", value, "hello world")
	}

	learnValue := learners[0].chosen()
	if learnValue != value {
		t.Errorf("learnValue = %s, expected %s", learnValue, "hello world")
	}
}

func TestTwoProposers(t *testing.T) {
	acceptorIds := []int {1001, 1002, 1003}
	learnerIds := []int {2001}
	acceptors, learners := start(acceptorIds, learnerIds)

	defer cleanup(acceptors, learners)

	p1 := &Proposer{
		id: 1,
		acceptors: acceptorIds,
	}
	v1 := p1.propose("hello world")

	p2 := &Proposer{
		id: 2,
		acceptors: acceptorIds,
	}
	v2 := p2.propose("hello can")

	if v1 != v2 {
		t.Errorf("value1 = %s, value2 = %s", v1, v2)
	}

	learnValue := learners[0].chosen()
	if learnValue != v1 {
		t.Errorf("learnValue = %s, excepted %s", learnValue, v1)
	}
}