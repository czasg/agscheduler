package agscheduler

import "testing"

func TestSTATUS_IsPaused(t *testing.T) {
	tests := []struct {
		name string
		s    STATUS
		want bool
	}{
		{
			name: "paused",
			s:    STATUS_PAUSED,
			want: true,
		},
		{
			name: "not paused",
			s:    STATUS_RUNNING,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.IsPaused(); got != tt.want {
				t.Errorf("IsPaused() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSTATUS_IsRunning(t *testing.T) {
	tests := []struct {
		name string
		s    STATUS
		want bool
	}{
		{
			name: "running",
			s:    STATUS_RUNNING,
			want: true,
		},
		{
			name: "not running",
			s:    STATUS_PAUSED,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.IsRunning(); got != tt.want {
				t.Errorf("IsRunning() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSTATUS_IsStopped(t *testing.T) {
	tests := []struct {
		name string
		s    STATUS
		want bool
	}{
		{
			name: "stopped",
			s:    STATUS_STOPPED,
			want: true,
		},
		{
			name: "not stopped",
			s:    STATUS_RUNNING,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.IsStopped(); got != tt.want {
				t.Errorf("IsStopped() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSTATUS_SetPaused(t *testing.T) {
	tests := []struct {
		name string
		s    STATUS
	}{
		{
			name: "pass",
			s:    STATUS_RUNNING,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.SetPaused()
			if got := tt.s.IsPaused(); got != true {
				t.Errorf("IsStopped() = %v, want %v", got, true)
			}
		})
	}
}

func TestSTATUS_SetRunning(t *testing.T) {
	tests := []struct {
		name string
		s    STATUS
	}{
		{
			name: "pass",
			s:    STATUS_STOPPED,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.SetRunning()
			if got := tt.s.IsRunning(); got != true {
				t.Errorf("IsStopped() = %v, want %v", got, true)
			}
		})
	}
}

func TestSTATUS_SetStopped(t *testing.T) {
	tests := []struct {
		name string
		s    STATUS
	}{
		{
			name: "pass",
			s:    STATUS_RUNNING,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.SetStopped()
			if got := tt.s.IsStopped(); got != true {
				t.Errorf("IsStopped() = %v, want %v", got, true)
			}
		})
	}
}
