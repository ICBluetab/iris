package main

import (
	"fmt"
	// "encoding/json"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	// "github.com/dgrijalva/jwt-go"
	// jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	"github.com/iris-contrib/middleware/cors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// type Sub struct {
// 	Id int `json:"id"`
// 	UserName string `json:"user_name"`
// 	IsAdmin bool `json:"is_admin"`
// }
//
// type User struct {
// 	UserName string
// 	IsAdmin bool
// 	Id int
// 	HasPermissions bool
// 	GroupIds []string
// 	Jti string
// }

type Concept struct {
	gorm.Model
	Name        string
	Description string
}

func open_connection() (db *gorm.DB, err error) {
	dialect := "postgres"
	config  := "host=localhost port=5432 user=postgres dbname=iris password=postgres sslmode=disable"
	return gorm.Open(dialect, config)
}

func create_tables() {
	db, err := open_connection()
	if err != nil {
		fmt.Println("%s\n", err)
		return
	}

	defer db.Close()
	db.AutoMigrate(&Concept{})
	fmt.Sprintf("AutoMigrate \n")
}

func main() {
	fmt.Println("%s\n", "Entering Main")

	create_tables()

	app := iris.New()
	app.Logger().SetLevel("debug")
	app.Use(recover.New())
	app.Use(logger.New())
	app.Configure(iris.WithConfiguration(iris.YAML("./config/config.yml")))

	crs := cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowCredentials: true,
	})

	// jwtHandler := jwtmiddleware.New(jwtmiddleware.Config{
	// 	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
	// 		return []byte("SuperSecretTruedat"), nil
	// 	},
	// 	SigningMethod: jwt.SigningMethodHS512,
	// })

	// app.Use(jwtHandler.Serve)

	// app.Handle("GET", "/api", func(ctx iris.Context) {
	// 	token := ctx.Values().Get("jwt").(*jwt.Token)
	// 	claims := token.Claims.(jwt.MapClaims)
	//
	// 	var s Sub
	// 	json.Unmarshal([]byte(claims["sub"].(string)), &s)
	//
	// 	user := User{UserName: s.UserName,
	// 		           IsAdmin: s.IsAdmin,
	// 							 Id: s.Id,
	// 							 HasPermissions: claims["has_permissions"].(bool),
	// 							 //GroupIds: claims["gids"], TODO: resolve this
	// 							 Jti: claims["jti"].(string)}
	//
	// 	ctx.Writef("%s", user)
	// })

	api := app.Party("/api", crs).AllowMethods(iris.MethodOptions)
	{

		api.Get("/ping", func(ctx iris.Context) {
			ctx.WriteString("pong")
		})

		api.Get("/concepts/{id:int}", func(ctx iris.Context) {
			db, err := open_connection()
			if err != nil {
				app.Logger().Infof("%s\n", err)
				return
			}
			defer db.Close()

			id, _ := ctx.Params().GetInt("id")
			var concept Concept
			db.Find(&concept, id)
			ctx.Writef("%d\n", concept)
			ctx.WriteString("Show Concept")
		})

		api.Get("/concepts", func(ctx iris.Context) {
			db, err := open_connection()
			if err != nil {
				app.Logger().Infof("%s\n", err)
				return
			}
			defer db.Close()
			var concepts []Concept
			db.Find(&concepts)
			ctx.Writef("%d\n", concepts)
			ctx.WriteString("concepts")
		})

		api.Delete("/concepts/{id:int}", func(ctx iris.Context) {
			db, err := open_connection()
			if err != nil {
				app.Logger().Infof("%s\n", err)
				return
			}
			defer db.Close()

			id, _ := ctx.Params().GetInt("id")
			var concept Concept
			db.Find(&concept, id)
			db.Delete(concept)

			ctx.Writef("%d\n", concept)
			ctx.WriteString("Delete Concept")
		})

		api.Post("/concepts", func(ctx iris.Context) {
			db, err := open_connection()
			if err != nil {
				app.Logger().Infof("%s\n", err)
				return
			}
			defer db.Close()
			concept := Concept{Name: "Nombre del Concepto", Description: "Descripción del Concepto"}
			id := db.Create(&concept)
			ctx.Writef("%d\n", id)
			ctx.WriteString("New Concept")
		})

		api.Put("/concepts/{id:int}", func(ctx iris.Context) {
			db, err := open_connection()
			if err != nil {
				app.Logger().Infof("%s\n", err)
				return
			}
			defer db.Close()

			id, _ := ctx.Params().GetInt("id")
			var concept Concept
			db.Find(&concept, id)

			concept.Name = "Nombre Modificado"
			concept.Description = "Descripción Modificada"

			db.Save(concept)

			ctx.Writef("%d\n", concept)
			ctx.WriteString("Delete Concept")
		})

		api.Get("/hello", func(ctx iris.Context) {
			ctx.JSON(iris.Map{"message": "Hello Iris!"})
		})

	}

	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}
