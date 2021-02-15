package model

const (
	ProjectName = "project"
)

type AccelDTO struct {
	ID string `json:"id"`
	X  string `json:"x"`
	Y  string `json:"y"`
	Z  string `json:"z"`
}

type Accel struct {
	ID string
	X  float64
	Y  float64
	Z  float64
}
