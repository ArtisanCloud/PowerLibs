test: test-cache

test-cache:
	go test -v cache/memory.go cache/memory_test.go
	go test -v cache/errors.go cache/redis.go cache/redis_test.go

