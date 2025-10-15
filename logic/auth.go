package logic

import (
	"asapgiri/golib/session"
	"dunakeke/config"
	"dunakeke/dbase"
	"dunakeke/dictionary"
	"errors"
	"net/mail"
	"slices"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func Authenticate(a *session.Auth) {
    // TODO: Check if user can be mocked out to existing and be used for unauthenticated login...
    if a.Username != "" {
        user := User{}
        user.FindByUsername(a.Username)
        if "" == user.Username {
            *a = session.Auth{}
            return
        }

        a.Id = user.Id
        a.Username = user.Username
        a.Name = user.Name
        a.Email = user.Email
        a.Roles = user.Roles
        a.Error = ""
        a.IsAdmin = slices.Contains(user.Roles, ROLES.ADMIN)
        a.IsMod = a.IsAdmin || slices.Contains(user.Roles, ROLES.MODERATOR)
    }
}

func (user *User) Register(dict dictionary.Dictionary, password_clear_a string, password_clear_b string) error {
    new_user := dbase.User{}

    if new_user.FindByUsername(user.Username) == nil {
        return errors.New(dict.Auth.RegErrUsernameExists)
    }
    if new_user.FindByEmail(user.Email) == nil {
        return errors.New(dict.Auth.RegErrEmailExists)
    }

    if len(user.Username) < config.Config.User.MinPasswordLen {
        return errors.New(strings.Replace(dict.Auth.RegErrUsernameMinLen, "{}", string(config.Config.User.MinPasswordLen), 1))
    }
    for _, bword := range(config.Config.User.NameCantContain) {
        if strings.Contains(user.Username, bword) {
            return errors.New(dict.Auth.RegErrUsernameCantCont + bword)
        }
    }

    _, err := mail.ParseAddress(user.Email)
    if nil != err {
        return errors.New(dict.Auth.RegErrEmailValidation)
    }
    if len(password_clear_a) < 6 {
        return errors.New(dict.Auth.RegErrPasswordValidation)
    }
    // FIXME: Also validate on the web site...
    if password_clear_a != password_clear_b {
        return errors.New(dict.Auth.RegErrPasswordDoNotMatch)
    }

    pwh, _ := bcrypt.GenerateFromPassword([]byte(password_clear_a), 0)
    new_user = user.UnMap()
    new_user.Id = primitive.NewObjectID()
    new_user.PasswordHash = string(pwh)
    new_user.Roles = []string{ROLES.USER}
    new_user.RegDate = time.Now()
    new_user.EditDate = time.Now()

    new_user.Add()
    log.Printf("Registerd with %s:%s\n", new_user.Id, string(pwh))

    addExistingDonationsToNewUser(new_user)

    return nil
}

func (user *User) Login(dict dictionary.Dictionary, uname_or_email string, password_clear string) error {
    duser := dbase.User{}
    duser_uname := dbase.User{}
    duser_email := dbase.User{}
    err_uname := duser_uname.FindByUsername(uname_or_email)
    err_email := duser_email.FindByEmail(uname_or_email)

    if nil != err_uname && nil != err_email {
        return errors.New(dict.Auth.LoginErrBadUsernameOrEmail)
    }

    if nil == err_uname {
        duser = duser_uname
    } else {
        duser = duser_email
    }

    log.Println(duser)
    if nil != bcrypt.CompareHashAndPassword([]byte(duser.PasswordHash), []byte(password_clear)) {
        return errors.New(dict.Auth.LoginErrBadPassword)
    }

    user.Map(duser)

    return nil
}

func (user *User) Logout() {
    // nothing to do here...
}

func (user *User) List() []User {
    duser := dbase.User{}
    dusers, _ := duser.List()

    users := make([]User, len(dusers))
    for i, u := range(dusers) {
        users[i].Map(u)
    }

    return users
}
