package main

import (
	"app/handlers/auth"
	authService "app/internal/core/services/auth"
	"context"
	"crypto/rsa"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ohmspeed777/go-pkg/logx"
	"github.com/ohmspeed777/go-pkg/middlewares"
	"github.com/spf13/viper"
	"github.com/tylerb/graceful"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func init() {
	logx.Init("trace", true)
	initViper()
}

func main() {
	var (
		mongoDB = newMongoDB()
		privKey = initRsa()
		e       = initEcho(privKey)
	)

	api := e.Group("/api/v1")

	// handler zone
	authHandler := auth.NewHandler(authService.NewService(mongoDB))

	// endpoint zone
	authGroup := api.Group("/auth")
	authGroup.POST("/sign-up", authHandler.SignUp)
	authGroup.POST("/sign-in", authHandler.SignIn)

	logx.GetLog().Infof("server starting on port: %d", viper.GetInt("app.port"))
	_ = graceful.ListenAndServe(e.Server, 5*time.Second)
}

func initEcho(key *rsa.PrivateKey) *echo.Echo {
	e := echo.New()
	e.Use(middlewares.LogRequestMiddleware(key))
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	e.Use(middlewares.LogResponseMiddleware())
	e.HTTPErrorHandler = middlewares.CustomHTTPErrorHandler
	e.Server.Addr = fmt.Sprintf(":%d", viper.GetInt("app.port"))
	logx.GetLog().Infof("http server started on port: %d", viper.GetInt("app.port"))
	return e
}

func initViper() {
	viper.SetConfigName("config")
	viper.AddConfigPath("configs")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		logx.GetLog().Fatalf("cannot read in viper config:%s", err)
	}
}

func newMongoDB() *mongo.Database {
	clientOptions := options.Client()
	clientOptions.SetHosts([]string{viper.GetString("mongo.host")})
	clientOptions.SetAuth(options.Credential{
		Username: viper.GetString("mongo.username"),
		Password: viper.GetString("mongo.password"),
	})

	logx.GetLog().Info("[CONFIG] Mongo host:", []string{viper.GetString("mongo.host")})
	logx.GetLog().Info("[CONFIG] Mongo database:", viper.GetString("mongo.database"))

	clientOptions.SetAppName(viper.GetString("app.name"))

	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		logx.GetLog().Fatalf("cannot NewClient MongoDB :%v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		logx.GetLog().Fatalf("cannot Connect MongoDB :%v", err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		logx.GetLog().Fatalf("cannot Ping MongoDB :%v", err)
	}

	return client.Database(viper.GetString("mongo.database"))
}

func initRsa() *rsa.PrivateKey {
	pkbytes := []byte(viper.GetString("jwt.key"))

	privateKeyImported, err := jwt.ParseRSAPrivateKeyFromPEM(pkbytes)
	if err != nil {
		logx.GetLog().Fatal(err)
		return nil
	}

	return privateKeyImported
}
