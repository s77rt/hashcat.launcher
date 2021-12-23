package hashcatlauncher

import (
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func (a *App) NewWatcher() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Create != fsnotify.Create && event.Op&fsnotify.Write != fsnotify.Write && event.Op&fsnotify.Remove != fsnotify.Remove && event.Op&fsnotify.Rename != fsnotify.Rename {
					continue
				}
				if event.Name == a.Hashcat.BinaryFile {
					a.WatcherHashcatCallback()
				} else {
					dir, _ := filepath.Split(event.Name)
					dir = filepath.Join(dir) // to be compatible with below directories as they are all constructed by filepath (and to avoid trailing slash issue)
					switch dir {
					case a.HashesDir:
						a.WatcherHashesCallback()
					case a.DictionariesDir:
						a.WatcherDictionariesCallback()
					case a.RulesDir:
						a.WatcherRulesCallback()
					case a.MasksDir:
						a.WatcherMasksCallback()
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("watcher error:", err)
			}
		}
	}()

	if err := watcher.Add(a.HashcatDir); err != nil {
		return err
	}

	if err := watcher.Add(a.HashesDir); err != nil {
		return err
	}

	if err := watcher.Add(a.DictionariesDir); err != nil {
		return err
	}

	if err := watcher.Add(a.RulesDir); err != nil {
		return err
	}

	if err := watcher.Add(a.MasksDir); err != nil {
		return err
	}

	a.Watcher = watcher

	return nil
}
