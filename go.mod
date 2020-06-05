module github.com/platinasystems/goes-build

require (
	github.com/platinasystems/go-cpio v0.0.1
	github.com/platinasystems/jobserver v0.0.0-00010101000000-000000000000
)

replace github.com/platinasystems/jobserver => ../jobserver

go 1.14
