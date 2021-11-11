package auth

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/game-sales-analytics/users-service/internal/db/repository"
	"github.com/game-sales-analytics/users-service/internal/passhash"
)

type LoginDefaultCreds struct {
	UserIPAddress       string
	UserDeviceUserAgent string
}

type LoginWithEmailCreds struct {
	LoginDefaultCreds
	Email    string
	Password string
}

func (a authsrv) LoginWithEmail(ctx context.Context, creds LoginWithEmailCreds) (LoginResult, error) {
	user, err := a.repo.GetUserLoginInfo(ctx, creds.Email)
	if nil != err {
		if errors.Is(err, repository.ErrUserNotExists) {
			randSpan := rand.Int63n(251)
			time.Sleep(time.Millisecond * time.Duration((4736 + randSpan)))
			return LoginResult{}, errors.New("unauthorized")
		}

		return LoginResult{}, errors.New("internal")
	}

	matched, err := passhash.Verify(creds.Password, user.Password)
	if nil != err {
		return LoginResult{}, errors.New("internal")
	}
	if !matched {
		return LoginResult{}, errors.New("unauthorized")
	}

	token, err := generateToken(user.ID, a.cfg.Key)
	if nil != err {
		return LoginResult{}, errors.New("internal")
	}

	loginRecord := repository.NewUserLoginToSave{
		ID:                  "",
		UserID:              user.ID,
		LoggedInAt:          time.Now(),
		UserIPAddress:       creds.UserIPAddress,
		UserDeviceUserAgent: creds.UserDeviceUserAgent,
	}
	if err := a.repo.SaveNewUserLogin(ctx, loginRecord); nil != err {
		return LoginResult{}, err
	}

	return LoginResult{
		Token: token,
	}, nil
}