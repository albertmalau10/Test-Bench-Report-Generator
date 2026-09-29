package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"valve_database/controllers"
	"valve_database/database"
	"valve_database/middleware"
)

// proxyPrefix MUST match services.proxyMapping[].url in
// packaging/configs/package-assets/valve-database-app.package-manifest.json
// exactly. ctrlX CORE's reverse proxy forwards the full request path —
// prefix included — to this app's unix socket, so the app strips it
// itself before routing. Empty / unused outside a snap.
const proxyPrefix = "/valve-database-app"

// snapName MUST match snapcraft.yaml's top-level "name" and the
// package-manifest.json's "id" — it's the sub-path both the socket
// file and the package-run content slot are published under.
const snapName = "valve-database-app"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env tidak ditemukan, pakai default")
	}

	// Inside a snap, $SNAP is the read-only, versioned app bundle and
	// $SNAP_DATA is the writable directory snapd guarantees persists
	// across restarts/updates. Every relative path this app already
	// uses (data.db, ./images, ./datasheets, client_cert.pem,
	// client_key.pem) resolves against the process's working
	// directory, so switching that directory once here — before
	// anything else runs — makes the whole app snap-safe without
	// touching those files individually.
	if snapData := os.Getenv("SNAP_DATA"); snapData != "" {
		if err := os.Chdir(snapData); err != nil {
			log.Fatalf("failed to chdir into SNAP_DATA %s: %v", snapData, err)
		}
	}
	isSnap := os.Getenv("SNAP") != ""

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080"
	}

	db := database.ConnectDB()
	database.SeedUsers(db)

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("gagal mengambil koneksi sql:", err)
	}
	defer sqlDB.Close()

<<<<<<< HEAD
	// FIX
	os.MkdirAll("./images", os.ModePerm)
	os.MkdirAll("./datasheets", os.ModePerm)

	router := gin.Default()
=======
	os.MkdirAll("./images", os.ModePerm)
	os.MkdirAll("./datasheets", os.ModePerm)

	router := gin.Default()

	if !isSnap {
		// Local development only: the Vite dev server (a different
		// origin, :5173) calls this API directly. Once behind ctrlX
		// CORE's reverse proxy, frontend and API share one origin and
		// no CORS headers are needed at all.
		router.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"http://localhost:5173"},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))
	}
