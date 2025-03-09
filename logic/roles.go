package logic

type roles_t struct {
    USER        string
    MODERATOR   string
    ADMIN       string
    EDITOR      string
}

var ROLES = roles_t {
    USER:       "USER",
    MODERATOR:  "MODERATOR",
    ADMIN:      "ADMIN",
    EDITOR:     "EDITOR",
}
