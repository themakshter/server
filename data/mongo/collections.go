package mongo

import "gopkg.in/mgo.v2"

type sessionEnder func()

func (m *mongo) getOutcomeCollection() (*mgo.Collection, sessionEnder) {
	session := m.baseSession.Copy()
	return session.DB("").C("outcomesets"), session.Close
}

func (m *mongo) getMeetingCollection() (*mgo.Collection, sessionEnder) {
	session := m.baseSession.Copy()
	return session.DB("").C("meetings"), session.Close
}

func (m *mongo) getOrganisationCollection() (*mgo.Collection, sessionEnder) {
	session := m.baseSession.Copy()
	return session.DB("").C("organisations"), session.Close
}

func (m *mongo) getJWTCollection() (*mgo.Collection, sessionEnder) {
	session := m.baseSession.Copy()
	return session.DB("").C("jwts"), session.Close
}
