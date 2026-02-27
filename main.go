package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"orstax/plugins"
	"orstax/store"
	"orstax/store/sqlstore"

	"github.com/joho/godotenv"
	"go.mau.fi/whatsmeow"
	waLog "go.mau.fi/whatsmeow/util/log"

	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
)

// sourceDir is injected at build time via:
//
//	-ldflags "-X main.sourceDir=/path/to/src"
var sourceDir string

// loadEnv loads .env if present, otherwise falls back to .env.example.
func loadEnv() {
	if err := godotenv.Load(".env"); err != nil {
		_ = godotenv.Load(".env.example")
	}
}

// dbConfig returns the sql dialect and connection address derived from DATABASE_URL.
// A bare filename (no scheme) or a path ending in .db is treated as SQLite;
// anything starting with postgres:// or postgresql:// is treated as PostgreSQL.
func dbConfig() (dialect, addr string) {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		url = "database.db"
	}

	if strings.HasPrefix(url, "postgres://") || strings.HasPrefix(url, "postgresql://") {
		return "postgres", url
	}

	// SQLite – build the connection string with recommended pragmas.
	// Strip a leading "file:" if present so we can normalise the path.
	path := strings.TrimPrefix(url, "file:")
	addr = "file:" + path +
		"?_pragma=foreign_keys(1)" +
		"&_pragma=journal_mode(WAL)" +
		"&_pragma=synchronous(NORMAL)" +
		"&_pragma=busy_timeout(10000)" +
		"&_pragma=cache_size(-64000)" +
		"&_pragma=mmap_size(2147483648)" +
		"&_pragma=temp_store(MEMORY)"
	return "sqlite", addr
}

// getDevice returns the device for the given phone number.
// If phone is empty it falls back to the first stored device (or a new one).
// If phone is provided and no matching device exists, a new (unpaired) device is returned.
func getDevice(ctx context.Context, container *sqlstore.Container, phone string) (*store.Device, error) {
	if phone == "" {
		return container.GetFirstDevice(ctx)
	}

	devices, err := container.GetAllDevices(ctx)
	if err != nil {
		return nil, err
	}
	for _, dev := range devices {
		if dev.ID == nil {
			continue
		}
		// Device JID User field may be "phone.deviceIndex" – compare only the phone part.
		userPhone := strings.SplitN(dev.ID.User, ".", 2)[0]
		if userPhone == phone {
			return dev, nil
		}
	}
	// No existing session for this number – return a fresh device for pairing.
	return container.NewDevice(), nil
}

