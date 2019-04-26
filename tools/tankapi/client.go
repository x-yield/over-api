package tankapi

type Firestarter struct{}

func NewFirestarter() *Firestarter {
	return &Firestarter{}
}

// Close - does nothing. made for overall consistency
// May be it could cleanup sessions or something.
func (*Firestarter) Close() error {
	return nil
}

// Validate - sends config into corresponding tank apis to validate them
func (*Firestarter) Validate(sessions []*TankSession) []*TankSession {
	c := make(chan error, len(sessions))
	for _, s := range sessions {
		go func(s *TankSession, c chan<- error) {
			c <- s.validate()
		}(s, c)
	}
	for i := 0; i < len(sessions); i++ {
		<-c
	}
	return sessions
}

// Prepare - starts tankapi sessions with breakpoint set to run, so that tanks will prepare to be started
// ??? validate before preparing ???
func (*Firestarter) Prepare(sessions []*TankSession) []*TankSession {
	c := make(chan error, len(sessions))
	for _, s := range sessions {
		go func(s *TankSession, c chan<- error) {
			c <- s.prepare()
		}(s, c)
	}
	for i := 0; i < len(sessions); i++ {
		<-c
	}
	return sessions
}

// Run - sets tankapi sessions breakpoint to finished, so that tanks will run at once
func (*Firestarter) Run(sessions []*TankSession) []*TankSession {
	c := make(chan error, len(sessions))
	for _, s := range sessions {
		go func(s *TankSession, c chan<- error) {
			c <- s.run()
		}(s, c)
	}
	for i := 0; i < len(sessions); i++ {
		<-c
	}
	return sessions
}

// Stop - stops tankapi sessions
func (*Firestarter) Stop(sessions []*TankSession) []*TankSession {
	c := make(chan error, len(sessions))
	for _, s := range sessions {
		go func(s *TankSession, c chan<- error) {
			c <- s.stop()
		}(s, c)
	}
	for i := 0; i < len(sessions); i++ {
		<-c
	}
	return sessions
}

// Poll - polls tankapi sessions' status
func (*Firestarter) Poll(sessions []*TankSession) []*TankSession {
	c := make(chan error, len(sessions))
	for _, s := range sessions {
		go func(s *TankSession, c chan<- error) {
			c <- s.poll()
		}(s, c)
	}
	for i := 0; i < len(sessions); i++ {
		<-c
	}
	return sessions
}
