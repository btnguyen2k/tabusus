package tabusus

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func getSession(c echo.Context) *sessions.Session {
	sess, _ := session.Get("session", c)
	return sess
}

func actionLogout(c echo.Context) error {
	sess := getSession(c)
	delete(sess.Values, "uid")
	sess.Save(c.Request(), c.Response())
	return c.Redirect(http.StatusFound, c.Echo().Reverse("home"))
}

func actionLogin(c echo.Context) error {
	return c.Render(http.StatusOK, "login", map[string]interface{}{
	})
}

func actionLoginSubmit(c echo.Context) error {
	id := c.FormValue("user")
	pwd := c.FormValue("password")

	//TODO
	if id != "admin" || pwd != "secret" {
		return c.Render(http.StatusOK, "login", map[string]interface{}{
			"error": "Login failed!",
		})
	}

	sess := getSession(c)
	sess.Values["uid"] = id
	sess.Save(c.Request(), c.Response())
	return c.Redirect(http.StatusFound, c.Echo().Reverse("home"))
}

func actionHome(c echo.Context) error {
	return c.Render(http.StatusOK, "layout:home", map[string]interface{}{
		"active": "home",
	})
}

/*----------------------------------------------------------------------*/

func transformFormData(c echo.Context) map[string]string {
	form, err := c.FormParams()
	if err != nil {
		form = map[string][]string{}
	}
	formData := map[string]string{}
	for k, v := range form {
		if len(v) > 0 {
			formData[k] = v[0]
		}
	}
	return formData
}

func parseRsaPublicKey(keyDataBase64 string) *rsa.PublicKey {
	if !strings.HasPrefix(keyDataBase64, "-----BEGIN PUBLIC KEY-----") && !strings.HasSuffix(keyDataBase64, "-----END PUBLIC KEY-----") {
		keyDataBase64 = "-----BEGIN PUBLIC KEY-----\n" + keyDataBase64 + "\n-----END PUBLIC KEY-----"
	}
	block, _ := pem.Decode([]byte(keyDataBase64))
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil
	}
	if pubkey, err := x509.ParsePKIXPublicKey(block.Bytes); err == nil {
		switch pubkey.(type) {
		case *rsa.PublicKey:
			return pubkey.(*rsa.PublicKey)
		}
	}
	return nil
}

func actionAppList(c echo.Context) error {
	return c.Render(http.StatusOK, "layout:apps", map[string]interface{}{
		"active": "apps",
		"apps":   AppDao.List(),
	})
}

func actionCreateApp(c echo.Context) error {
	formData := transformFormData(c)
	return c.Render(http.StatusOK, "layout:create_edit_app", map[string]interface{}{
		"active": "apps",
		"form":   formData,
	})
}

var validAppId = regexp.MustCompile(`^[a-z0-9_-]+$`)

func actionCreateAppSubmit(c echo.Context) error {
	formData := transformFormData(c)
	var error string

	appId := strings.ToLower(strings.TrimSpace(formData["id"]))
	rsaPubKey := parseRsaPublicKey(formData["pubkey"])
	if !validAppId.MatchString(appId) {
		error = "Invalid application id (must contains only a-z, 0-9, _, -)"
	} else if rsaPubKey == nil {
		error = "Error parsing RSA Public Key data!"
	} else if rsaPubKey.Size() < 1024/8 {
		keySize := strconv.Itoa(rsaPubKey.Size() * 8)
		error = "Key size (" + keySize + ") is less than 1024 bits!"
	} else {
		app, err := AppDao.Get(appId)
		if err != nil {
			error = "Error while checking app [" + appId + "]: " + err.Error() + "!"
		} else if app != nil {
			error = "App [" + appId + "] already existed!"
		}
	}
	if error == "" {
		app := NewApp(appId)
		if formData["enabled"] != "" {
			app.SetStatus(1)
		} else {
			app.SetStatus(0)
		}
		app.SetDescription(formData["desc"])
		app.SetRsaPubKey(formData["pubkey"])
		err := AppDao.Save(app)
		if err != nil {
			error = "Error while saving application [" + appId + "]: " + error
		}
	}
	if error != "" {
		return c.Render(http.StatusOK, "layout:create_edit_app", map[string]interface{}{
			"active": "apps",
			"form":   formData,
			"error":  error,
		})
	} else {
		sess := getSession(c)
		sess.AddFlash("Application [" + appId + "] has been created successfully.")
		sess.Save(c.Request(), c.Response())
		return c.Redirect(http.StatusFound, c.Echo().Reverse("apps"))
	}
}

