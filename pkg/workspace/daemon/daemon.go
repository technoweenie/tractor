package daemon

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/manifold/qtalk/libmux/mux"
	"github.com/manifold/qtalk/qrpc"
	"github.com/rjeczalik/notify"
)

var logBus = NewMulticastWriteCloser()
var currentServer *exec.Cmd

const addr = "localhost:4242"

func runServer() error {
	bin, err := exec.LookPath("go")
	if err != nil {
		panic(err)
	}
	cmd := exec.Command(bin, "run", "workspace.go")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Stdout = logBus
	cmd.Stderr = logBus
	currentServer = cmd
	return cmd.Run()
}

func extensionIn(path string, exts []string) bool {
	for _, ext := range exts {
		if filepath.Ext(path) == ext {
			return true
		}
	}
	return false
}

func notifyChanges(dir string, exts []string, onlyCreate bool, cb func(path string)) {
	c := make(chan notify.EventInfo, 1)
	types := notify.All
	if onlyCreate {
		types = notify.Create
	}
	if err := notify.Watch(dir, c, types); err != nil {
		log.Fatal(err)
	}
	defer notify.Stop(c)
	for event := range c {
		path := event.Path()
		dir, file := filepath.Split(path)
		if filepath.Base(dir) == ".git" {
			continue
		}
		if filepath.Base(file)[0] == '.' {
			continue
		}
		if extensionIn(path, exts) {
			cb(path)
		}
	}
}

// func watchDelegates(watcher *fsnotify.Watcher) {
// 	files, err := ioutil.ReadDir("./delegates")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	for _, f := range files {
// 		if f.IsDir() {
// 			if err := watcher.Add(path.Join("./delegates", f.Name(), "delegate.go")); err != nil {
// 				log.Fatal(err)
// 			}
// 		}
// 	}
// }

func Run() {
	log.SetOutput(logBus)
	go logBus.WriteTo(os.Stdout)

	go notifyChanges("./...", []string{".go"}, false, func(path string) {
		if currentServer != nil {
			syscall.Kill(-currentServer.Process.Pid, syscall.SIGTERM)
		}
	})

	// watcher, err := fsnotify.NewWatcher()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer watcher.Close()
	// go func() {
	// 	for {
	// 		select {
	// 		case event, ok := <-watcher.Events:
	// 			if !ok {
	// 				return
	// 			}
	// 			if event.Op&fsnotify.Write == fsnotify.Write {
	// 				if strings.HasSuffix(event.Name, "delegates/delegates.go") {
	// 					watchDelegates(watcher)
	// 				}
	// 				if currentServer != nil {
	// 					syscall.Kill(-currentServer.Process.Pid, syscall.SIGTERM)
	// 				}
	// 			}
	// 		case err, ok := <-watcher.Errors:
	// 			if !ok {
	// 				return
	// 			}
	// 			log.Println(err)
	// 		}
	// 	}
	// }()
	// if err := watcher.Add("."); err != nil {
	// 	log.Fatal(err)
	// }
	// if err := watcher.Add("./delegates/delegates.go"); err != nil {
	// 	log.Fatal(err)
	// }
	// watchDelegates(watcher)
	// // TODO: watch all delegate directories

	go func() {
		for {
			err := runServer()
			if err != nil {
				log.Println(err)
			}
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			if currentServer != nil {
				syscall.Kill(-currentServer.Process.Pid, syscall.SIGTERM)
			}
			os.Exit(0)
		}
	}()

	log.Println("running daemon v0.1...")
	log.Fatal(ListenAndServe(addr))
}

func ListenAndServe(addr string) error {
	server := &qrpc.Server{}
	l, err := mux.ListenWebsocket(addr)
	if err != nil {
		panic(err)
	}
	api := qrpc.NewAPI()
	api.HandleFunc("console", func(r qrpc.Responder, c *qrpc.Call) {
		ch, err := r.Hijack(nil)
		if err != nil {
			log.Println(err)
		}
		logBus.WriteTo(ch)
	})
	return server.Serve(l, api)
}
