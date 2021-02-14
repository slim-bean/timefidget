package model

const (
	ProjectName = "project"
)

type AccelDTO struct {
	X string `json:"x"`
	Y string `json:"y"`
	Z string `json:"z"`
}

type Accel struct {
	X float64
	Y float64
	Z float64
}
