package main

import (
	"database/sql"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB
func initDB() {
	log.Println("Initializing...")
	log.Println(DB)
	r, err := DB.Exec(`
CREATE TABLE "Mood" (
	"id"	INTEGER NOT NULL UNIQUE,
	"mood_name"	INTEGER NOT NULL,
	PRIMARY KEY("id" AUTOINCREMENT)
);`)
	if err != nil {
		log.Println(err)
	}

	rowsAffected, err := r.RowsAffected()
	if err != nil {
		log.Println(err)
	}
	log.Println("Rows affected:", rowsAffected)
}

type Task struct {
	Name string
}

type Mood struct {
	Name string
}

type TemplateData struct {
	ServerIP string
}

func createTask(task Task) {
	DB.Exec("INSERT", task.Name)
}

var CLIENT *websocket.Conn

func main() {
	var ipv4 string
	page := template.Must(template.ParseFiles("view/page.gotmpl"))
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	//listen for events
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) {
					log.Println("file changed:", event.Name)
					// refresh template variable after change
					if CLIENT != nil {
						page = template.Must(template.ParseFiles("view/page.gotmpl"))
						err = CLIENT.WriteMessage(websocket.TextMessage, []byte("FILE_CHANGE"))
						if err != nil {
							log.Println(err)
						}
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()
	err = watcher.Add("/home/violet/repo/lifeos/view/")
	if err != nil {
		log.Println("error:", err)
	}

	err = os.Remove("lifeos.db")
	if os.IsNotExist(err) {
		log.Println("lifeos.db doesn't exist")

		err = os.Remove("lifeos.db-shm")
		if os.IsNotExist(err) {
			log.Println("lifeos.db-shm doesn't exist")
		}

		err = os.Remove("lifeos.db-wal")
		if os.IsNotExist(err) {
			log.Println("lifeos.db-wal doesn't exist")
		}
	}

	dsn := "file:lifeos.db?_journal=WAL"
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		log.Fatal(err)
	}
	if db.Ping() != nil {
		log.Fatal(err)
	}

	DB = db
	initDB()

	upgrader := websocket.Upgrader{}

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		data := TemplateData{
			ServerIP: ipv4,
		}
		page.Execute(w, data)
	})

	http.HandleFunc("GET /websocket", func(w http.ResponseWriter, r *http.Request) {
		log.Println("connection received from", r.RemoteAddr)
		conn, err := upgrader.Upgrade(w, r, nil)
		CLIENT = conn
		if err != nil {
			log.Println(err)
		}
	})

	// grab net interfaces to listen on private networks (192.x.x.x)
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			log.Println("Failed getting addresses for interface", iface.Name, ":", err)
		}
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && ipnet.IP.To4() != nil {
				ip := ipnet.IP.String()
				if strings.HasPrefix(ip, "192") {
					ipv4 = ip
				}
			}
		}
	}

	if ipv4 == "" {
		log.Fatal("Couldn't find any private addresses")
	}

	log.Println("Starting server at", ipv4 + ":1337")
	log.Fatal(http.ListenAndServe(ipv4+":1337", nil))

	db.Close()
}

