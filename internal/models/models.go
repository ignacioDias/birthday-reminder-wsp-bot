package models

import "time"

type Birthday struct {
	Month time.Month
	Day   int
	Name  string
}
