package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type HttpConfig struct {
    Url     string
    Port    string
    Cert    string
    Key     string
}

type DbConfig struct {
    Url     string
    Name    string
}

type UserConfig struct {
    MinPasswordLen      int
    NameCantContain     []string
}

type SiteConfig struct {
    Title   string
}

type ConfigT struct {
    Http        HttpConfig
    Dbase       DbConfig
    User        UserConfig
    Site        SiteConfig
}

var Config = ConfigT{
    Http: HttpConfig{
        Url:    "",
        Port:   "3000",
        Cert:   "",
        Key:    "",
    },
    Dbase: DbConfig{
        Url:    "mongodb://localhost:27017",
        Name:   "dunakeke",
    },
    User: UserConfig{
        MinPasswordLen: 8,
        NameCantContain: []string{},
    },
    Site: SiteConfig{
        Title: "Dunakéke",
    },
}

func InitConfig() {
    ex, err := os.Executable()
    if nil != err {
        panic(err)
    }
    expath := filepath.Dir(ex)
    configfile := expath + "/.config.json"

    dat, err := os.ReadFile(configfile)
    if nil != err {
        log.Println(err.Error())
        configdat, _ := json.MarshalIndent(Config, "", "  ")
        os.WriteFile(configfile, configdat, 0644)
    } else {
        err = json.Unmarshal(dat, &Config)
        if nil != err {
            log.Println(err.Error())
        }
    }
}
