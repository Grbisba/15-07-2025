package fiber

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
	"go.uber.org/zap"

	"github.com/Grbisba/15-07-2025/internal/config"
	"github.com/Grbisba/15-07-2025/internal/controller"
	"github.com/Grbisba/15-07-2025/internal/service"
)

var (
	_ controller.Controller = (*Controller)(nil)
)

type Controller struct {
	server  *fiber.App
	log     *zap.Logger
	config  *config.Controller
	service service.Service
}

func New(log *zap.Logger, config *config.Config, srv service.Service) (*Controller, error) {
	return newWithConfig(log.Named("controller.http.fiber.insecure"), config.Controller, srv)
}

func newWithConfig(log *zap.Logger, cfg *config.Controller, srv service.Service) (*Controller, error) {
	switch {
	case cfg == nil:
		return nil, errors.Wrap(errNilRef, "nil config")
	case log == nil:
		log = zap.NewNop()
	}

	c := &Controller{
		log:     log,
		config:  cfg,
		service: srv,
	}

	err := c.configure()

	return c, err
}

func (ctrl *Controller) configure() error {
	if err := ctrl.createServer(); err != nil {
		return err
	}

	ctrl.initRoutes()
	ctrl.initMiddlewares()

	return nil
}

func (ctrl *Controller) createServer() error {
	if ctrl == nil {
		return errors.Wrap(errNilRef, "nil controller")
	}

	ctrl.server = fiber.New(fiber.Config{
		CaseSensitive:     true,
		AppName:           "downloader-backend",
		EnablePrintRoutes: true,
		ReadTimeout:       time.Second * time.Duration(ctrl.config.ReadTimeout),
		WriteTimeout:      time.Second * time.Duration(ctrl.config.WriteTimeout),
		IdleTimeout:       time.Second * time.Duration(ctrl.config.IdleTimeout),
	})

	return nil
}

func (ctrl *Controller) initRoutes() {
	ctrl.server.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})
	users := ctrl.server.Group("/users")
	users.Post("/register", ctrl.CreateUser)
	users.Post("/login", ctrl.Login)
	users.Post("/refresh", ctrl.RefreshToken)

	tasks := ctrl.server.Group("tasks")
	tasks.Post("/create", ctrl.CreateTask)
	tasks.Get("/status/:task_id", ctrl.GetStatus)
	tasks.Post("/add-file/:task_id", ctrl.UploadFile)
	tasks.Get("/download/:task_id", ctrl.Download)
}

func (ctrl *Controller) initMiddlewares() {
	ctrl.server.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
	}))
}

func (ctrl *Controller) Start(ctx context.Context) error {
	ch := make(chan error, 1)
	go ctrl.start(ch)
	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Millisecond * time.Duration(ctrl.config.ShutdownTimeout)):
		return nil
	}
}

func (ctrl *Controller) start(ch chan<- error) {
	if ctrl.config.UseTLS {
		ch <- ctrl.server.ListenTLS(ctrl.config.BindAddress(), ctrl.config.CertFile, ctrl.config.KeyFile)
	}
	ch <- ctrl.server.Listen(ctrl.config.BindAddress())
}

func (ctrl *Controller) Stop(_ context.Context) error {
	return multierr.Combine(
		ctrl.server.Shutdown(),
		ctrl.service.Shutdown(),
	)
}

func (ctrl *Controller) ShouldBeRunning() bool {
	return ctrl.config.Enabled
}
