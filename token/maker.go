package token

import "time"

// managin token makers
type Maker interface {
	CreateToken(userid int64, security_key string, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}
