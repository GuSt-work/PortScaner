package argsparser

import (
	"reflect"
	"testing"
)

func TestParsePorts(t *testing.T) {

	tests := []struct {
		name  string
		input string
		want  []uint16
	}{
		{
			name:  "single port",
			input: "80",
			want:  []uint16{80},
		},

		{
			name:  "multiple ports",
			input: "22,80,443",
			want:  []uint16{22, 80, 443},
		},

		{
			name:  "port range",
			input: "20-25",
			want:  []uint16{20, 21, 22, 23, 24, 25},
		},

		{
			name:  "mixed ports",
			input: "20-25,80,443",
			want:  []uint16{20, 21, 22, 23, 24, 25, 80, 443},
		},

		{
			name:  "duplicate ports",
			input: "80,80,443,80",
			want:  []uint16{80, 443},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			got, err := ParsePorts(tt.input)

			if err != nil {
				t.Fatalf(
					"unexpected error: %v",
					err,
				)
			}

			if !reflect.DeepEqual(got, tt.want) {

				t.Errorf(
					"got %v, want %v",
					got,
					tt.want,
				)
			}
		})
	}
}

func TestParsePortsErrors(t *testing.T) {

	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "zero port",
			input: "0",
		},

		{
			name:  "too large port",
			input: "65536",
		},

		{
			name:  "invalid text",
			input: "abc",
		},

		{
			name:  "invalid range",
			input: "100-50",
		},

		{
			name:  "range overflow",
			input: "65530-70000",
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			_, err := ParsePorts(tt.input)

			if err == nil {

				t.Errorf(
					"expected error for %s",
					tt.input,
				)

			}

		})
	}
}
