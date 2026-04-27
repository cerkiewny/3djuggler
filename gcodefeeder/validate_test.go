package gcodefeeder

import (
	"strings"
	"testing"
)

func TestValidateGcode(t *testing.T) {
	tests := []struct {
		name    string
		gcode   string
		wantErr bool
	}{
		{
			name:    "safe temperature",
			gcode:   "M104 S200\nG28\n",
			wantErr: false,
		},
		{
			name:    "exactly at limit",
			gcode:   "M104 S220\nG28\n",
			wantErr: false,
		},
		{
			name:    "M104 over limit",
			gcode:   "M104 S250\nG28\n",
			wantErr: true,
		},
		{
			name:    "M109 over limit",
			gcode:   "G28\nM109 S260\n",
			wantErr: true,
		},
		{
			name:    "temperature in comment ignored",
			gcode:   "; M104 S300\nG28\n",
			wantErr: false,
		},
		{
			name:    "no temperature commands",
			gcode:   "G28\nG1 X10 Y10\n",
			wantErr: false,
		},
		{
			name:    "bed temperature ignored",
			gcode:   "M140 S300\nG28\n",
			wantErr: false,
		},
		{
			name:    "inline comment with high temp",
			gcode:   "M104 S250 ; set temp\n",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGcode(strings.NewReader(tt.gcode))
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGcode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && err.Error() != "too high temperature set, please use PLA presets only" {
				t.Errorf("unexpected error message: %v", err)
			}
		})
	}
}
