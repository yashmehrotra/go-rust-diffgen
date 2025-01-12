rustdiffgen:
	cd lib/diffgen && cargo build --release

build:
	go build main.go

benchmark:
	go test -bench=. -benchmem
