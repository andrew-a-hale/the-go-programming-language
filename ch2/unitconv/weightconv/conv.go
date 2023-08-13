package weightconv

// lbToKg converts from Pounds to Kilograms
func LbToKg(lb Pound) Kilogram { return Kilogram(lb * lbToKgFactor) }

// kgToLb converts from Kilograms to Pounds
func KgToLb(kg Kilogram) Pound { return Pound(kg * kgToLbFactor) }
