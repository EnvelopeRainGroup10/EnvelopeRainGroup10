package model

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type User struct {
	ID    int64
	Count int64
}

type Envelope struct {
	ID         int64 `json:"envelope_id"`
	UID        int64 `json:"uid"`
	Opened     bool  `json:"opened"`
	Value      int64 `json:"value"`
	SnatchTime int64 `json:"snatch_time"`
}
