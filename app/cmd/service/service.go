package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html"
	"github.com/jinzhu/configor"
	"github.com/laterius/service_architecture_hw3/app/internal/transport/server/httpmw"
	"net/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	dblogger "gorm.io/gorm/logger"

	"github.com/laterius/service_architecture_hw3/app/internal/domain"
	"github.com/laterius/service_architecture_hw3/app/internal/service"
	"github.com/laterius/service_architecture_hw3/app/internal/transport/client/dbrepo"
	transport "github.com/laterius/service_architecture_hw3/app/internal/transport/server/http"
	_ "github.com/laterius/service_architecture_hw3/app/migrations"
)

func main() {
	var cfg domain.Config
	err := configor.New(&configor.Config{Silent: true}).Load(&cfg, "config/config.yaml", "./config.yaml")
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dbrepo.Dsn(cfg.Db),
	}), &gorm.Config{
		Logger: dblogger.Default.LogMode(dblogger.Info),
	})
	if err != nil {
		panic(err)
	}

	userRepo := dbrepo.NewUserRepo(db)
	userService := service.NewUserService(userRepo)
	getUserHandler := transport.NewGetUser(userService)
	postUserHandler := transport.NewPostUser(userService)
	putUserHandler := transport.NewPutUser(userService)
	patchUserHandler := transport.NewPatchUser(userService)
	deleteUserHandler := transport.NewDeleteUser(userService)
	getContactHandler := transport.GetContact()
	getHomeHandler := transport.GetHomePage()
	signUpGet := transport.SignUpGet()
	signUpPost := transport.SignUpPost(userService)
	loginGet := transport.LoginGet()
	loginPost := transport.LoginPost(userService)
	profileGet := transport.NewGetProfile(userService)
	profilePost := transport.NewPostProfile(userService, userService)

	engine := html.New("./views", ".html")
	srv := fiber.New(fiber.Config{Views: engine})
	//srv.Static("/static")

	prometheus := httpmw.New("otus-msa-hw5")
	prometheus.RegisterAt(srv, "/metrics")
	srv.Use(prometheus.Middleware)

	srv.Use(logger.New())
	srv.Use(favicon.New())
	srv.Use(recover.New())
	//srv.Use(httpmw.NewChaosMonkeyMw())

	api := srv.Group("/api")
	api.Post("/user", postUserHandler.Handle())
	api.Get("/user/:id", getUserHandler.Handle())
	api.Put("/user/:id", putUserHandler.Handle())
	api.Patch("/user/:id", patchUserHandler.Handle())
	api.Delete("/user/:id", deleteUserHandler.Handle())

	srv.Get("/probe/live", transport.RespondOk)
	srv.Get("/probe/ready", transport.RespondOk)
	srv.Get("/contact", getContactHandler.Handle())
	srv.Get("/", getHomeHandler.Handle())
	srv.Get("/signup", signUpGet.Handle())
	srv.Post("/signup", signUpPost.Handle())
	srv.Get("/login", loginGet.Handle())
	srv.Post("/login", loginPost.Handle())
	srv.Get("/profile/:token", profileGet.Handle())
	srv.Post("/profile/:id", profilePost.Handle())

	srv.All("/*", transport.DefaultResponse)

	err = srv.Listen(fmt.Sprintf(":%s", cfg.Http.Port))
	if err != nil {
		panic(err)
	}
}

// HTTP 404 NotFound
func notfound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "404 Page not found")
}
