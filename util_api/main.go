package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"fmt"
	"html/template"
	"net/http"
	"util_api/config"
	"util_api/content"
	"util_api/database/mysql"
	"util_api/factory"
	"util_api/kafka"
	"util_api/middleware"
	"util_api/mypay"
	"util_api/newebpay"
	"util_api/url_click"
	"util_api/util"
	"util_api/util/log"

	// "util_api/repository"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

//go:embed templates/*
var f embed.FS

func init() {
	_ = godotenv.Load(".env")
	gin.SetMode(gin.DebugMode)
	config.InitConfig()
	mysql.DatabaseOpen()
}

func main() {

	h := hmac.New(sha256.New, []byte(config.ServerInfo.SecretKey))
	h.Write([]byte("app=Api&key=d1cc009c11a68d8784ffcd4572193115"))
	fmt.Println(hex.EncodeToString(h.Sum(nil)))

	kafka.SetupListen()

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"POST"},
		AllowHeaders:    []string{"Content-Type", "Access-Control-Allow-Origin", "App-Sign-Type", "App-Sign", "App-Key", "Auth-Token", "App-Name"},
	}))

	//網頁模版
	templ := template.Must(template.New("").ParseFS(f, "templates/*.tmpl"))
	router.SetHTMLTemplate(templ)

	apiRouter := router.Group("/")

	//接收高鉅(mypay)訂單完成資料
	router.POST("/mypay-top-up", mypay.MypayTopUp)
	router.POST("/payment", mypay.MypayPayment)

	//藍新金流用
	router.POST("/pagreementNotify", newebpay.PAgreementNotifyForPOST)
	router.POST("/ndnf101Notify", newebpay.Ndnf101NotifyForPOST)
	router.GET("/newebpayPost", newebpay.NewebpayPost)

	//設定 api路徑及連結
	//點擊轉址 api
	router.POST("/url-click", url_click.UrlClickApi)

	middleware.Secret = config.ServerInfo.SecretKey
	apiRouter.Use(middleware.AuthToken())
	{
		apiRouter.POST("/auth", api)
		apiRouter.POST("", api)
	}

	log.Error(router.Run(":" + config.ServerInfo.Port))
}

func api(c *gin.Context) {
	var context content.Context

	userId := c.GetInt64("user_id")
	if err := c.ShouldBindJSON(&context); err != nil {
		log.Error(err)
		c.JSON(http.StatusOK, util.RS{Message: "should bind JSON error", Status: false})
		return
	} else if a, ok := factory.ActionFactoryAuth[context.Action]; ok && userId <= 0 {
		c.JSON(http.StatusOK, util.RS{Message: "api auth failure", Status: false})
		return
	} else {
		if !ok {
			a = factory.ActionFactory[context.Action]
		}
		log.Verbose(context.Action + "\n")
		log.Verbose(context.Parameters + "\n")

		result := factory.LaunchHandler(a, context, c)

		if len(result) > 0 {
			c.SecureJSON(http.StatusOK, result[0].Interface())
			return
		}
	}
}