>>>>>>> pr-theme-update

	// All backend endpoints live under /api so they never collide with
	// a client-side (Vue Router) page of the same name — GET /valves
	// and GET /settings are both real frontend pages AND real API
	// routes; without this split, a hard refresh on those pages would
	// hit the API instead of the app shell once frontend and backend
	// share one origin behind ctrlX CORE's reverse proxy.
	api := router.Group("/api")
	{
		api.GET("/ping", controllers.Ping)
		api.POST("/login", controllers.Login(db))

		api.GET("/valves", controllers.GetValves(db))
		api.GET("/valves/:id", controllers.GetValveByID(db))

		api.PUT("/valves/:id", middleware.AuthRequired(), controllers.UpdateValve(db))

		api.POST("/valves", middleware.AuthRequired(), middleware.AdminOnly(), controllers.CreateValve(db))
		api.DELETE("/valves/:id", middleware.AuthRequired(), middleware.AdminOnly(), controllers.DeleteValve(db))
		api.POST("/valves/:id/image", middleware.AuthRequired(), middleware.AdminOnly(), controllers.UploadValveImage(db))
		api.POST("/valves/:id/datasheet", middleware.AuthRequired(), middleware.AdminOnly(), controllers.UploadValveDatasheet(db))

		api.GET("/settings", middleware.AuthRequired(), middleware.AdminOnly(), controllers.GetSettings(db))
		api.PUT("/settings", middleware.AuthRequired(), middleware.AdminOnly(), controllers.UpdateSettings(db))

		// OPC UA routes
		api.GET("/opcua/data", middleware.AuthRequired(), controllers.GetOpcData(db))
		api.POST("/ctrlx/start", middleware.AuthRequired(), controllers.StartOutput(db))
		api.GET("/opcua/status", middleware.AuthRequired(), controllers.GetOpcStatus(db))
		api.POST("/ctrlx/stop", middleware.AuthRequired(), controllers.StopOutput(db))

<<<<<<< HEAD
	router.GET("/settings", middleware.AuthRequired(), middleware.AdminOnly(), controllers.GetSettings(db))
	router.PUT("/settings", middleware.AuthRequired(), middleware.AdminOnly(), controllers.UpdateSettings(db))
	
	//OPC UA Routes
	router.GET("/opcua/data", middleware.AuthRequired(), controllers.GetOpcData(db))
	router.POST("/ctrlx/start", middleware.AuthRequired(), controllers.StartOutput(db))
	router.GET("/opcua/status", middleware.AuthRequired(), controllers.GetOpcStatus(db))
	router.POST("/ctrlx/stop", middleware.AuthRequired(), controllers.StopOutput(db))

	router.Run(":" + appPort)
=======
		api.Static("/images", "./images")
		api.Static("/datasheets", "./datasheets")
	}

	// Everything else: the built Vue frontend (frontend/dist, copied
	// into wwwDir at packaging time — see build-frontend.sh). Falls
	// back to index.html for any path that isn't a real file, so a
	// hard refresh on a client-side route (e.g. /valves/5) still works
	// — Vue Router then takes over from index.html.
	router.NoRoute(gin.WrapH(spaHandler(filepath.Join(wwwRoot(), "www"))))

	var handler http.Handler = router
	if isSnap {
		handler = http.StripPrefix(proxyPrefix, router)
		listenUnixSocket(handler)
		return
	}

	log.Printf("Valve Database App backend listening on :%s", appPort)
	if err := http.ListenAndServe(":"+appPort, handler); err != nil {
		log.Fatal("server stopped:", err)
	}
}

// spaHandler serves the built frontend and falls back to index.html
// for any path that isn't a real file, so a hard refresh on a
// client-side route still works (Vue Router then takes over).
func spaHandler(wwwDir string) http.Handler {
	fileServer := http.FileServer(http.Dir(wwwDir))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		full := filepath.Join(wwwDir, filepath.Clean(r.URL.Path))
		if info, err := os.Stat(full); err == nil && !info.IsDir() {
			fileServer.ServeHTTP(w, r)
			return
		}
		http.ServeFile(w, r, filepath.Join(wwwDir, "index.html"))
	})
}

// listenUnixSocket binds the reverse-proxy socket ctrlX CORE expects,
// at exactly the path package-manifest.json's "binding" declares:
// unix://{$SNAP_DATA}/package-run/valve-database-app/valve-database-app.sock
func listenUnixSocket(handler http.Handler) {
	sockDir := filepath.Join(os.Getenv("SNAP_DATA"), "package-run", snapName)
	sockFile := filepath.Join(sockDir, snapName+".sock")

	if err := os.MkdirAll(sockDir, 0o755); err != nil {
		log.Fatalf("failed to create socket directory %s: %v", sockDir, err)
	}
	os.Remove(sockFile) // clean up a stale socket left by a previous run

	listener, err := net.Listen("unix", sockFile)
	if err != nil {
		log.Fatalf("failed to listen on unix socket %s: %v", sockFile, err)
	}

	log.Printf("Valve Database App backend listening on unix socket %s", sockFile)
	if err := http.Serve(listener, handler); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}

// wwwRoot is where the built frontend lives relative to the binary:
// $SNAP inside a snap (read-only, from the build), the backend's own
// working directory in local development.
func wwwRoot() string {
	if dir := os.Getenv("SNAP"); dir != "" {
		return dir
	}
	return "."
>>>>>>> pr-theme-update
}
