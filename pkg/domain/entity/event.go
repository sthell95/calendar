package entity

import (
	"strconv"
	"strings"
	"time"
)

const ISOLayout = "2006-01-02T15:04:05.000Z"

type Event struct {
	ID          string     `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	Title       string     `json:"title" gorm:"type:varchar"`
	Description string     `json:"description" gorm:"type:varchar"`
	Timezone    string     `json:"timezone" gorm:"type:varchar"`
	Time        *time.Time `json:"time" gorm:"type:timestamp"`
	Duration    duration   `json:"duration" gorm:"type:time"`
	Notes       []string   `json:"notes" gorm:"string[]"`
}

type duration int32

func (d *duration) UnmarshalJSON(b []byte) error {
	dur, _ := strconv.Atoi(strings.Trim(string(b), `"`))
	*d = duration(dur)
	return nil
}
