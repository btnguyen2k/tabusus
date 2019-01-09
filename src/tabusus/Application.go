package tabusus

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/gommon/log"
	"os"
	"strconv"
)

const defaultConfigFile = "./config/application.conf"
const defaultListenAddr = "127.0.0.1"
const defaultListenPort = 8080
const staticPath = "/static"

var (
	AppConfig *HoconConfig
	AppDao    ApplicationDao
)

func loadAppConfig() *HoconConfig {
	configFile := os.Getenv("APP_CONFIG")
	if configFile == "" {
		log.Info("No environment APP_CONFIG found, fallback to [", defaultConfigFile, "]")
		configFile = defaultConfigFile
	}
	return LoadAppConfig(configFile)
}

func initDaos(appConfig *HoconConfig) {
	url := appConfig.Conf.GetString("db.mongo.url")
	db := appConfig.Conf.GetString("db.mongo.db")
	AppDao = NewMongoApplicationDao(url, db)
}

func initEcho() *echo.Echo {
	e := echo.New()

	// register static route
	e.Static(staticPath, "public")

	// register template renderer
	e.Renderer = newTemplateRenderer("./views", ".html")

	// register controllers
	s := NewStats()
	e.Use(s.Process)
	e.GET("/stats", s.Handle) // Endpoint to get stats

	e.GET("/logout", actionLogout).Name = "logout"
	e.GET("/login", actionLogin).Name = "login"
	e.POST("/login", actionLoginSubmit).Name = "login"
	e.GET("/apps", actionAppList, RequiredAuthMiddleWare).Name = "apps"
	e.GET("/createApp", actionCreateApp, RequiredAuthMiddleWare).Name = "createApp"
	e.POST("/createApp", actionCreateAppSubmit, RequiredAuthMiddleWare).Name = "createApp"
	e.GET("/editApp/:id", actionEditApp, RequiredAuthMiddleWare).Name = "editApp"
	e.POST("/editApp/:id", actionEditAppSubmit, RequiredAuthMiddleWare).Name = "editApp"
	e.GET("/deleteApp/:id", actionDeleteApp, RequiredAuthMiddleWare).Name = "deleteApp"
	e.POST("/deleteApp/:id", actionDeleteAppSubmit, RequiredAuthMiddleWare).Name = "deleteApp"
	e.GET("/", actionHome, RequiredAuthMiddleWare).Name = "home"

	// register session middleware
	sessionKey := AppConfig.Conf.GetString("session.key", "secret")
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(sessionKey))))
	return e
}

func Start() {
	AppConfig = loadAppConfig()

	initDaos(AppConfig)
	e := initEcho()

	listenAddr := AppConfig.Conf.GetString("http.listen_addr", defaultListenAddr)
	listenPort := AppConfig.Conf.GetInt32("http.listen_port", defaultListenPort)
	e.Logger.Fatal(e.Start(listenAddr + ":" + strconv.Itoa(int(listenPort))))
}
