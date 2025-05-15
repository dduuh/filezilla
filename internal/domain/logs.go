package domain

const (
	USER  = "USER"
	ADMIN = "ADMIN"

	REGISTER = "REGISTER"
	LOGIN    = "LOGIN"
)

type Log struct {
	Id     string `json:"-"       bson:"_id,omitempty"`
	UserId int    `json:"user_id" bson:"user_id"`
	Role   string `json:"role"    bson:"role"`
	Action string `json:"action"  bson:"action"`
}
