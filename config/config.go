package config

import (
	"asapgiri/golib/logger"
	"asapgiri/golib/session"
	"encoding/json"
	"os"
	"path/filepath"
)

var log = logger.Logger {
    Color: logger.Colors.Yellow,
    Pretext: "config",
}

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
    MinUsernameLen      int
    MinPasswordLen      int
    NameCantContain     []string
}

type DonationConfig struct {
    Merchant            string
    SecretKey           string
    SimplePayURL        string
    SimplePayReturnURL  string
}

type ConfigT struct {
    Http        HttpConfig
    Dbase       DbConfig
    User        UserConfig
    Site        session.Config
    Donation    DonationConfig
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
        MinUsernameLen: 3,
        MinPasswordLen: 5,
        NameCantContain: []string{},
    },
    Site: session.Config{
        Title: "Dunakéke",
        SiteTitle: "Dunakéke",
        TitleSeparator: " - ",
        MaxImgUploadMB: 10,
    },
    Donation: DonationConfig{
        Merchant: "",
        SecretKey: "",
        SimplePayURL: "https://simplepay.hu",
        SimplePayReturnURL: "",
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
            log.Println("Check your `.config.json` format!")
            panic(err)
        }
    }
}
