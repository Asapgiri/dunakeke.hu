package logic

type roles_t struct {
    USER        string
    MODERATOR   string
    ADMIN       string
}

var roles = roles_t {
    USER:       "USER",
    MODERATOR:  "MODERATOR",
    ADMIN:      "ADMIN",
}
