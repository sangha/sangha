package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"gitlab.techcultivation.org/techcultivation/sangha/config"

	"github.com/lib/pq"
	"github.com/muesli/cache2go"
	uuid "github.com/nu7hatch/gouuid"
)

var (
	pgDB   *sql.DB
	pgConn config.PostgreSQLConnection

	projectsCache = cache2go.Cache("project")
	budgetsCache  = cache2go.Cache("budget")
	usersCache    = cache2go.Cache("user")

	// ErrInvalidID is the error returned when encountering an invalid database ID
	ErrInvalidID = errors.New("Invalid id")
)

// SetupPostgres sets the db configuration
func SetupPostgres(pc config.PostgreSQLConnection) {
	pgConn = pc
}

// GetDatabase connects to the database on first run and returns the existing
// connection on further calls
func GetDatabase() *sql.DB {
	if pgDB == nil {
		var err error
		pgDB, err = sql.Open("postgres", pgConn.Marshal())
		if err != nil {
			panic(err)
		}

		tables := []string{
			`CREATE TABLE IF NOT EXISTS users
				(
				  id          	bigserial 	PRIMARY KEY,
				  email       	text		NOT NULL,
				  nickname    	text      	NOT NULL,
				  password		text		NOT NULL,
				  about       	text		DEFAULT '',
				  activated   	bool		DEFAULT false,
				  authtoken   	text[]     	NOT NULL,
				  CONSTRAINT  	uk_email 	UNIQUE (email)
				)`,
			`CREATE TABLE IF NOT EXISTS projects
				(
				  id          	bigserial 	PRIMARY KEY,
				  name       	text      	NOT NULL,
				  about			text      	NOT NULL,
				  website      	text		DEFAULT '',
				  license      	text		DEFAULT '',
				  repository	text		DEFAULT '',
				  activated   	bool		DEFAULT false
				)`,
			`CREATE TABLE IF NOT EXISTS budgets
				(
				  id          	bigserial 	PRIMARY KEY,
				  project_id    bigserial   NOT NULL,
				  name       	text      	NOT NULL,
				  CONSTRAINT    fk_project  FOREIGN KEY (project_id) REFERENCES projects (id) MATCH SIMPLE ON UPDATE CASCADE ON DELETE CASCADE
				)`,
		}

		// FIXME: add IF NOT EXISTS to CREATE INDEX statements (coming in v9.5)
		// See: http://www.postgresql.org/docs/devel/static/sql-createindex.html
		indexes := []string{
			`CREATE INDEX idx_users_email ON users(email)`,
			`CREATE INDEX idx_users_authtoken ON users(authtoken)`,
			`CREATE INDEX idx_projects_name ON projects(name)`,
			`CREATE INDEX idx_budgets_name ON budgets(name)`,
			`CREATE INDEX idx_budgets_project_id ON budgets(project_id)`,
		}

		for _, v := range tables {
			fmt.Println("Creating table:", v)
			_, err = pgDB.Exec(v)
			if err != nil {
				panic(err)
			}
		}
		for _, v := range indexes {
			fmt.Println("Creating index:", v)
			_, err = pgDB.Exec(v)
			if err != nil && strings.Index(err.Error(), "already exists") < 0 {
				fmt.Println("Error:", err)
			}
		}
	}

	return pgDB
}

// WipeDatabase drops all database tables - use carefully!
func WipeDatabase() {
	// Commented out to prevent accidental usage

	/*
		drops := []string{
			`DROP TABLE budgets`,
			`DROP TABLE projects`,
			`DROP TABLE users`,
		}

		for _, v := range drops {
			fmt.Println("Dropping table:", v)
			_, err := pgDB.Exec(v)
			if err != nil {
				panic(err)
			}
		}
	*/
}

func init() {
	fmt.Println("db.init")
	initCaches()

	negativeInf := time.Time{}
	positiveInf, _ := time.Parse("2006", "3000")

	pq.EnableInfinityTs(negativeInf, positiveInf)
}

// UUID returns a new unique identifier
func UUID() (string, error) {
	u, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	uuid := strings.Join(strings.Split(u.String(), "-"), "")
	return uuid, nil
}

func initCaches() {
	usersCache.SetAddedItemCallback(func(item *cache2go.CacheItem) {
		// fmt.Println("Now in users-cache:", item.Key().(string), item.Data().(*DbUser).Username)
	})
	usersCache.SetAboutToDeleteItemCallback(func(item *cache2go.CacheItem) {
		// fmt.Println("Deleting from users-cache:", item.Key().(string), item.Data().(*DbUser).Username, item.CreatedOn())
	})
	usersCache.SetDataLoader(func(key interface{}, args ...interface{}) *cache2go.CacheItem {
		if len(args) == 1 {
			if context, ok := args[0].(*APIContext); ok {
				user, err := context.LoadUserByID(key.(int64))
				if err != nil {
					fmt.Println("usersCache ERROR for key", key, ":", err)
					return nil
				}

				entry := cache2go.NewCacheItem(key, 10*time.Minute, &user)
				return entry
			}
		}
		fmt.Println("Got no APIContext passed in")
		return nil
	})

	projectsCache.SetAddedItemCallback(func(item *cache2go.CacheItem) {
		// fmt.Println("Now in projects-cache:", item.Key().(string), item.Data().(*DbProject).Name)
	})
	projectsCache.SetAboutToDeleteItemCallback(func(item *cache2go.CacheItem) {
		// fmt.Println("Deleting from projects-cache:", item.Key().(string), item.Data().(*DbProject).Name, item.CreatedOn())
	})
	projectsCache.SetDataLoader(func(key interface{}, args ...interface{}) *cache2go.CacheItem {
		if len(args) == 1 {
			if context, ok := args[0].(*APIContext); ok {
				project, err := context.LoadProjectByID(key.(int64))
				if err != nil {
					fmt.Println("projectsCache ERROR for key", key, ":", err)
					return nil
				}

				entry := cache2go.NewCacheItem(key, 10*time.Minute, &project)
				return entry
			}
		}
		fmt.Println("Got no APIContext passed in")
		return nil
	})

	budgetsCache.SetAddedItemCallback(func(item *cache2go.CacheItem) {
		// fmt.Println("Now in budgets-cache:", item.Key().(string), item.Data().(*DbProject).Name)
	})
	budgetsCache.SetAboutToDeleteItemCallback(func(item *cache2go.CacheItem) {
		// fmt.Println("Deleting from budgets-cache:", item.Key().(string), item.Data().(*DbProject).Name, item.CreatedOn())
	})
	budgetsCache.SetDataLoader(func(key interface{}, args ...interface{}) *cache2go.CacheItem {
		if len(args) == 1 {
			if context, ok := args[0].(*APIContext); ok {
				budget, err := context.LoadBudgetByID(key.(int64))
				if err != nil {
					fmt.Println("budgetsCache ERROR for key", key, ":", err)
					return nil
				}

				entry := cache2go.NewCacheItem(key, 10*time.Minute, &budget)
				return entry
			}
		}
		fmt.Println("Got no APIContext passed in")
		return nil
	})
}
