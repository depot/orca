package chunk

import (
	"reflect"
	"testing"
)

func TestNewChunksFromString(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		want    Chunks
		wantErr bool
	}{
		{s: "", want: Chunks{}},
		{s: "000002654127863bf5fcd2ad", want: Chunks{{Size: 613, Hash: 4694868728645210797}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewChunksFromString(tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewChunksFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewChunksFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChunks_String(t *testing.T) {
	tests := []struct {
		name string
		c    Chunks
		want string
	}{
		{c: Chunks{}, want: ""},
		{c: Chunks{{Size: 613, Hash: 4694868728645210797}}, want: "000002654127863bf5fcd2ad"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(); got != tt.want {
				t.Errorf("Chunks.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
