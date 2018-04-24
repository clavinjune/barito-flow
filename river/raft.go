package river

import "github.com/BaritoLog/go-boilerplate/errkit"

const (
	StoreError = errkit.Error("Error when store timber")
)

// Raft
type Raft interface {
	Start()
	ErrorChannel() (errCh chan error)
}

type raft struct {
	from  Upstream
	to    Downstream
	errCh chan error
}

func NewRaft(from Upstream, to Downstream) Raft {
	return &raft{
		from:  from,
		to:    to,
		errCh: make(chan error),
	}
}

// Drifting
func (r *raft) Start() {
	r.from.SetErrorChannel(r.errCh)

	go r.from.StartTransport()
	go r.start()

}

func (r *raft) start() {
	timberCh := r.from.TimberChannel()
	for {
		select {
		case timber := <-timberCh:
			err := r.to.Store(timber)
			if err != nil {
				r.errCh <- errkit.Concat(StoreError, err)
			}
		}
	}
}

func (r *raft) ErrorChannel() (errCh chan error) {
	return r.errCh
}
