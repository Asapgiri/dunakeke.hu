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
    http.HandleFunc("GET /admin/links",         pages.AdminLinks)
    http.HandleFunc("GET /admin/users",         pages.AdminUsers)
    http.HandleFunc("GET /admin/posts",         pages.AdminPosts)
    http.HandleFunc("GET /admin/tags",          pages.AdminTags)
    http.HandleFunc("GET /admin/comments",      pages.NotFound)
    http.HandleFunc("GET /admin/newsletter",    pages.NotFound)
    http.HandleFunc("GET /admin/statistics",    pages.NotFound)
    http.HandleFunc("GET /admin/donations",     pages.AdminDonations)

    http.HandleFunc("POST /admin/links/update",         pages.AdminLinksUpdate)
    http.HandleFunc("GET /admin/links/delete/{id}",     pages.AdminLinksDelete)

    http.HandleFunc("POST /admin/tag/update",           pages.AdminTagsUpdate)
    http.HandleFunc("GET /admin/tag/tl/{id}",           pages.AdminTagsToggleListable)
    http.HandleFunc("GET /admin/tag/delete/{id}",       pages.AdminTagsDelete)

    http.HandleFunc("GET /admin/user/setrole/{id}/{role}",  pages.AdminUserSetRole)

    http.HandleFunc("GET /donate",              pages.DonationRoot)
    http.HandleFunc("POST /donate",             pages.DonationInProgress)
    http.HandleFunc("GET /donate/return",       pages.DonationReturn)
    http.HandleFunc("GET /donate/{id}",         pages.DonationShowStatus)

    http.HandleFunc("GET /user/{id}",           pages.NotFound)
    http.HandleFunc("GET /user/edit/{id}",      pages.NotFound)
    http.HandleFunc("GET /user/delete/{id}",    pages.NotFound)

    http.HandleFunc("GET /post/{id}",           pages.PostShow)
    http.HandleFunc("GET /post/new",            pages.PostNew)
    http.HandleFunc("GET /post/edit/{id}",      pages.PostEdit)
    http.HandleFunc("POST /post/edit/{id}",     pages.PostEditPhotoSave)
    http.HandleFunc("GET /post/delete/{id}",    pages.PostDelete)
    http.HandleFunc("GET /post/pub/{id}/{val}", pages.PostPublish)

    http.HandleFunc("GET /tag/{tagname}",       pages.TagList)
    http.HandleFunc("POST /tag/add",            pages.TagAdd)

    http.HandleFunc("POST /api/post/save",      pages.PostSave)
    http.HandleFunc("POST /api/post/image",     pages.PostSaveImage)

    http.HandleFunc("POST /api/comment/add",    pages.NotFound)
    http.HandleFunc("POST /api/comment/edit",   pages.NotFound)
    http.HandleFunc("POST /api/comment/delete", pages.NotFound)

    http.HandleFunc("GET /api/newsletter/register",     pages.NotFound)
    http.HandleFunc("GET /api/newsletter/unregister",   pages.NotFound)

    http.HandleFunc("GET /access-violation",    pages.AccessViolation)
}
