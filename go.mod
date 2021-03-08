module github.com/y-yagi/tsudura

go 1.15

require (
	github.com/avast/retry-go v3.0.0+incompatible
	github.com/aws/aws-sdk-go v1.37.26
	github.com/dgraph-io/badger v1.6.2
	github.com/dgraph-io/badger/v2 v2.2007.2 // indirect
	github.com/dgraph-io/ristretto v0.0.3 // indirect
	github.com/fsnotify/fsnotify v1.4.9
	github.com/getlantern/golog v0.0.0-20201105130739-9586b8bde3a9 // indirect
	github.com/getlantern/hidden v0.0.0-20201229170000-e66e7f878730 // indirect
	github.com/getlantern/ops v0.0.0-20200403153110-8476b16edcd6 // indirect
	github.com/getlantern/systray v1.1.0
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/lunixbochs/vtclean v1.0.0 // indirect
	github.com/manifoldco/promptui v0.8.0
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/y-yagi/configure v0.2.0
	github.com/y-yagi/goext v0.6.0
	github.com/y-yagi/rnotify v0.1.0
	golang.org/x/net v0.0.0-20210226172049-e18ecbb05110 // indirect
	golang.org/x/sys v0.0.0-20210308170721-88b6017d0656
	google.golang.org/protobuf v1.25.0 // indirect
)

replace github.com/fsnotify/fsnotify => github.com/y-yagi/fsnotify v1.4.10-0.20201227062311-078207fcf401
