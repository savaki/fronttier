package sessions

type memory struct {
	repo map[string]*Session
}

func Memory() Store {
	repo := make(map[string]*Session)
	return &memory{
		repo: repo,
	}
}

func (self *memory) Create(values map[string]string) (*Session, error) {
	session := create(values)
	err := self.Set(session.Id, session)
	return session, err
}

func (self *memory) Set(key string, session *Session) error {
	self.repo[key] = session
	return nil
}

func (self *memory) Get(key string) (*Session, error) {
	session := self.repo[key]
	if session == nil {
		return nil, nil
	} else {
		return session, nil
	}
}

func (self *memory) Delete(key string) error {
	delete(self.repo, key)
	return nil
}
