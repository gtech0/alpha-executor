package model

import (
	"flag"
	"gopkg.in/ini.v1"
	"log"
)

type Config struct {
	TestCount int
	Source    string
	Tests     string
	Output    string
}

func GetConfig() (*Config, error) {
	config := flag.Lookup("config-path").Value.String()

	inidata, err := ini.Load(config)
	if err != nil {
		log.Fatalf("Fail to read file: %v", err)
	}

	section := inidata.Section("")

	tests, err := section.Key("tests").Int()
	if err != nil {
		log.Fatal(err)
	}

	data := &Config{
		TestCount: tests,
		Source:    section.Key("source").String(),
		Tests:     section.Key("tests_dir").String(),
		Output:    section.Key("output_dir").String(),
	}
	return data, err
}
