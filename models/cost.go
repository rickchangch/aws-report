package models

import (
	"time"

	"gorm.io/gorm"
)

// Cost represents the schema of table "costs".
type Cost struct {
	TS      int64
	Service string
	Value   float64
}

type Costs []Cost

var CostModel Cost

// TableName overrides the table name pluralized by the GORM.
func (Cost) TableName() string {
	return "costs"
}

func (c Cost) GetByMultiDateRange(db *gorm.DB, weekDates [][]time.Time) ([]Costs, error) {
	var weekCosts []Costs
	var costs Costs

	for _, weekDate := range weekDates {
		tx := db.
			Where("ts BETWEEN ? AND ?", weekDate[0].Unix(), weekDate[1].Unix()).
			Find(&costs)

		if tx.Error != nil {
			return []Costs{}, tx.Error
		}

		weekCosts = append(weekCosts, costs)
		costs = Costs{}
	}
	return weekCosts, nil
}
