package session

type Session interface {
	Set(key, value interface{}) error
	Get(Key interface{}) interface{}
	Delete(key interface{}) error
	SessionID() string
}
