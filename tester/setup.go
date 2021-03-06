package tester

import (
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/configor"
	"github.com/jinzhu/gorm"
	"github.com/qor/validations"

	"golang-starter-kit/router"

	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func SetupDB() (*gorm.DB, error) {

	var testconfig = struct {
		Test struct {
			Adapter  string `required:"true"`
			Database string `required:"true"`
		}
	}{}

	configor.Load(&testconfig, "../config.yml")

	db, err := gorm.Open(testconfig.Test.Adapter, testconfig.Test.Database)

	if err == nil {
		validations.RegisterCallbacks(db)
	}

	return db, err
}

func GinEngine(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(setDb(db))
	router.InitializeRoutes(r)

	return r
}

func setDb(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("DB", db)
		c.Next()
	}
}

func PerformRequest(r http.Handler, method, path string, args ...io.Reader) *httptest.ResponseRecorder {

	req, _ := http.NewRequest(method, path, nil)
	if len(args) > 0 {
		req, _ = http.NewRequest(method, path, args[0])
		req.Header.Set("Content-Type", "application/json")
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
