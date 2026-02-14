package distance

import "math"

const (
	EarthRadiusMeters = 6371000.0
	YardsPerMeter = 1.09361
)

// Distance between two geographic coordinates using the Haversine formula (yay trig)
func Calculate(lat1, lon1, lat2, lon2 float64) float64 {
	dLat := toRadians(lat2 - lat1)
	dLon := toRadians(lon2 - lon1)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(toRadians(lat1))*math.Cos(toRadians(lat2))*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return EarthRadiusMeters * c
}

func YardsToMeters(yards float64) float64 {
	return yards / YardsPerMeter
}

func MetersToYards(meters float64) float64 {
	return meters * YardsPerMeter
}

func toRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}