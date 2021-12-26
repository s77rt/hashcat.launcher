package hashcatlauncher

import (
	"log"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

func (a *App) NewWatcher() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	lazyWatcher := func(c <-chan bool, callback func()) {
		call := false
		for {
			select {
			case <-c:
				call = true
			case <-time.After(1 * time.Second):
				if call {
					call = false
					callback()
				}
			}
		}
	}

	watcherHashcatChan := make(chan bool)
	watcherHashesChan := make(chan bool)
	watcherDictionariesChan := make(chan bool)
	watcherRulesChan := make(chan bool)
	watcherMasksChan := make(chan bool)

	go lazyWatcher(watcherHashcatChan, a.WatcherHashcatCallback)
	go lazyWatcher(watcherHashesChan, a.WatcherHashesCallback)
	go lazyWatcher(watcherDictionariesChan, a.WatcherDictionariesCallback)
	go lazyWatcher(watcherRulesChan, a.WatcherRulesCallback)
	go lazyWatcher(watcherMasksChan, a.WatcherMasksCallback)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Create != fsnotify.Create && event.Op&fsnotify.Remove != fsnotify.Remove && event.Op&fsnotify.Rename != fsnotify.Rename {
					continue
				}
				if event.Name == a.Hashcat.BinaryFile {
					watcherHashcatChan <- true
				} else {
					dir, _ := filepath.Split(event.Name)
					dir = filepath.Join(dir) // to be compatible with below directories as they are all constructed by filepath (and to avoid trailing slash issue)
					switch dir {
					case a.HashesDir:
						watcherHashesChan <- true
					case a.DictionariesDir:
						watcherDictionariesChan <- true
					case a.RulesDir:
						watcherRulesChan <- true
					case a.MasksDir:
						watcherMasksChan <- true
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
