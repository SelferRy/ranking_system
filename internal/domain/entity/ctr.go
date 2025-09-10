package entity

import (
	"fmt"
)

type CTR struct {
	value float64
}

func NewCTR(impressions, clicks int64) (CTR, error) {
	if impressions < 0 || clicks < 0 {
		return CTR{}, fmt.Errorf("impressions and clicks must be >= 0")
	}
	var val float64
	if impressions == 0 {
		val = 0.0
	} else {
		val = float64(clicks) / float64(impressions)
	}
	if val < 0 || val > 1 {
		return CTR{}, fmt.Errorf("value must be between 0 and 1")
	}
	return CTR{val}, nil
}

func (c CTR) Value() float64 {
	return c.value
}

// WithImpression increment CTR via impression in FP style
func (c CTR) WithImpression(impressions, clicks int) CTR {
	val := float64(clicks) / float64(impressions+1)
	return CTR{val}
}

// WithClicks increment CTR via click in FP style
func (c CTR) WithClicks(clicks int) CTR {
	val := float64(clicks+1) / float64(clicks)
	return CTR{val}
}
