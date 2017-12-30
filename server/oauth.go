package server

import (
	"net/http"
	"net/url"

	"github.com/flosch/pongo2"
)

// TODO: make redirect_uri a configurable setting

func (s *Server) oauthBegin(w http.ResponseWriter, r *http.Request) {
	q := url.Values{}
	q.Set("client_id", s.clientID)
	q.Set("scope", "chat:write:user")
	q.Set("redirect_uri", "https://slouch.quickmediasolutions.com/oauth/return")
	u, err := url.Parse("https://slack.com/oauth/authorize")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	u.RawQuery = q.Encode()
	s.render(w, r, "oauth/begin.html", pongo2.Context{
		"title": "Begin",
		"url":   u,
	})
}

func (s *Server) oauthReturn(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, "oauth/return.html", pongo2.Context{
		"title": "Return",
	})
}
