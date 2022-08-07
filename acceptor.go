package paxos

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)

type Acceptor struct {
	lis net.Listener

	// server id
	id int

	// proposal number that the acceptor has commited, start from 0
	minProposal int

	// proposal number that the acceptor has accepted, start from 0
	acceptedNumber int

	// proposal value that the acceptor has accepted, nil if no 
	acceptedValue interface{}

	// id list of learners 
	learners []int
}


func newAcceptor(id int, learners []int) *Acceptor {
	acceptor := &Acceptor{
		id: id,
		learners: learners,
	}
	acceptor.server()
	return acceptor
}

func (a *Acceptor) Prepare(args *MsgArgs, reply *MsgReply) error {
	if args.Number > a.minProposal {
		a.minProposal = args.Number
		reply.Number = a.acceptedNumber
		reply.Value = a.acceptedValue
		reply.Ok = true
	} else {
		reply.Ok = false
	}
	return nil
}

func (a *Acceptor) Accept(args *MsgArgs, reply *MsgReply) error {
	// can > be correct? for catch the missing may be right, but seems this will
	// also introduce other error cases?
	if args.Number >= a.minProposal {
		a.minProposal = args.Number
		a.acceptedNumber = args.Number
		a.acceptedValue = args.Value
		reply.Ok = true
		// forwarding the accepted proposal to learners
		for _, lid := range a.learners {
			go func(learner int) {
				addr := fmt.Sprintf("127.0.0.1:%d", learner)
				args.From = a.id
				args.To = learner
				resp := new(MsgReply)
				ok := call(addr, "Learner.Learn", args, resp)
				if !ok {
					return 
				}
			}(lid)
		}
	} else {
		reply.Ok = false
	}
	return nil
}

func (a *Acceptor) server() {
	rpcs := rpc.NewServer()
	rpcs.Register(a)
	addr := fmt.Sprintf(":%d", a.id)
	l, e := net.Listen("tcp", addr)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	a.lis = l 
	go func() {
		for {
			conn, err := a.lis.Accept()
			if err != nil {
				continue
			}
			go rpcs.ServeConn(conn)
		}
	}()
}

// like a wrapper method
func (a *Acceptor) close() {
	a.lis.Close()
}