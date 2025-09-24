package dictionary

import (
	"net/http"
)

type Meta struct {
    CountryCode                             string
}

type Page struct {
    BaseHome                                string
    BasePosts                               string
    BaseDonate                              string
    NotFound                                string
}

type Auth struct {
    Account                                 string
    AdminSite                               string
    Login                                   string
    Register                                string
    Logout                                  string
    ForgottenPassword                       string
    UsernameOrEmail                         string
    Password                                string
    PasswordAgain                           string
    Username                                string
    Name                                    string
    Email                                   string
    Phone                                   string
    RegDate                                 string
    Roles                                   string
    AlreadyHaveAnAccount                    string
    RegErrUsernameExists                    string
    RegErrEmailExists                       string
    RegErrUsernameMinLen                    string
    RegErrUsernameCantCont                  string
    RegErrEmailValidation                   string
    RegErrPasswordValidation                string
    RegErrPasswordDoNotMatch                string
    LoginErrBadUsernameOrEmail              string
    LoginErrBadPassword                     string
    AccessViolation                         string
}

type Editor struct {
    Description                             string
    TocTitle                                string
    ToolbarUndo                             string
    ToolbarRedo                             string
    ToolbarBold                             string
    ToolbarDel                              string
    ToolbarItalic                           string
    ToolbarQuote                            string
    ToolbarUcwords                          string
    ToolbarUppercase                        string
    ToolbarLowercase                        string
    ToolbarH1                               string
    ToolbarH2                               string
    ToolbarH3                               string
    ToolbarH4                               string
    ToolbarH5                               string
    ToolbarH6                               string
    ToolbarListUl                           string
    ToolbarListOl                           string
    ToolbarHr                               string
    ToolbarLink                             string
    ToolbarReferenceLink                    string
    ToolbarImage                            string
    ToolbarCode                             string
    ToolbarPreformattedText                 string
    ToolbarCodeBlock                        string
    ToolbarTable                            string
    ToolbarDatetime                         string
    ToolbarEmoji                            string
    ToolbarHtmlEntities                     string
    ToolbarPagebreak                        string
    ToolbarWatch                            string
    ToolbarUnwatch                          string
    ToolbarPreview                          string
    ToolbarFullscreen                       string
    ToolbarClear                            string
    ToolbarSearch                           string
    ToolbarHelp                             string
    ToolbarInfo                             string
    ButtonsEnter                            string
    ButtonsCancel                           string
    ButtonsClose                            string
    DialogLinkTitle                         string
    DialogLinkUrl                           string
    DialogLinkUrlTitle                      string
    DialogLinkUrlEmpty                      string
    DialogReferenceLinkTitle                string
    DialogReferenceLinkName                 string
    DialogReferenceLinkUrl                  string
    DialogReferenceLinkUrlId                string
    DialogReferenceLinkUrlTitle             string
    DialogReferenceLinkNameEmpty            string
    DialogReferenceLinkIdEmpty              string
    DialogReferenceLinkUrlEmpty             string
    DialogImageTitle                        string
    DialogImageUrl                          string
    DialogImageLink                         string
    DialogImageAlt                          string
    DialogImageUploadButton                 string
    DialogImageImageURLEmpty                string
    DialogImageUploadFileEmpty              string
    DialogImageFormatNotAllowed             string
    DialogPreformattedTextTitle             string
    DialogPreformattedTextEmptyAlert        string
    DialogPreformattedTextPlaceholder       string
    DialogCodeBlockTitle                    string
    DialogCodeBlockSelectLabel              string
    DialogCodeBlockSelectDefaultText        string
    DialogCodeBlockOtherLanguage            string
    DialogCodeBlockUnselectedLanguageAlert  string
    DialogCodeBlockCodeEmptyAlert           string
    DialogCodeBlockPlaceholder              string
    DialogHtmlEntitiesTitle                 string
    DialogHelpTitle                         string
}

type Admin struct {
    Users                                   string
    Posts                                   string
    Comments                                string
    Links                                   string
    Newsletter                              string
    Statistics                              string
    Donations                               string
    Settings                                string
}

type Donate struct {
    Redirect                                string
    Header                                  string
    Description                             string
    Other                                   string
    Amount                                  string
    Name                                    string
    Email                                   string
    DonateButton                            string
    Newsletter                              string
    GDPRpre                                 string
    GDPRaszf                                string
    GDPRmid                                 string
    GDPRavsz                                string
}

type Dictionary struct {
    Meta        Meta
    Page        Page
    Auth        Auth
    Editor      Editor
    Admin       Admin
    Donate      Donate
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
