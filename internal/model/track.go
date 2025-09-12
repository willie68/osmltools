package model

import "time"

type Waypoint struct {
	Name  string    `json:"name,omitempty"`
	Lat   float64   `json:"latitude,omitempty"`
	Lon   float64   `json:"longitude,omitempty"`
	Time  time.Time `json:"time,omitempty"`
	Speed float64   `json:"speed,omitempty"`
	Ele   float64   `json:"elevation,omitempty"`
	Depth float64   `json:"depth,omitempty"`
}

type Track struct {
	Name      string      `json:"name,omitempty"`
	Waypoints []*Waypoint `json:"waypoints,omitempty"`
	Start     *Waypoint   `json:"start,omitempty"`
	End       *Waypoint   `json:"end,omitempty"`
	LogLines  []*LogLine  `json:"log_lines,omitempty"`
}
