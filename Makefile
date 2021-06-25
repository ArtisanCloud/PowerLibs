test: test-cache

test-cache:
	go test -v cache/memory.go cache/memory_test.go
