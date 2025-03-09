package dictionary

var dict_hu = Dictionary{
    Meta: Meta{
        CountryCode:                "hu",
    },
    Page: Page{
        BaseHome:                   "Kezdőlap",
        BasePosts:                  "Bejegyzések",
        BaseDonate:                 "Támogatás",
    },
    Auth: Auth{
        Login:                      "Bejelentkezés",
        Register:                   "Regisztráció",
        Logout:                     "Kijelentkezés",
        ForgottenPassword:          "Elfelejtett jelszó",
        UsernameOrEmail:            "Felhasználónév vagy Email-cím",
        Password:                   "Jelszó",
        PasswordAgain:              "Jelszó mégegyszer",
        Username:                   "Felhasználónév",
        Name:                       "Megnevezés",
        Email:                      "Email-cím",
        Phone:                      "Telefonszám",
        AlreadyHaveAnAccount:       "Már van fiókod?:",
        RegErrUsernameExists:       "A felhasználónév már létezik!",
        RegErrEmailExists:          "Az email-cím már regisztrált!",
        RegErrUsernameMinLen:       "A felhasználónévnek minimum {} karakter hosszúnak kell lennie!",
        RegErrUsernameCantCont:     "A felhasználónév nem tartalmazhat: ",
        RegErrEmailValidation:      "Rossz email cím!",
        RegErrPasswordValidation:   "A jelszó nem megfelelő!",
        RegErrPasswordDoNotMatch:   "A megadott két jelszó helytelen!",
        LoginErrBadUsernameOrEmail: "Rossz felhasználónév vagy email-cím!",
        LoginErrBadPassword:        "Rossz jelszó!",
    },
}
