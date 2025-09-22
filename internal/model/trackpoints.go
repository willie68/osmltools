package model

import (
	"time"

	"github.com/adrianmo/go-nmea"
)

type ThreePoints struct {
	X int64 `json:"x,omitempty"`
	Y int64 `json:"y,omitempty"`
	Z int64 `json:"z,omitempty"`
}

// Waypoint internal waypoint structure
type Waypoint struct {
	Name         string       `json:"name,omitempty"`
	Lat          float64      `json:"latitude,omitempty"`
	Lon          float64      `json:"longitude,omitempty"`
	Time         time.Time    `json:"time,omitempty"`
	Speed        float64      `json:"speed,omitempty"`
	Ele          float64      `json:"elevation,omitempty"`
	Depth        float64      `json:"depth,omitempty"`
	Acceleration *ThreePoints `json:"acc,omitempty"`
	GyroLocation *ThreePoints `json:"gyro,omitempty"`
	Supply       int64        `json:"supply,omitempty"`
}

type TrackPoints struct {
	Name      string      `json:"name,omitempty"`
	Waypoints []*Waypoint `json:"waypoints,omitempty"`
	Start     *Waypoint   `json:"start,omitempty"`
	End       *Waypoint   `json:"end,omitempty"`
	LogLines  []*LogLine  `json:"log_lines,omitempty"`
}

// GetWaypoints extracts the waypoints from the log lines of the track
func GetWaypoints(track *TrackPoints) (*TrackPoints, error) {
	track.Waypoints = make([]*Waypoint, 0)

	for _, ll := range track.LogLines {
		if ll.NMEAMessage != nil {
			switch ll.NMEAMessage.Prefix() {
			case "GPRMC":
				rmc, ok := ll.NMEAMessage.(nmea.RMC)
				if ok && rmc.Validity == "A" { // only valid
					track.End = &Waypoint{
						Lat:   rmc.Latitude,
						Lon:   rmc.Longitude,
						Time:  ll.CorrectTimeStamp,
						Speed: rmc.Speed,
						Ele:   0.0,
					}
					track.Waypoints = append(track.Waypoints, track.End)
					if track.Start == nil {
						track.Start = track.End
					}
				}
			case "GPGGA":
				if track.End != nil {
					gga, ok := ll.NMEAMessage.(nmea.GGA)
					if ok {
						if track.End.Ele == 0.0 {
							track.End.Ele = gga.Altitude
						}
					}
				}
			case "SDDBT":
				if track.End != nil {
					dbt, ok := ll.NMEAMessage.(nmea.DBT)
					if ok {
						depth := dbt.DepthFeet * 0.3048 // convert feet to meters
						if track.End.Depth == 0.0 {
							track.End.Depth = depth
						}
					}
				}
			case "SDDPT":
				if track.End != nil {
					dpt, ok := ll.NMEAMessage.(nmea.DPT)
					if ok {
						if track.End.Depth == 0.0 {
							track.End.Depth = dpt.Depth
						}
					}
				}
			}
		}
	}
	if track.Start != nil {
		track.Start.Name = "Start"
	}
	if track.End != nil {
		track.End.Name = "End"
	}
	return track, nil
}
