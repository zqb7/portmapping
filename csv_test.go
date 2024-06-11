package portmapping

import (
	"reflect"
	"testing"
)

func TestLoadCSV(t *testing.T) {
	tests := []struct {
		name    string
		want    []*Item
		wantErr bool
	}{
		{name: "1", want: []*Item{
			{8080, "tcp", true, "127.0.0.1", 80, ""},
		}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadCSV("mapping.csv.example")
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadCSV() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadCSV() = %v, want %v", got, tt.want)
			}
		})
	}
}
