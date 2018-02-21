package accessor

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/db/lock"
	jwt "github.com/dgrijalva/jwt-go"
)

//go:generate counterfeiter . AccessorFactory
type AccessorFactory interface {
	Create(*http.Request) Accessor
}

type accessorFactory struct {
	conn        db.Conn
	lockFactory lock.LockFactory
	publicKey   *rsa.PublicKey
	logger      lager.Logger
}

func NewAccessorFactory(conn db.Conn, lockFactory lock.LockFactory, publicKey *rsa.PublicKey, logger lager.Logger) AccessorFactory {
	return &accessorFactory{
		conn:        conn,
		lockFactory: lockFactory,
		publicKey:   publicKey,
		logger:      logger,
	}
}

func (f *accessorFactory) Create(r *http.Request) Accessor {
	return &accessor{
		db.NewTeamFactory(f.conn, f.lockFactory),
		f.teamsFromRequest(r),
		f.isAdminFromRequest(r),
		f.logger,
	}
}

func (f *accessorFactory) teamsFromRequest(r *http.Request) []string {

	token, err := getJWT(r, f.publicKey)
	if err != nil {
		return []string{}
	}

	claims := token.Claims.(jwt.MapClaims)
	team, ok := claims["teamName"]

	if ok {
		return []string{team.(string)}
	} else {
		return []string{}
	}
}

func (f *accessorFactory) isAdminFromRequest(r *http.Request) bool {

	token, err := getJWT(r, f.publicKey)
	if err != nil {
		return false
	}

	claims := token.Claims.(jwt.MapClaims)
	isAdmin, ok := claims["isAdmin"]

	return ok && isAdmin.(bool)
}

func getJWT(r *http.Request, publicKey *rsa.PublicKey) (token *jwt.Token, err error) {
	fun := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	}

	if ah := r.Header.Get("Authorization"); ah != "" {
		if len(ah) > 6 && strings.ToUpper(ah[0:6]) == "BEARER" {
			return jwt.Parse(ah[7:], fun)
		}
	}

	return nil, errors.New("unable to parse authorization header")
}
