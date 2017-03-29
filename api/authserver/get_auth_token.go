package authserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc"
)

const CookieName = "ATC-Authorization"

func (s *Server) GetAuthToken(w http.ResponseWriter, r *http.Request) {
	logger := s.logger.Session("get-auth-token")
	logger.Debug("getting-auth-token")

	var token atc.AuthToken
	teamName := r.FormValue(":team_name")
	teamDB := s.teamDBFactory.GetTeamDB(teamName)
	team, found, err := teamDB.GetTeam()
	if err != nil {
		logger.Error("get-team-by-name", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !found {
		logger.Info("cannot-find-team-by-name", lager.Data{
			"teamName": teamName,
		})
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	tokenType, tokenValue, err := s.tokenGenerator.GenerateToken(time.Now().Add(s.expire), team.Name, team.Admin)
	if err != nil {
		logger.Error("generate-token", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	token.Type = string(tokenType)
	token.Value = string(tokenValue)

	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    fmt.Sprintf("%s %s", token.Type, token.Value),
		Path:     "/",
		Expires:  time.Now().Add(s.expire),
		HttpOnly: s.httpOnly,
		Secure:   s.secure,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}
