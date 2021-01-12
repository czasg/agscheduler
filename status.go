package agscheduler

type STATUS string

var (
	STATUS_RUNNING = STATUS("RUNNING")
	STATUS_PAUSED  = STATUS("PAUSED")
	STATUS_STOPPED = STATUS("STOPPED")
)

func (s *STATUS) IsRunning() bool {
	return *s == STATUS_RUNNING
}

func (s *STATUS) IsPaused() bool {
	return *s == STATUS_PAUSED
}

func (s *STATUS) IsStopped() bool {
	return *s == STATUS_STOPPED
}

func (s *STATUS) SetRunning() {
	*s = STATUS_RUNNING
}

func (s *STATUS) SetPaused() {
	*s = STATUS_PAUSED
}

func (s *STATUS) SetStopped() {
	*s = STATUS_STOPPED
}
