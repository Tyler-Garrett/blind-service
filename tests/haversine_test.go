package tests

import (
	"math"
	"testing"
	"github.com/tyler-garrett/blind-service/pkg/distance"
)

func TestCalculate(t *testing.T) {
	tests := []struct {
		name     string
		lat1     float64
		lon1     float64
		lat2     float64
		lon2     float64
		expected float64
		delta    float64 // acceptable difference in meters
	}{
		{
			name:     "Same location",
			lat1:     37.5,
			lon1:     -76.7,
			lat2:     37.5,
			lon2:     -76.7,
			expected: 0,
			delta:    0.1,
		},
		{
			name:     "500 yards apart (approximately)",
			lat1:     37.5,
			lon1:     -76.7,
			lat2:     37.50414,
			lon2:     -76.7,
			expected: 457.2, // 500 yards in meters
			delta:    10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := distance.Calculate(tt.lat1, tt.lon1, tt.lat2, tt.lon2)
			if math.Abs(result-tt.expected) > tt.delta {
				t.Errorf("Calculate() = %v, want %v (±%v)", result, tt.expected, tt.delta)
			}
		})
	}
}

func TestYardsToMeters(t *testing.T) {
	yards := 500.0
	expected := 457.2
	result := distance.YardsToMeters(yards)
	if math.Abs(result-expected) > 0.1 {
		t.Errorf("YardsToMeters(%v) = %v, want %v", yards, result, expected)
	}
}

func TestMetersToYards(t *testing.T) {
	meters := 457.2
	expected := 500.0
	result := distance.MetersToYards(meters)
	if math.Abs(result-expected) > 0.1 {
		t.Errorf("MetersToYards(%v) = %v, want %v", meters, result, expected)
	}
}