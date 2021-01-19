package agscheduler

import (
	"github.com/sirupsen/logrus"
	"reflect"
	"testing"
)

func TestGenAGSDetails(t *testing.T) {
	tests := []struct {
		name string
		want logrus.Fields
	}{
		{
			name: "pass",
			want: logrus.Fields{
				"ASGAuthor":  Author,
				"AGSGitHub":  GitHub,
				"AGSGitBook": GitBook,
				"AGSEmail":   Email,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenAGSDetails(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenAGSDetails() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenAGSVersion(t *testing.T) {
	tests := []struct {
		name string
		want logrus.Fields
	}{
		{
			name: "pass",
			want: logrus.Fields{"AGSVersion": Version},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenAGSVersion(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenAGSVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenASGModule(t *testing.T) {
	type args struct {
		module string
	}
	tests := []struct {
		name string
		args args
		want logrus.Fields
	}{
		{
			name: "pass",
			args: args{
				module: "test",
			},
			want: logrus.Fields{"ASGModule": "test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenASGModule(tt.args.module); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenASGModule() = %v, want %v", got, tt.want)
			}
		})
	}
}
