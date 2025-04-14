package common

func ListenChannel[T any](channel chan T, cb func(value T)) {
	go func() {
		for value := range channel {
			cb(value)
		}
	}()
}
