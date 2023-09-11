package sync

import (
	"reflect"
	"testing"
)

func Test_compareMaps(t *testing.T) {
	type args struct {
		ori map[string]string
		oth map[string]string
	}
	tests := []struct {
		name      string
		args      args
		wantToAdd map[string]string
		wantToDel map[string]string
	}{
		{
			name: "",
			args: args{
				ori: map[string]string{"abc": "", "def": ""},
				oth: map[string]string{"def": "fgi", "klm": ""},
			},
			wantToAdd: map[string]string{"def": "fgi", "klm": ""},
			wantToDel: map[string]string{"abc": ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToAdd, gotToDel := compareMaps(tt.args.ori, tt.args.oth)
			if !reflect.DeepEqual(gotToAdd, tt.wantToAdd) {
				t.Errorf("compareMaps() gotToAdd = %v, want %v", gotToAdd, tt.wantToAdd)
			}
			if !reflect.DeepEqual(gotToDel, tt.wantToDel) {
				t.Errorf("compareMaps() gotToDel = %v, want %v", gotToDel, tt.wantToDel)
			}
		})
	}
}
