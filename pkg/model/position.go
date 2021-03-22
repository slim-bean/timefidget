package model

import (
	"math"

	"github.com/go-kit/kit/log/level"

	"timefidget/pkg/util"
)

const (
	ON_MIN   = 8
	ON_MAX   = 11
	OFF_MIN  = -1
	OFF_MAX  = 1
	HALF_MIN = 5
	HALF_MAX = 8
	Z_THRESH = 5

	P0 = ""
	P1 = "Loki Ops"
	P2 = "Loki Community"
	P3 = "Loki"
	P4 = "Hiring"
	P5 = "Sales"
	P6 = "1-1"
	P7 = "Management"
	P8 = "BAU"
)

func LogPosition(id string, x, y, z float64) {
	if math.Abs(z) > Z_THRESH {
		// Off
		//level.Info(util.Logger).Log("pos", "0")
	} else if x > OFF_MIN && x < OFF_MAX && y < -ON_MIN && y > -ON_MAX {
		// Position 1
		level.Info(util.Logger).Log("id", id, "type", "add", "pos", "1", ProjectName, P1)
	} else if x > HALF_MIN && x < HALF_MAX && y < -HALF_MIN && y > -HALF_MAX {
		// Position 2
		level.Info(util.Logger).Log("id", id, "type", "add", "pos", "2", ProjectName, P2)
	} else if x > ON_MIN && x < ON_MAX && y > OFF_MIN && y < OFF_MAX {
		// Position 3
		level.Info(util.Logger).Log("id", id, "type", "add", "pos", "3", ProjectName, P3)
	} else if x > HALF_MIN && x < HALF_MAX && y > HALF_MIN && y < HALF_MAX {
		// Position 4
		level.Info(util.Logger).Log("id", id, "type", "add", "pos", "4", ProjectName, P4)
	} else if x > OFF_MIN && x < OFF_MAX && y > ON_MIN && y < ON_MAX {
		// Position 5
		level.Info(util.Logger).Log("id", id, "type", "add", "pos", "5", ProjectName, P5)
	} else if x < -HALF_MIN && x > -HALF_MAX && y > HALF_MIN && y < HALF_MAX {
		// Position 6
		level.Info(util.Logger).Log("id", id, "type", "add", "pos", "6", ProjectName, P6)
	} else if x < -ON_MIN && x > -ON_MAX && y > OFF_MIN && y < ON_MIN {
		// Position 7
		level.Info(util.Logger).Log("id", id, "type", "add", "pos", "7", ProjectName, P7)
	} else if x < -HALF_MIN && x > -HALF_MAX && y < -HALF_MIN && y > -HALF_MAX {
		// Position 8
		level.Info(util.Logger).Log("id", id, "type", "add", "pos", "8", ProjectName, P8)
	}
}
