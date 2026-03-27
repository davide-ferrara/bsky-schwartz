package scorer

import (
	"os"
	"sync"
)

var (
	Cache        sync.Map
	cachedReader FileReader
)

type FileReader func(string) ([]byte, error)

func initCachedReader() {
	cachedReader = newCachedReader(os.ReadFile)
}

func GetCachedReader() FileReader {
	return cachedReader
}

func newCachedReader(reader FileReader) FileReader {
	return func(filename string) ([]byte, error) {
		if value, ok := Cache.Load(filename); ok {
			return value.([]byte), nil
		}

		content, err := reader(filename)
		if err != nil {
			return nil, err
		}

		Cache.Store(filename, content)
		return content, nil
	}
}
