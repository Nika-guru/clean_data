package auth

import (
	"base/pkg/log"
	"fmt"
	"strconv"
	"strings"
)

type Payload struct {
	UserId    int `json:"userId"`
	AccountId int `json:"accountId"`
}

func (payload *Payload) GetUserIdAndAcountIdFromClaims(claims string) error {
	var err error
	userIdString := strings.Split(claims, ",")[0]
	accountIdString := strings.Split(claims, ",")[1]
	payload.UserId, err = strconv.Atoi(strings.Split(userIdString, ":")[1])
	if err != nil {
		log.Println(log.LogLevelDebug, "GetUserIdAndAcountIdFromClaims", err)
		return err
	}
	payload.AccountId, err = strconv.Atoi(strings.Split(accountIdString, ":")[1])
	if err != nil {
		log.Println(log.LogLevelDebug, "GetUserIdAndAcountIdFromClaims", err)
		return err
	}
	return nil
}

func (payload *Payload) CreatePayloadForJWT() string {
	return fmt.Sprintf("userId:%s,accountId:%s", strconv.Itoa(payload.UserId), strconv.Itoa(payload.AccountId))

}

func (payload *Payload) GetTokenData() (any, error) {
	token, err := GetJWTToken(payload.CreatePayloadForJWT())
	if err != nil {
		return nil, err
	}
	tokenData := struct {
		Token string `json:"token"`
		Type  string `json:"type"`
	}{
		Token: token,
		Type:  "Bearer",
	}

	return tokenData, nil
}