func actionEditApp(c echo.Context) error {
	appId := c.Param("id")
	app, err := AppDao.Get(appId)
	var error string
	if err != nil {
		error = "Error while getting application info [" + appId + "]!"
	} else if app == nil {
		error = "Application not found [" + appId + "]!"
	}
	formData := transformFormData(c)
	if app.GetStatus() == 1 {
		formData["enabled"] = "1"
	}
	formData["id"] = app.GetId()
	formData["desc"] = app.GetDescription()
	formData["pubkey"] = app.GetRsaPubKey()
	return c.Render(http.StatusOK, "layout:create_edit_app", map[string]interface{}{
		"active":   "apps",
		"form":     formData,
		"error":    error,
		"editMode": true,
	})
}

func actionEditAppSubmit(c echo.Context) error {
	appId := c.Param("id")
	app, err := AppDao.Get(appId)
	var error string
	if err != nil {
		error = "Error while getting application info [" + appId + "]!"
	} else if app == nil {
		error = "Application not found [" + appId + "]!"
	}

	formData := transformFormData(c)
	if error == "" {
		rsaPubKey := parseRsaPublicKey(formData["pubkey"])
		if rsaPubKey == nil {
			error = "Error parsing RSA Public Key data!"
		} else if rsaPubKey.Size() < 1024/8 {
			keySize := strconv.Itoa(rsaPubKey.Size() * 8)
			error = "Key size (" + keySize + ") is less than 1024 bits!"
		}
	}
	if error == "" {
		if formData["enabled"] != "" {
			app.SetStatus(1)
		} else {
			app.SetStatus(0)
		}
		app.SetDescription(formData["desc"])
		app.SetRsaPubKey(formData["pubkey"])
		app.SetTimeUpdated(time.Now())
		err := AppDao.Save(app)
		if err != nil {
			error = "Error while saving application [" + appId + "]: " + error
		}
	}
	if error != "" {
		return c.Render(http.StatusOK, "layout:create_edit_app", map[string]interface{}{
			"active":   "apps",
			"form":     formData,
			"error":    error,
			"editMode": true,
		})
	} else {
		sess := getSession(c)
		sess.AddFlash("Application [" + appId + "] has been updated successfully.")
		sess.Save(c.Request(), c.Response())
		return c.Redirect(http.StatusFound, c.Echo().Reverse("apps"))
	}
}

func actionDeleteApp(c echo.Context) error {
	appId := c.Param("id")
	app, err := AppDao.Get(appId)
	var error string
	if err != nil {
		error = "Error while getting application info [" + appId + "]!"
	} else if app == nil {
		error = "Application not found [" + appId + "]!"
	}
	return c.Render(http.StatusOK, "layout:delete_app", map[string]interface{}{
		"active": "apps",
		"error":  error,
		"app":    app,
	})
}

func actionDeleteAppSubmit(c echo.Context) error {
	appId := c.Param("id")
	app, err := AppDao.Get(appId)
	var error string
	if err != nil {
		error = "Error while getting application info [" + appId + "]: " + err.Error()
	} else if app == nil {
		error = "Application not found [" + appId + "]!"
	}
	if error == "" {
		err := AppDao.Delete(app)
		if err != nil {
			error = "Error while deleting application [" + appId + "]: " + err.Error()
		}
	}
	if error != "" {
		return c.Render(http.StatusOK, "layout:delete_app", map[string]interface{}{
			"active": "apps",
			"error":  error,
			"app":    app,
		})
	} else {
		sess := getSession(c)
		sess.AddFlash("Application [" + appId + "] has been deleted successfully.")
		sess.Save(c.Request(), c.Response())
		return c.Redirect(http.StatusFound, c.Echo().Reverse("apps"))
	}
}
