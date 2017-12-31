package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/flosch/pongo2"
)

// TODO: make this a configurable setting
var redirectUri = "https://slouch.quickmediasolutions.com/oauth/return"

func (s *Server) oauthBegin(w http.ResponseWriter, r *http.Request) {
	q := url.Values{}
	q.Set("client_id", s.clientID)
	q.Set("scope", "chat:write:user")
	q.Set("redirect_uri", redirectUri)
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
	accessToken, err := func() (string, error) {
		code := r.URL.Query().Get("code")
		if len(code) == 0 {
			return "", errors.New("empty verification code")
		}
		q := url.Values{}
		q.Set("client_id", s.clientID)
		q.Set("client_secret", s.clientSecret)
		q.Set("code", code)
		q.Set("redirect_uri", redirectUri)
		u, err := url.Parse("https://slack.com/api/oauth.access")
		if err != nil {
			return "", err
		}
		u.RawQuery = q.Encode()
		resp, err := http.Get(u.String())
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
		var accessResponse struct {
			OK          bool   `json:"ok"`
			Error       string `json:"error"`
			AccessToken string `json:"access_token"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&accessResponse); err != nil {
			return "", err
		}
		if !accessResponse.OK {
			return "", errors.New(accessResponse.Error)
		}
		return accessResponse.AccessToken, nil
	}()
	if err != nil {
		s.render(w, r, "error.html", pongo2.Context{
			"title": "Error",
			"error": err.Error(),
		})
	} else {
		s.render(w, r, "oauth/return.html", pongo2.Context{
			"title":        "Return",
			"access_token": accessToken,
		})
	}
}
