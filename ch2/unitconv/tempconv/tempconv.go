// Package tempconv performs Celsius, Fahrenheit, and Kelvin conversions.
package tempconv

import "fmt"

type (
	Celsius    float64
	Fahrenheit float64
	Kelvin     float64
)

const (
	AbsoluteZeroC Celsius = -273.15
	FreezingC     Celsius = 0
	BoilingC      Celsius = 100
)

// Celsius to string
func (c Celsius) String() string { return fmt.Sprintf("%g°C", c) }

// Fahrenheit to string
func (f Fahrenheit) String() string { return fmt.Sprintf("%g°F", f) }

// Kelvin to string
func (k Kelvin) String() string { return fmt.Sprintf("%gK", k) }
