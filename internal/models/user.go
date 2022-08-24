package models

import "github.com/dgrijalva/jwt-go"

type UserFromRequest struct {
	ID          int64  `json:"id"`
	Username    string `json:"username"`
	AccessToken string `json:"-"`
}

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

type UserFullName struct {
	ID       int64   `json:"id"`
	FullName *string `json:"fullName"`
}

// Claims описывает кодинг для jwt
type Claims struct {
	jwt.StandardClaims
	UserName string `json:"username"`
	ID       int64  `json:"id"` // id пользователя
	UID      string `json:"uid"`
}

// AccessOrgs описывает доступные ИД организаций
type AccessOrgs struct {
	UserID     int64
	FactoryIDs []int
	SupplyIDs  []int
}
