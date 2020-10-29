package AGScheduler

import "testing"

func TestRegisterWorksMap(t *testing.T) {
	type args struct {
		worksMap map[string]WorkDetail
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "register",
			args: args{
				worksMap: map[string]WorkDetail{
					"work": {
						Func: func(args ...interface{}) {
						},
						Args: []interface{}{},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RegisterWorksMap(tt.args.worksMap); (err != nil) != tt.wantErr {
				t.Errorf("RegisterWorksMap() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(WorksMap) == 0 {
				t.Error("register error")
			}
		})
	}
}
