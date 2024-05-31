package vcdusage

import (
	"math"
	"strconv"
)

// KB factor
const kb = 1_024

// MB factor
const mb = 1_048_576

// GB factor
const gb = 1_073_741_820

// TB factor
const tb = 1_099_511_630_000

// DataStorage is a representation of a data storage amount in bytes.
type DataStorage float64

// Float64 converts the DataStorage value to float64.
func (d DataStorage) Float64() float64 {
	return math.Round(float64(d))
}

// Uint64 converts the DataStorage value to uint64.
func (d DataStorage) Uint64() uint64 {
	return uint64(math.Round(d.Float64()))
}

// Int64 converts the DataStorage value to int64.
func (d DataStorage) Int64() int64 {
	return int64(math.Round(d.Float64()))
}

// String converts the DataStorage value to string.
func (d DataStorage) String() string {
	return strconv.FormatFloat(d.Float64(), 'f', -1, 64)
}

// KB converts the DataStorage value to KB as float64 with 0 decimal points.
func (d DataStorage) KB() float64 {
	return d.convert(kb)
}

// MB converts the DataStorage value to MB as float64 with 0 decimal points.
func (d DataStorage) MB() float64 {
	return d.convert(mb)
}

// GB converts the DataStorage value to GB as float64 with 0 decimal points.
func (d DataStorage) GB() float64 {
	return d.convert(gb)
}

// TB converts the DataStorage value to TB as float64 with a maximum of 2 decimal points.
func (d DataStorage) TB() float64 {
	return math.Round((float64(d)/tb)*10) / 10
}

func (d DataStorage) convert(factor float64) float64 {
	return math.Round(float64(d) / factor)
}
