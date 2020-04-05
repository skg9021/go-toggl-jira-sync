package main

import (
	"os"
	"strconv"
)

// Config to run the program
type Config struct {
	Days      int
	Frequency int
}

// NewConfig is the constructor for config to set
// default values
func NewConfig() Config {
	config := Config{}
	days, _ := strconv.Atoi(os.Getenv("DAYS"))
	if days == 0 {
		config.Days = 7
	} else {
		config.Days = days
	}

	frequency, _ := strconv.Atoi(os.Getenv("FREQUENCY"))
	if frequency == 0 {
		config.Frequency = 60
	} else {
		config.Frequency = frequency
	}

	return config
}
