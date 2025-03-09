package dictionary

import (
	"net/http"
)

type Meta struct {
    CountryCode                 string
}

type Page struct {
    BaseHome                    string
    BasePosts                   string
    BaseDonate                  string
}

type Auth struct {
    Login                       string
    Register                    string
    Logout                      string
    ForgottenPassword           string
    UsernameOrEmail             string
    Password                    string
    PasswordAgain               string
    Username                    string
    Name                        string
    Email                       string
    Phone                       string
    AlreadyHaveAnAccount        string
    RegErrUsernameExists        string
    RegErrEmailExists           string
    RegErrUsernameMinLen        string
    RegErrUsernameCantCont      string
    RegErrEmailValidation       string
    RegErrPasswordValidation    string
    RegErrPasswordDoNotMatch    string
    LoginErrBadUsernameOrEmail  string
    LoginErrBadPassword         string
}

type Dictionary struct {
    Meta        Meta
    Page        Page
    Auth        Auth
}

type DictCollection struct {
    Hungarian   Dictionary
    English     Dictionary
}

var Dict = DictCollection{
    Hungarian: dict_hu,
    English: dict_en,
}

var langCookieName = "lang"

func SetLanguage(w http.ResponseWriter, val string) {
    http.SetCookie(w, &http.Cookie{Name: langCookieName, Value: val, Path: "/", HttpOnly: false})
}

func GetLanguage(r *http.Request) Dictionary {
    cookie, err := r.Cookie(langCookieName)
    lang := ""

    if nil != err {
        lang = "hu"
    } else {
        lang = cookie.Value
    }

    switch lang {
    case Dict.Hungarian.Meta.CountryCode:
        return Dict.Hungarian
    case Dict.English.Meta.CountryCode:
        return Dict.English
    default:
        return Dict.Hungarian
    }
}
