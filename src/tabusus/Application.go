package tabusus

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"net/http"
	"os"
	"strconv"
)

const defaultConfigFile = "./config/application.conf"
const defaultListenAddr = "127.0.0.1"
const defaultListenPort = 8080

func Start() {
	configFile := os.Getenv("APP_CONFIG")
	if configFile == "" {
		log.Info("No environment APP_CONFIG found, fallback to [", defaultConfigFile, "]")
		configFile = defaultConfigFile
	}
	appConf := LoadAppConfig(configFile)
	fmt.Println(appConf.Conf)

	listenAddr := appConf.Conf.GetString("http.listen_addr", defaultListenAddr)
	listenPort := appConf.Conf.GetInt32("http.listen_port", defaultListenPort)

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(listenAddr + ":" + strconv.Itoa(int(listenPort))))
}
