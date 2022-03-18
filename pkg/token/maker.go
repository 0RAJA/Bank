package token

import "time"

type Maker interface {
	// CreateToken 生成Token
	CreateToken(username string, expireDate time.Duration) (string, error)
	// VerifyToken 解析Token
	VerifyToken(token string) (*Payload, error)
}
