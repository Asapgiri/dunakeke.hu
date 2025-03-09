package main

import (
	"dunakeke/apis"
	"dunakeke/pages"
	"net/http"
)

func setup_routes() {
    http.HandleFunc("GET /",                    pages.Root)
    http.HandleFunc("GET /index",               pages.Root)
    http.HandleFunc("GET /index.html",          pages.Root)

    http.HandleFunc("GET /api/lang/{lang}",     apis.SetLanguage)

    http.HandleFunc("GET /login",               pages.Login)
    http.HandleFunc("POST /login",              pages.Login)
    http.HandleFunc("GET /register",            pages.Register)
    http.HandleFunc("POST /register",           pages.Register)
    http.HandleFunc("GET /logout",              pages.Logout)

    http.HandleFunc("GET /admin",               pages.AdminPage)
    http.HandleFunc("GET /admin/links",         pages.NotFound)
    http.HandleFunc("GET /admin/users",         pages.AdminUsers)
    http.HandleFunc("GET /admin/posts",         pages.NotFound)
    http.HandleFunc("GET /admin/comments",      pages.NotFound)
    http.HandleFunc("GET /admin/newsletter",    pages.NotFound)
    http.HandleFunc("GET /admin/statistics",    pages.NotFound)
    http.HandleFunc("GET /admin/donations",     pages.NotFound)

    http.HandleFunc("GET /admin/api/addlink",   pages.NotFound)
    http.HandleFunc("GET /admin/api/dellink",   pages.NotFound)

    http.HandleFunc("GET /donate",              pages.NotFound)

    http.HandleFunc("GET /user/{id}",           pages.NotFound)
    http.HandleFunc("GET /user/edit/{id}",      pages.NotFound)
    http.HandleFunc("GET /user/delete/{id}",    pages.NotFound)

    http.HandleFunc("GET /post/{id}",           pages.PostShow)
    http.HandleFunc("GET /post/edit/{id}",      pages.PostEdit)
    http.HandleFunc("GET /post/delete/{id}",    pages.PostDelete)

    http.HandleFunc("POST /api/comment/add",    pages.NotFound)
    http.HandleFunc("POST /api/comment/edit",   pages.NotFound)
    http.HandleFunc("POST /api/comment/delete", pages.NotFound)

    http.HandleFunc("GET /api/newsletter/register",     pages.NotFound)
    http.HandleFunc("GET /api/newsletter/unregister",   pages.NotFound)

    http.HandleFunc("GET /access-violation",    pages.AccessViolation)
}
