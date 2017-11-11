package mongo

import (
	"github.com/impactasaurus/server/data"
	"gopkg.in/mgo.v2"
)

type jwtKV struct {
	JTI string `bson:"_id"`
	JWT string `bson:"jwt"`
}

func (m *mongo) SaveJWT(jti string, jwt string) error {
	col, closer := m.getJWTCollection()
	defer closer()

	if err := col.Insert(jwtKV{
		JTI: jti,
		JWT: jwt,
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
