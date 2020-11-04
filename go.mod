module github.com/platinasystems/goes-build

require (
	github.com/petermattis/goid v0.0.0-20180202154549-b0b1615b78e5 // indirect
	github.com/platinasystems/go-cpio v0.0.1
	github.com/platinasystems/jobserver v0.0.0-00010101000000-000000000000
	github.com/sasha-s/go-deadlock v0.2.0 // indirect
)

replace github.com/platinasystems/jobserver => ../jobserver

go 1.14
