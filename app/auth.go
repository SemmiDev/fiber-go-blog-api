package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

type AuthInterface interface {
	CreateAuth(context.Context, uint64, *TokenDetails) error
	FetchAuth(context.Context, string) (uint64, error)
	DeleteRefresh(context.Context, string) error
	DeleteTokens(context.Context, *AccessDetails) error
}

type ClientData struct {
	client *redis.Client
}

var _ AuthInterface = &ClientData{}

func NewAuth(client *redis.Client) *ClientData {
	return &ClientData{client: client}
}

type AccessDetails struct {
	TokenUuid string
	UserId    uint64
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	TokenUuid    string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

// CreateAuth Save token metadata to Redis
func (tk *ClientData) CreateAuth(c context.Context, userid uint64, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	atCreated, err := tk.client.Set(c, td.TokenUuid, strconv.Itoa(int(userid)), at.Sub(now)).Result()
	if err != nil {
		return err
	}
	rtCreated, err := tk.client.Set(c, td.RefreshUuid, strconv.Itoa(int(userid)), rt.Sub(now)).Result()
	if err != nil {
		return err
	}
	if atCreated == "0" || rtCreated == "0" {
		return errors.New("no record inserted")
	}
	return nil
}

// FetchAuth Check the metadata saved
func (tk *ClientData) FetchAuth(c context.Context, tokenUuid string) (uint64, error) {
	userid, err := tk.client.Get(c, tokenUuid).Result()
	if err != nil {
		return 0, err
	}
	userID, _ := strconv.ParseUint(userid, 10, 64)
	return userID, nil
}

// DeleteTokens Once a user row in the token table
func (tk *ClientData) DeleteTokens(c context.Context, authD *AccessDetails) error {
	//get the refresh uuid
	refreshUuid := fmt.Sprintf("%s++%d", authD.TokenUuid, authD.UserId)
	//delete access token
	deletedAt, err := tk.client.Del(c, authD.TokenUuid).Result()
	if err != nil {
		return err
	}
	//delete refresh token
	deletedRt, err := tk.client.Del(c, refreshUuid).Result()
	if err != nil {
		return err
	}
	//When the record is deleted, the return value is 1
	if deletedAt != 1 || deletedRt != 1 {
		return errors.New("something went wrong")
	}
	return nil
}

func (tk *ClientData) DeleteRefresh(c context.Context, refreshUuid string) error {
	//delete refresh token
	deleted, err := tk.client.Del(c, refreshUuid).Result()
	if err != nil || deleted == 0 {
		return err
	}
	return nil
}
