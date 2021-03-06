package sessions

import (
	"encoding/json"
	redigo "github.com/garyburd/redigo/redis"
)

type redis struct {
	namespace string
	pool      *redigo.Pool
}

type RedisConfig struct {
	namespace   string
	host        string
	maxIdle     int
	connFactory func() (redigo.Conn, error)
}

func Redis() *RedisConfig {
	return &RedisConfig{}
}

func (self *RedisConfig) MaxIdle(maxIdle int) *RedisConfig {
	self.maxIdle = maxIdle
	return self
}

func (self *RedisConfig) Host(host string) *RedisConfig {
	self.host = host
	return self
}

func (self *RedisConfig) ConnFactory(connFactory func() (redigo.Conn, error)) *RedisConfig {
	self.connFactory = connFactory
	return self
}

func (self *RedisConfig) Namespace(namespace string) *RedisConfig {
	self.namespace = namespace
	return self
}

func (self *RedisConfig) getMaxIdle() int {
	maxIdle := self.maxIdle
	if maxIdle < 1 {
		maxIdle = 3
	}
	return maxIdle
}

func (self *RedisConfig) getHost() string {
	host := self.host
	if host == "" {
		host = ":6379"
	}
	return host
}

func (self *RedisConfig) getNamespace() string {
	namespace := self.namespace
	if namespace == "" {
		namespace = "sessions"
	}
	return namespace
}

func (self *RedisConfig) getConnFactory() func() (redigo.Conn, error) {
	connFactory := self.connFactory
	if connFactory == nil {
		host := self.getHost()
		connFactory = func() (redigo.Conn, error) {
			return redigo.Dial("tcp", host)
		}
	}
	return connFactory
}

func (self *RedisConfig) Build() Store {
	pool := redigo.NewPool(self.getConnFactory(), self.getMaxIdle())

	return &redis{
		namespace: self.getNamespace(),
		pool:      pool,
	}
}

func (self *redis) toKey(key string) string {
	return self.namespace + ":" + key
}

func (self *redis) Create(values map[string]string) (*Session, error) {
	session := create(values)
	err := self.Set(session.Id, session)
	return session, err
}

func (self *redis) Set(key string, session *Session) error {
	conn := self.pool.Get()
	defer conn.Close()

	data, err := json.Marshal(session)
	if err != nil {
		return err
	}

	err = conn.Send("SET", self.toKey(key), string(data))
	if err != nil {
		return err
	}

	err = conn.Flush()
	if err != nil {
		return err
	}

	_, err = conn.Receive()
	if err != nil {
		return err
	}

	return nil
}

func (self *redis) Get(key string) (*Session, error) {
	conn := self.pool.Get()
	defer conn.Close()

	err := conn.Send("GET", self.toKey(key))
	if err != nil {
		return nil, err
	}

	err = conn.Flush()
	if err != nil {
		return nil, err
	}

	value, err := conn.Receive()
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, nil
	}

	result := &Session{}
	data := value.([]byte)
	err = json.Unmarshal(data, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (self *redis) Delete(key string) error {
	conn := self.pool.Get()
	defer conn.Close()

	err := conn.Send("DEL", self.toKey(key))
	if err != nil {
		return err
	}

	err = conn.Flush()
	if err != nil {
		return err
	}

	_, err = conn.Receive()
	if err != nil {
		return err
	}

	return nil
}
