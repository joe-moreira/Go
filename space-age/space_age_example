package space

// Planet represents the planet on which we want to calculate the age
type Planet string

// The duration of one complete revolution of earth in seconds
const earthRevolution = 31557600

// Age calculates the age relative to a planet rotation cycle
func Age(ageInSeconds float64, planet Planet) float64 {
	// Planets represents the association between a planet and it's revolution
	Planets := map[Planet]float64{
		"Earth":   earthRevolution,
		"Mercury": earthRevolution * 0.2408467,
		"Venus":   earthRevolution * 0.61519726,
		"Mars":    earthRevolution * 1.8808158,
		"Jupiter": earthRevolution * 11.862615,
		"Saturn":  earthRevolution * 29.447498,
		"Uranus":  earthRevolution * 84.016846,
		"Neptune": earthRevolution * 164.79132,
	}
	return calculateAge(Planets[planet], ageInSeconds)
}

func calculateAge(planetRevolution float64, ageInSeconds float64) float64 {
	return ageInSeconds / planetRevolution
}
