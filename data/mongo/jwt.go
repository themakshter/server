package mongo

import (
	"time"

	"github.com/impactasaurus/server/auth"
	"github.com/impactasaurus/server/data"
	"gopkg.in/mgo.v2"
)

type jwtKV struct {
	JTI     string    `bson:"_id"`
	JWT     string    `bson:"jwt"`
	User    string    `bson:"user"`
	Created time.Time `bson:"created"`
}

func (m *mongo) SaveJWT(jti, jwt string, u auth.User) error {
	col, closer := m.getJWTCollection()
	defer closer()

	if err := col.Insert(jwtKV{
		JTI:     jti,
		JWT:     jwt,
		Created: time.Now(),
		User:    u.UserID(),
	}); err != nil {
		return err
	}

	return nil
}
func (m *mongo) GetJWT(jti string) (string, error) {
	col, closer := m.getJWTCollection()
	defer closer()

	kv := jwtKV{}
	err := col.FindId(jti).One(&kv)
	if err != nil {
		if mgo.ErrNotFound == err {
			return "", data.NewNotFoundError("JWT")
		}
		return "", err
	}
	return kv.JWT, nil
}
