package pgorm

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Connection struct {
	pool           *pgxpool.Pool
	context        context.Context
	database       string
	databaseMaster string
	server         string
	user           string
	password       string
	port           string
	ssl            string
	appName        string
}

func (c *Connection) New(database string) *Connection {
	_ = godotenv.Load()

	server := os.Getenv("ENV_DDBB_SERVER")
	user := os.Getenv("ENV_DDBB_USER")
	password := os.Getenv("ENV_DDBB_PASSWORD")
	port := os.Getenv("ENV_DDBB_PORT")
	ssl := os.Getenv("ENV_DDBB_SSL")
	databaseMaster := os.Getenv("ENV_DDBB_DATABASE")
	appName := os.Getenv("ENV_DDBB_APP")
	if appName == "" {
		appName = "pgorm"
	}
	sslMode := "disable"
	if ssl == "true" {
		sslMode = "require"
	}
	return &Connection{
		context:        context.Background(),
		database:       database,
		databaseMaster: databaseMaster,
		server:         server,
		user:           user,
		password:       password,
		ssl:            sslMode,
		port:           port,
		appName:        appName,
	}
}

func (c *Connection) newPool(isMaster bool) (*Connection, error) {
	database := c.database
	if isMaster {
		database = c.databaseMaster
	}
	host := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&application_name=%s",
		c.user,
		c.password,
		c.server,
		c.port,
		database,
		c.ssl,
		c.appName,
	)

	cnn, err := pgxpool.New(c.context, host)
	if err != nil {
		return nil, fmt.Errorf("conexión error: %w", ManagerErrors{}.SqlConnections(err))
	}

	// Forzar conexión real para detectar errores, porque pool no realiza la conexión hasta realizar la primera consulta
	if err := cnn.Ping(c.context); err != nil {
		return nil, fmt.Errorf("conexión error: %w", ManagerErrors{}.SqlConnections(err))
	}

	c.pool = cnn
	return c, nil
}

func (c *Connection) Pool() (*Connection, error) {
	return c.newPool(false)
}

func (c *Connection) PoolMaster() (*Connection, error) {
	return c.newPool(true)
}

func (c *Connection) GetPool() *pgxpool.Pool {
	return c.pool
}

func (c *Connection) GetContext() context.Context {
	return c.context
}

func (c *Connection) Close() {
	c.pool.Close()
}
