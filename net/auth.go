package net

import ()

type Auth interface {
	SendKey() (key *[]byte)
	RecvKey(key *[]byte) (name string, err error)
}
