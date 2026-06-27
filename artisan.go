package main

import (
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/zhitoo/golang-web-api/app"
	"github.com/zhitoo/golang-web-api/config"
	"github.com/zhitoo/golang-web-api/database/db"
	"github.com/zhitoo/golang-web-api/database/migrations"
)

func run(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("%s failed: %v", name, err)
	}
}

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "swagger:generate":
			run("swag", "init", "-g", "cmd/main.go")
			return
		case "serve:dev":
			run("air")
			return
		}
	}

	cfg := config.GetConfig()

	err := db.InitDb(cfg)
	defer db.CloseDb()
	if err != nil {
		log.Fatal(err.Error())
	}

	// بررسی آرگومان‌های خط فرمان
	if len(os.Args) > 1 {
		command := os.Args[1]

		dbClient := db.GetDb()
		defer db.CloseDb()

		switch command {
		case "migrate:up":
			if err := migrations.Up(dbClient); err != nil {
				log.Fatalf("Migration up failed: %v", err)
			}
			return // خروج پس از پایان مایگریشن

		case "migrate:down":
			steps := 1
			if len(os.Args) > 2 {
				stepsInt, err := strconv.Atoi(os.Args[2])
				if err != nil {
					log.Fatalf("Migration down failed: %v", err)
				}
				steps = stepsInt
			}

			if err := migrations.Down(dbClient, steps); err != nil {
				log.Fatalf("Migration down failed: %v", err)
			}
			return

		case "migrate:create":
			if len(os.Args) < 3 {
				log.Fatal("Please provide a migration name. Example: go run main.go migrate-create init_users")
			}
			name := os.Args[2]
			if err := migrations.Create(name); err != nil {
				log.Fatalf("Create migration failed: %v", err)
			}
			return
		case "migrate:force":
			if len(os.Args) < 3 {
				log.Fatal("Please provide a version number. Example: go run main.go migrate:force 1780506576")
			}
			version, err := strconv.Atoi(os.Args[2])
			if err != nil {
				log.Fatalf("Invalid version number: %v", err)
			}
			if err := migrations.Force(dbClient, version); err != nil {
				log.Fatalf("Migration force failed: %v", err)
			}
			return
		case "serve":
			app.InitServer(cfg)
		default:
			log.Fatalf("Unknown command: %s", command)

		}

	}

}