func main() {
	loadEnv()

	// ── CLI flags ────────────────────────────────────────────────────────────
	phoneArg      := flag.String("phone-number",    "", "Phone number (international format) used to identify or pair a device")
	updateFlag    := flag.Bool("update",             false, "Pull latest source and rebuild the binary in-place")
	listFlag      := flag.Bool("list-sessions",      false, "List all paired sessions stored in the database")
	deleteFlag    := flag.String("delete-session",  "", "Permanently delete the session for the given phone number")
	resetFlag     := flag.String("reset-session",   "", "Reset the session for the given phone number so it can be re-paired")
	flag.Parse()

	ctx := context.Background()

	// ── Management commands (exit after completion) ───────────────────────────
	if *updateFlag {
		runUpdate()
		return
	}

	dialect, dbAddr := dbConfig()

	if *listFlag {
		runListSessions(ctx, dialect, dbAddr)
		return
	}
	if *deleteFlag != "" {
		runDeleteSession(ctx, dialect, dbAddr, *deleteFlag, false)
		return
	}
	if *resetFlag != "" {
		runDeleteSession(ctx, dialect, dbAddr, *resetFlag, true)
		return
	}

	// ── Normal bot startup ────────────────────────────────────────────────────
	dbLog := waLog.Stdout("Database", "ERROR", true)

	container, err := sqlstore.New(ctx, dialect, dbAddr, dbLog)
	if err != nil {
		panic(err)
	}

	if err := plugins.InitDB(container.DB()); err != nil {
		panic(fmt.Errorf("settings db init: %w", err))
	}

	plugins.InitLIDStore(container.LIDMap, "")

	deviceStore, err := getDevice(ctx, container, *phoneArg)
	if err != nil {
		panic(err)
	}

	clientLog := waLog.Stdout("Client", "ERROR", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)
	client.AddEventHandler(plugins.NewHandler(client))

	err = client.Connect()
	if err != nil {
		panic(err)
	}

	if client.Store.ID == nil {
		if *phoneArg == "" {
			fmt.Println("No session found. Please provide a phone number using --phone-number")
			return
		}

		fmt.Println("Waiting 10 seconds before generating pairing code...")
		time.Sleep(10 * time.Second)

		code, err := client.PairPhone(ctx, *phoneArg, true, whatsmeow.PairClientChrome, "Chrome (Linux)")
		if err != nil {
			panic(err)
		}
		fmt.Printf("Your pairing code is: %s\n", code)
	} else {
		ownerPhone := strings.SplitN(client.Store.ID.User, ".", 2)[0]
		plugins.InitLIDStore(container.LIDMap, ownerPhone)
		if err := plugins.InitSettings(ownerPhone); err != nil {
			panic(fmt.Errorf("settings load: %w", err))
		}
		plugins.BootstrapOwnerSudoers()
		fmt.Println("Already logged in.")
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
}

// ── Management command handlers ───────────────────────────────────────────────

// runUpdate pulls the latest source and rebuilds the binary in-place.
func runUpdate() {
	if sourceDir == "" {
		fmt.Fprintln(os.Stderr, "error: this binary was not built with a sourceDir.\nPlease reinstall using the install script.")
		os.Exit(1)
	}

	fmt.Println("Pulling latest changes...")
	pull := exec.Command("git", "-C", sourceDir, "pull")
	pull.Stdout = os.Stdout
	pull.Stderr = os.Stderr
	if err := pull.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "git pull failed: %v\n", err)
		os.Exit(1)
	}

	exePath, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not determine executable path: %v\n", err)
		os.Exit(1)
	}
	exePath, err = filepath.EvalSymlinks(exePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not resolve executable path: %v\n", err)
		os.Exit(1)
	}

	// Build to a temp file first so we never leave a half-written binary.
	tmpPath := exePath + ".new"
	ldflags := fmt.Sprintf("-s -w -X main.sourceDir=%s", sourceDir)

	fmt.Println("Building new binary...")
	build := exec.Command("go", "build",
		"-ldflags", ldflags,
		"-trimpath",
		"-o", tmpPath,
		".",
	)
	build.Dir = sourceDir
	build.Stdout = os.Stdout
	build.Stderr = os.Stderr
	if err := build.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "build failed: %v\n", err)
		_ = os.Remove(tmpPath)
		os.Exit(1)
	}

	if err := os.Rename(tmpPath, exePath); err != nil {
		// On Windows the running binary may be locked; give clear guidance.
		fmt.Fprintf(os.Stderr, "could not replace binary (stop the bot first if it is running): %v\n", err)
		fmt.Fprintf(os.Stderr, "New binary saved at: %s\n", tmpPath)
		fmt.Fprintf(os.Stderr, "Replace manually with:\n  mv %s %s\n", tmpPath, exePath)
		os.Exit(1)
	}

	fmt.Println("Orstax updated successfully.")
}

// runListSessions opens the database and prints all paired sessions.
func runListSessions(ctx context.Context, dialect, dbAddr string) {
	dbLog := waLog.Stdout("Database", "ERROR", true)
	container, err := sqlstore.New(ctx, dialect, dbAddr, dbLog)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open database: %v\n", err)
		os.Exit(1)
	}

	devices, err := container.GetAllDevices(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to list sessions: %v\n", err)
		os.Exit(1)
	}

	if len(devices) == 0 {
		fmt.Println("No sessions found.")
		return
	}

	fmt.Printf("%-4s  %-20s  %s\n", "No.", "Phone", "JID")
	fmt.Println(strings.Repeat("-", 60))
	for i, dev := range devices {
		phone := "(unknown)"
		jid := "(unpaired)"
		if dev.ID != nil {
			phone = strings.SplitN(dev.ID.User, ".", 2)[0]
			jid = dev.ID.String()
		}
		fmt.Printf("%-4d  %-20s  %s\n", i+1, phone, jid)
	}
}

// runDeleteSession removes the stored session for the given phone number.
// When reset is true the message instructs the user to re-pair.
func runDeleteSession(ctx context.Context, dialect, dbAddr, phone string, reset bool) {
	dbLog := waLog.Stdout("Database", "ERROR", true)
	container, err := sqlstore.New(ctx, dialect, dbAddr, dbLog)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open database: %v\n", err)
		os.Exit(1)
	}

	devices, err := container.GetAllDevices(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to query sessions: %v\n", err)
		os.Exit(1)
	}

	for _, dev := range devices {
		if dev.ID == nil {
			continue
		}
		if strings.SplitN(dev.ID.User, ".", 2)[0] == phone {
			if err := container.DeleteDevice(ctx, dev); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to delete session: %v\n", err)
				os.Exit(1)
			}
			if reset {
				fmt.Printf("Session for %s has been reset.\nRun with --phone-number %s to re-pair.\n", phone, phone)
			} else {
				fmt.Printf("Session for %s has been permanently deleted.\n", phone)
			}
			return
		}
	}

	fmt.Fprintf(os.Stderr, "No session found for phone number: %s\n", phone)
	os.Exit(1)
}
