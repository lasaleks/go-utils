package saveevent

import (
	"encoding/json"
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/lasaleks/gormq"
)

type Event struct {
	SystemKey    string           `json:"system_key"`
	TypeEventKey string           `json:"type_event_key"`
	TimeBegin    *string          `json:"time_begin,omitempty"`
	TimeEnd      *string          `json:"time_end,omitempty"`
	ListContext  [][3]interface{} `json:"list_context"`
}

func SaveEvents(out chan<- gormq.MessageAmpq, system_key string, ev Event) error {
	data, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	out <- gormq.MessageAmpq{
		Exchange:     "QPDB",
		Routing_key:  fmt.Sprintf("events.save.%s", system_key),
		Content_type: "text/plain",
		Data:         data,
	}
	return nil
}

func UnixTimeToStr(utime int64) string {
	return time.Unix(utime, 0).Format("2006-01-02 15:04:05 -07")
}

func GetFIO(family string, first_name string, last_name string) string {
	fio := family
	first_symbol, size_symbol := utf8.DecodeRuneInString(first_name)
	if size_symbol > 0 {
		fio += fmt.Sprintf(" %c", first_symbol) + "."
	}
	first_symbol, size_symbol = utf8.DecodeRuneInString(last_name)
	if size_symbol > 0 {
		fio += fmt.Sprintf("%c", first_symbol) + "."
	}
	return fio
}

func GetInfoUnit(type_unit int, tab_number int, family string, first_name string, last_name string, grname string, stname string) string {
	var info string
	switch type_unit {
	case 0:
		info = fmt.Sprintf("#%d %s", tab_number, first_name)
	case 2:
		info = fmt.Sprintf("#%d %s", tab_number, first_name)
	case 1:
		//  ФИО:Затуливетров Р.В.; Долж.:Работник сторонней организации; Подр.:ООО "Майн Радио Системз-Р"
		info = fmt.Sprintf("Таб.№:%d ФИО:%s", tab_number, GetFIO(family, first_name, last_name))
		if len(grname) > 0 {
			info += fmt.Sprintf("Долж:%s ", grname)
		}
		if len(stname) > 0 {
			info += fmt.Sprintf("Подр:%s", stname)
		}
	}
	return info
}
