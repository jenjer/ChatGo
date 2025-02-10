package iniFunc

import (
	"gopkg.in/ini.v1"
	"fmt"
)

func GetIni(Section, Key string)(string) {
	cfg, err := ini.Load("setting.ini")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	propertySection := cfg.Section(Section)
	Data := propertySection.Key(Key).String()
	return Data
}
