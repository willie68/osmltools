package model

import "time"

type Waypoint struct {
	Name  string
	Lat   float64
	Lon   float64
	Time  time.Time
	Speed float64
	Ele   float64
	Depth float64
}

type Track struct {
	Name      string
	Waypoints []*Waypoint
	Start     *Waypoint
	End       *Waypoint
	LogLines  []*LogLine
}
