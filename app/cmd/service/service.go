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

	//us, err := models.NewUserService(db)
	//if err != nil {
	//	panic(err)
	//}
	//staticC := controllers.NewStatic() // Parsing static templates
	//usersC := controllers.NewUsers(us) // Handling User Controller
	//r := mux.NewRouter()
	//r.Handle("/contact", staticC.Contact).Methods("GET")
	//r.HandleFunc("/", usersC.New).Methods("GET")
	//r.HandleFunc("/signup", usersC.New).Methods("GET")
	//r.HandleFunc("/signup", usersC.Create).Methods("POST")
	//r.Handle("/login", usersC.LoginView).Methods("GET")
	//r.HandleFunc("/login", usersC.Login).Methods("POST")
	//r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")
	//r.NotFoundHandler = http.HandlerFunc(notfound)
	//fmt.Println("Starting Server on 0.0.0.0:8080")
	//go func() {
	//	err := http.ListenAndServe(":8080", r)
	//	if err != nil {
	//		panic(err)
	//	}
	//}()

	engine := html.New("./views", ".html")
	srv := fiber.New(fiber.Config{Views: engine})

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
