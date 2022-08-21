package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"portalnesia.com/api/config"
	"portalnesia.com/api/models"
	"portalnesia.com/api/routes"
	util "portalnesia.com/api/utils"

	"github.com/cloudflare/tableflip"
	"github.com/joho/godotenv"
)

func main() {
	//log.SetPrefix(fmt.Sprintf("[PID: %d] ", os.Getpid()))

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	config.SetupConfig()
	config.SetupFirebase()
	models.SetupDB()
	if os.Getenv("NODE_ENV") == "production" {
		models.SetupDebugDB()
	}

	r := routes.SetupRouters()

	if runtime.GOOS == "windows" {
		debug()
		r.Listen(fmt.Sprintf(":%s", os.Getenv("PORT")))
	} else {

		var (
			pidFile = flag.String("pid-file", "", "`Path` to pid file")
		)
		upg, err := tableflip.New(tableflip.Options{PIDFile: *pidFile})
		if err != nil {
			log.Fatalf("Error setup graceful: %v", err)
		}
		defer upg.Stop()

		go func() {
			sig := make(chan os.Signal, 1)
			signal.Notify(sig, syscall.SIGHUP)
			for range sig {
				upg.Upgrade()
			}
		}()

		ln, _ := upg.Listen("tcp", fmt.Sprintf("localhost:%s", os.Getenv("PORT")))
		defer ln.Close()

		go r.Listener(ln)

		if err := upg.Ready(); err != nil {
			log.Fatalf("Error graceful ready: %v", err)
		}

		<-upg.Exit()
	}
}

func debug() {
	if os.Getenv("NODE_ENV") == "development" {
		keys := fmt.Sprintf("$2a$08$j9jNyZvS.KFPHIMRAEE4k.ckWmeTMdv17E3QvftgbxEfAO0K94nDm%s", os.Getenv("DEBUG_USERID"))
		token := util.CreateToken(map[string]interface{}{
			"key":      keys,
			"datetime": "2025-05-05 20:20:00",
		})
		fmt.Printf("Debug Session: %s\n", token)
		token = util.CreateToken(map[string]interface{}{
			"token": os.Getenv("AUTH_WEB_SECRET"),
			"date":  "2025-05-05 20:20:00",
		})
		fmt.Printf("Debug PN Auth: %s\n\n", token)
	}
}
