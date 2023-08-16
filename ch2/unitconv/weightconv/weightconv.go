// Package weightconv performs kilogram and pound conversions
package weightconv

import "fmt"

type (
	Pound    float64
	Kilogram float64
)

const (
	kgToLbFactor Kilogram = 2.20462
	lbToKgFactor Pound    = 0.453592
)

// Pound to string
func (lb Pound) String() string { return fmt.Sprintf("%g lb", lb) }

// Kilogram to string
func (kg Kilogram) String() string { return fmt.Sprintf("%g kg", kg) }
