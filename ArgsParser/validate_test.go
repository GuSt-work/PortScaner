package argsparser

import (
	"reflect"
	"testing"
	"time"
)

func TestValidateArgs(t *testing.T) {

	tests := []struct {
		name         string
		args         []string
		wantTargets  int
		wantPorts    []uint16
		wantAllPorts bool
		wantTimeout  time.Duration
		expectError  bool
	}{
		{
			name:        "default values",
			args:        []string{"192.168.1.10"},
			wantTargets: 1,
			wantPorts:   []uint16{80},
			wantTimeout: time.Second,
		},
		{
			name:        "custom ports",
			args:        []string{"-p", "22,80,443", "192.168.1.10"},
			wantTargets: 1,
			wantPorts:   []uint16{22, 80, 443},
			wantTimeout: time.Second,
		},
		{
			name:        "custom timeout",
			args:        []string{"-t", "500", "192.168.1.10"},
			wantTargets: 1,
			wantPorts:   []uint16{80},
			wantTimeout: 500 * time.Millisecond,
		},
		{
			name:        "multiple targets",
			args:        []string{"-p", "80", "192.168.1.10", "192.168.1.20"},
			wantTargets: 2,
			wantPorts:   []uint16{80},
			wantTimeout: time.Second,
		},
		{
			name:        "port range",
			args:        []string{"-p", "20-22", "192.168.1.10"},
			wantTargets: 1,
			wantPorts:   []uint16{20, 21, 22},
			wantTimeout: time.Second,
		},
		{
			name:         "all ports",
			args:         []string{"-p", "-", "192.168.1.10"},
			wantTargets:  1,
			wantAllPorts: true,
			wantTimeout:  time.Second,
		},
		{
			name:        "no targets",
			args:        []string{"-p", "80"},
			expectError: true,
		},
		{
			name:        "invalid port",
			args:        []string{"-p", "70000", "192.168.1.10"},
			expectError: true,
		},
		{
			name:        "invalid timeout",
			args:        []string{"-t", "0", "192.168.1.10"},
			expectError: true,
		},
		{
			name:        "invalid target",
			args:        []string{"999.999.999.999"},
			expectError: true,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			cfg, err := ValidateArgs(tt.args)

			if tt.expectError {

				if err == nil {
					t.Fatal("expected error")
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(cfg.Targets) != tt.wantTargets {
				t.Errorf(
					"targets: got %d want %d",
					len(cfg.Targets),
					tt.wantTargets,
				)
			}

			if tt.wantAllPorts {

				if len(cfg.Ports) != 65535 {
					t.Fatalf(
						"expected 65535 ports, got %d",
						len(cfg.Ports),
					)
				}

				if cfg.Ports[0] != 1 {
					t.Errorf(
						"first port: got %d want 1",
						cfg.Ports[0],
					)
				}

				if cfg.Ports[len(cfg.Ports)-1] != 65535 {
					t.Errorf(
						"last port: got %d want 65535",
						cfg.Ports[len(cfg.Ports)-1],
					)
				}

			} else {

				if !reflect.DeepEqual(cfg.Ports, tt.wantPorts) {
					t.Errorf(
						"ports: got %v want %v",
						cfg.Ports,
						tt.wantPorts,
					)
				}

			}

			if cfg.Timeout != tt.wantTimeout {
				t.Errorf(
					"timeout: got %v want %v",
					cfg.Timeout,
					tt.wantTimeout,
				)
			}

		})

	}

}
