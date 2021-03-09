package model

import "gorm.io/gorm"

type Vehicle struct {
	gorm.Model
	VNo    string `gorm:"unique" json:"vno"`
	Modal  string `json:"modal"`
	Brand  string `json:"brand"`
	Fuel   string  `json:"fuel"`
	Status bool   `json:"status"`
}


// DBMigrate will create and migrate the tables
func DBMigrate(db *gorm.DB) *gorm.DB {
	db.AutoMigrate(&Vehicle{})
	return db
}