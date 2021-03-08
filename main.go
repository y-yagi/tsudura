package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/fsnotify/fsnotify"
	"github.com/getlantern/systray"
	"github.com/y-yagi/configure"
	"github.com/y-yagi/rnotify"
	"github.com/y-yagi/tsudura/db"
	"github.com/y-yagi/tsudura/icon"
	log "github.com/y-yagi/tsudura/logger"
	"github.com/y-yagi/tsudura/storage"
	"github.com/y-yagi/tsudura/utils"
)

var (
	watcher       *rnotify.Watcher
	logger        *log.Logger
	cfg           utils.Config
	storageClient *storage.Client
	dbClient      *db.Client
)

const (
	app = "Tsudura"
)

func init() {
	err := configure.Load(app, &cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if len(cfg.Root) != 0 && !strings.HasSuffix(cfg.Root, string(os.PathSeparator)) {
		cfg.Root = cfg.Root + string(os.PathSeparator)
	}
}

func main() {
	var err error

	if len(cfg.Root) == 0 {
		if err = setupConfig(); err != nil {
			fmt.Printf("%v\n", err)
			return
		}
	}

	logger, err = log.NewLogger(app)
	if err != nil {
		fmt.Printf("logger build failed: %v\n", err)
		return
	}

	if storageClient, err = storage.Init(&cfg); err != nil {
		fmt.Printf("client build failed: %v\n", err)
		return
	}

	dir := configure.ConfigDir(app)
	dbClient, err = db.Init(filepath.Join(dir, app+".db"))
	if err != nil {
		fmt.Printf("storage client build failed: %v\n", err)
		return
	}

	watcher, err = rnotify.NewWatcher()
	if err != nil {
		fmt.Printf("watcher build failed: %v\n", err)
		return
	}

	if err = sync(); err != nil {
		fmt.Printf("directory sync failed: %v\n", err)
		return
	}

	systray.Run(onReady, onExit)
	return
}

func onReady() {
	run()
	systray.SetIcon(icon.Data)
	systray.SetTooltip(app)
	mQuit := systray.AddMenuItem("Quit", "Quit")
	<-mQuit.ClickedCh
	systray.Quit()
}

func onExit() {
	watcher.Close()
	logger.Close()
	dbClient.Term()
}

func sync() error {
	serverFiles := map[string]*s3.Object{}
	localFiles := map[string]int64{}

	res, err := storageClient.S3.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(cfg.Bucket),
	})
	if err != nil {
		return err
	}

	for _, object := range res.Contents {
		serverFiles[*object.Key] = object
	}

	err = filepath.Walk(cfg.Root, func(path string, info os.FileInfo, err error) error {
		if path == cfg.Root {
			return nil
		}

		if err != nil {
			return err
		}

		key := strings.Split(path, cfg.Root)[1]
		localFiles[key] = info.ModTime().Unix()
		return nil
	})

	for key, _ := range localFiles {
		if object, found := serverFiles[key]; found {
			delete(serverFiles, key)
			localFileEtag, err := dbClient.Get([]byte(key))
			// TODO(y-yagi) need to update file when etag changed.
			if err != nil && *object.ETag == string(localFileEtag) {
				delete(localFiles, key)
				continue
			}
		}
	}

	for key, _ := range localFiles {
		// TODO: consider conflict
		upload(filepath.Join(cfg.Root, key))
	}

	for key, serverFile := range serverFiles {
		localFileEtag, err := dbClient.Get([]byte(key))
		if err != nil {
			if err = download(filepath.Join(cfg.Root, key), *serverFile.ETag); err != nil {
				msg := fmt.Sprintf("download error %v\n", err)
				logger.Error(msg)
			}
		} else {
			// TODO: consider conflict
			if string(localFileEtag) == *serverFile.ETag {
				if err = destroy(filepath.Join(cfg.Root, key)); err != nil {
					msg := fmt.Sprintf("delete error %v\n", err)
					logger.Error(msg)
				}
			}
		}
	}

	return nil
}

func run() {
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if cfg.Debug {
					fmt.Printf("DEBUG: %v\n", event)
				}

				if event.Op&fsnotify.Remove == fsnotify.Remove || event.Op&fsnotify.Rename == fsnotify.Rename {
					if err := destroy(event.Name); err != nil {
						msg := fmt.Sprintf("delete error %v\n", err)
						logger.Error(msg)
					}
				} else if event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Write == fsnotify.Write {
					if err := upload(event.Name); err != nil {
						msg := fmt.Sprintf("upload error %v\n", err)
						logger.Error(msg)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				msg := fmt.Sprintf("watch error %v\n", err)
				logger.Error(msg)
			}
		}
	}()

	if err := watcher.Add(cfg.Root); err != nil {
		msg := fmt.Sprintf("watcher couldn't add path %v\n", err)
		logger.Error(msg)
	}
}

func upload(path string) error {
	result, err := storageClient.Upload(path)
	if err != nil {
		return err
	}

	return dbClient.Set([]byte(result.Key), []byte(result.ETag))
}

func download(path string, etag string) error {
	result, err := storageClient.Download(path, etag)
	if err != nil {
		return err
	}

	return dbClient.Set([]byte(result.Key), []byte(result.ETag))
}

func destroy(path string) error {
	if cfg.AddOnly {
		return nil
	}

	result, err := storageClient.Destroy(path)
	if err != nil {
		return err
	}

	return dbClient.Delete([]byte(result.Key))
}
