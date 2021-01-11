package web

import (
	"html/template"
	"net/http"
)

const (
	TITLE = "CoolConf Admin"
)

const INDEX = `<!DOCTYPE html>
<html lang="en">
<head>
<title>{{.Title}}</title>
<meta charset="utf-8">
<meta name="description" content="">
<meta name="format-detection" content="telephone=no">
<meta name="msapplication-tap-highlight" content="no">
<meta name="version" content="">
<meta name="viewport" content="user-scalable=no, initial-scale=1, maximum-scale=3, minimum-scale=1, width=device-width, viewport-fit=cover">
<link rel="icon" type="image/png" href="logo.png">
<link rel="icon" type="image/png" sizes="16x16" href="statics/icons/favicon-16x16.png">
<link rel="icon" type="image/png" sizes="32x32" href="statics/icons/favicon-32x32.png">
<link rel="icon" type="image/png" sizes="96x96" href="statics/icons/favicon-96x96.png">
<link rel="icon" type="image/ico" href="statics/icons/favicon.ico">
</head>
<body>
hello!
</body>
</html>`

type Router struct {
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) Index() func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		data := Data{Title: TITLE}
		t, _ := template.New("index").Parse(INDEX)
		t.Execute(res, data)
	}
}

func (r *Router) Upload() func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")
		if req.Method != "POST" {
			res.WriteHeader(http.StatusMethodNotAllowed)
			res.Write([]byte(`{"message": "Method not allowed"}`))
			return
		}
		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(`{"data": "POST ok"}`))
	}
}
