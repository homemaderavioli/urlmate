package main

type StubStorage struct {
	Key string
	url []byte
}

func (sc StubStorage) CreateURL(key string, value []byte) (string, error) {
	return sc.Key, nil
}

func (sc StubStorage) GetURL(key string) ([]byte, error) {
	return sc.url, nil
}
