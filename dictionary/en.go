package dictionary

var dict_en = Dictionary{
    Meta: Meta{
        CountryCode:                "en",
    },
    Page: Page{
        BaseHome:                   "Home",
        BasePosts:                  "Posts",
        BaseDonate:                 "Donate",
    },
    Auth: Auth{
        Login:                      "Login",
        Register:                   "Register",
        Logout:                     "Logout",
        ForgottenPassword:          "Forgotten Password",
        UsernameOrEmail:            "Username or Email",
        Password:                   "Password",
        PasswordAgain:              "Password again",
        Username:                   "Username",
        Name:                       "Name",
        Email:                      "Email",
        Phone:                      "Phone number",
        AlreadyHaveAnAccount:       "Already have an account?:",
        RegErrUsernameExists:       "Username already exists!",
        RegErrEmailExists:          "Email already used!",
        RegErrUsernameMinLen:       "Username must be minimum {} characters long!",
        RegErrUsernameCantCont:     "Username cant contain word: ",
        RegErrEmailValidation:      "Email validation error!",
        RegErrPasswordValidation:   "Password validation error!",
        RegErrPasswordDoNotMatch:   "Double password doesnt match!",
        LoginErrBadUsernameOrEmail: "Bad username or email!",
        LoginErrBadPassword:        "Bad password!",
    },
}
