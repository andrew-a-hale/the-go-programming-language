// Package weightconv performs kilogram and pound conversions
package weightconv

import "fmt"

type Pound float64
type Kilogram float64

const kgToLbFactor Kilogram = 2.20462
const lbToKgFactor Pound = 0.453592

func (lb Pound) String() string    { return fmt.Sprintf("%g lb", lb) }
func (kg Kilogram) String() string { return fmt.Sprintf("%g kg", kg) }
