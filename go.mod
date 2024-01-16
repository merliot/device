module github.com/merliot/device

go 1.21.5

replace tinygo.org/x/drivers => tinygo.org/x/drivers v0.26.1-0.20231206190939-3fabdc5c9680

require (
	github.com/merliot/dean v0.0.0-20240115013534-487d0167866b
	github.com/merliot/prime v0.0.0-20240116033225-aa3d019fb6be
	github.com/merliot/target v0.0.0-20240113233253-7bc2b49a202d
	github.com/merliot/uf2 v0.0.0-20231228035705-76e82a789f10
)

require (
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	golang.org/x/crypto v0.16.0 // indirect
	golang.org/x/net v0.19.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	tinygo.org/x/drivers v0.0.0-00010101000000-000000000000 // indirect
)
