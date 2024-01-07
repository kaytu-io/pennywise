package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog/log"
	"strings"
)

var (
	delimeter      = "."
	seperator      = "_"
	tagName        = "koanf"
	uptemplate     = "************************ loaded configuration ***************************"
	bottomTemplate = "*************************************************************************"
)

func load(print bool) *Config {
	k := koanf.New(delimeter)
	if err := loadEnv(k); err != nil {
		log.Printf("error loading environment variable : %v", err)
	}
	//This variable is set to make sure the config file is there
	configFile := true

	//load YAML config
	if err := k.Load(file.Provider("../config/config.yml"), yaml.Parser()); err != nil {
		configFile = false
		log.Printf("error loading config: %v", err)
	}

	// load JSON config
	if err := k.Load(file.Provider("../config/config.json"), json.Parser()); err != nil {
		if configFile == false {
			log.Printf("error loading config: %v", err)
		}
	}

	config := Config{}
	var tag = koanf.UnmarshalConf{Tag: tagName}
	if err := k.UnmarshalWithConf("", &config, tag); err != nil {
		panic(fmt.Errorf("error unmarshaling config : %v", err))
	}
	if print {
		fmt.Printf("%v \n %v %v\n", uptemplate, spew.Sdump(config), bottomTemplate)
	}
	return &config
}

func loadEnv(k *koanf.Koanf) error {
	callBack := func(source string) string {
		base := strings.ToLower(source)
		return strings.ReplaceAll(base, seperator, delimeter)
	}
	//load environment variable
	if err := k.Load(env.Provider("", delimeter, callBack), nil); err != nil {
		return fmt.Errorf("error loading environment variables : %s", err)
	}
	return nil
}
