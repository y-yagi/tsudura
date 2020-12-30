module github.com/y-yagi/tsudura

go 1.15

require (
	github.com/avast/retry-go v3.0.0+incompatible
	github.com/aws/aws-sdk-go v1.36.27
	github.com/dgraph-io/badger v1.6.2
	github.com/dgraph-io/badger/v2 v2.2007.2 // indirect
	github.com/fsnotify/fsnotify v1.4.9
	github.com/getlantern/systray v1.1.0
	github.com/manifoldco/promptui v0.8.0
	github.com/y-yagi/configure v0.2.0
	github.com/y-yagi/goext v0.6.0
	github.com/y-yagi/rnotify v0.0.0-20201227081429-1fdaf9c1f914
	golang.org/x/sys v0.0.0-20201231184435-2d18734c6014
)

replace github.com/fsnotify/fsnotify => github.com/y-yagi/fsnotify v1.4.10-0.20201227062311-078207fcf401
