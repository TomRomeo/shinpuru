package dmdialog

type Awaiter struct {
	answers       map[string]string
	removeHandler func()
	cFinished     chan Result
	finished      bool
}

func (a *Awaiter) setResult(r Result) {
	if a.finished {
		return
	}
	a.finished = true
	a.removeHandler()
	a.cFinished <- r
}

func (a *Awaiter) Await() (r Result, m map[string]string) {
	r = <-a.cFinished
	m = a.answers
	return
}
