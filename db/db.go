package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"gitlab.techcultivation.org/sangha/sangha/config"

	"github.com/lib/pq"
	"github.com/muesli/cache2go"
	"github.com/satori/go.uuid"
)

var (
	// SchemaVersion needs to be incremented every time the SQL schema changes
	SchemaVersion = 1

	pgDB     *sql.DB
	pgConfig config.PostgreSQLConnection

	projectsCache = cache2go.Cache("project")
	budgetsCache  = cache2go.Cache("budget")
	codesCache    = cache2go.Cache("code")
	usersCache    = cache2go.Cache("user")

	// ErrInvalidID is the error returned when encountering an invalid database ID
	ErrInvalidID = errors.New("Invalid ID")
)

// SetupPostgres sets the db configuration
func SetupPostgres(pc config.PostgreSQLConnection) {
	pgConfig = pc
}

// GetDatabase connects to the database on first run and returns the existing
// connection on further calls
func GetDatabase() *sql.DB {
	if pgDB == nil {
		c := pgConfig
		c.DbName = ""

		db, err := sql.Open("postgres", c.Marshal())
		if err != nil {
			panic(err)
		}
		db.Query("CREATE DATABASE " + pgConfig.DbName)

		pgDB, err = sql.Open("postgres", pgConfig.Marshal())
		if err != nil {
			panic(err)
		}
	}

	return pgDB
}

// InitDatabase sets up the database with all required tables and indexes
func InitDatabase() {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS config
			(
			  name		text	PRIMARY KEY,
			  value		text
			)`,
		`CREATE TABLE IF NOT EXISTS users
			(
			  id          	bigserial 	PRIMARY KEY,
			  uuid			text		NOT NULL,
			  email       	text		NOT NULL,
			  nickname    	text      	NOT NULL,
			  password		text		NOT NULL,
			  about       	text		DEFAULT '',
			  avatar		text		DEFAULT '',
			  address		text[],
			  zip			text		DEFAULT '',
			  city			text		DEFAULT '',
			  country		text		DEFAULT '',
			  activated   	bool		DEFAULT false,
			  authtoken   	text[]     	NOT NULL,
			  CONSTRAINT  	uk_users_uuid 	UNIQUE (uuid),
			  CONSTRAINT  	uk_users_email 	UNIQUE (email)
			)`,
		`CREATE TABLE IF NOT EXISTS projects
			(
			  id          		bigserial 		PRIMARY KEY,
			  uuid				text			NOT NULL,
			  slug				text			NOT NULL,
			  name       		text      		NOT NULL,
			  summary			text			NOT NULL,
			  about				text      		DEFAULT '',
			  website      		text			DEFAULT '',
			  license      		text			DEFAULT '',
			  repository		text			DEFAULT '',
			  logo				text			DEFAULT '',
			  created_at		timestamp		NOT NULL,
			  private			bool			DEFAULT false,
			  private_balance	bool			DEFAULT true,
			  processing_cut	int				DEFAULT 10,
			  activated   		bool			DEFAULT false,
			  CONSTRAINT  		uk_projects_uuid 		UNIQUE (uuid),
			  CONSTRAINT  		uk_projects_slug 		UNIQUE (slug),
			  CONSTRAINT  		uk_projects_repository	UNIQUE (repository)
			)`,
		`CREATE TABLE IF NOT EXISTS budgets
			(
			  id          		bigserial 	PRIMARY KEY,
			  uuid				text		NOT NULL,
			  project_id    	int,
			  user_id			int,
			  parent			bigserial,
			  name       		text      	NOT NULL,
			  description		text,
			  private			bool		DEFAULT false,
			  private_balance	bool		DEFAULT true,
			  CONSTRAINT  		uk_budgets_uuid 		UNIQUE (uuid),
			  CONSTRAINT    	fk_budgets_project_id	FOREIGN KEY (project_id) REFERENCES projects (id) MATCH SIMPLE ON UPDATE CASCADE ON DELETE CASCADE,
			  CONSTRAINT    	fk_budgets_user_id		FOREIGN KEY (user_id) REFERENCES users (id) MATCH SIMPLE ON UPDATE CASCADE ON DELETE CASCADE
			)`,

		`CREATE TABLE IF NOT EXISTS payments
			(
			  id          			bigserial 		PRIMARY KEY,
			  budget_id				bigserial   	NOT NULL,
			  created_at			timestamp		NOT NULL,
			  amount				int				NOT NULL,
			  currency				text			NOT NULL,
			  code					text			DEFAULT '',
			  purpose				text			DEFAULT '',
			  remote_account		text			NOT NULL,
			  remote_name			text			NOT NULL,
			  remote_transaction_id	text			DEFAULT '',
			  remote_bank_id		text			DEFAULT '',
			  source				text			NOT NULL,
			  pending				bool			DEFAULT true,
			  CONSTRAINT    		fk_payments_budget_id	FOREIGN KEY (budget_id) REFERENCES budgets (id) MATCH SIMPLE ON UPDATE CASCADE ON DELETE RESTRICT
			)`,
		`CREATE TABLE IF NOT EXISTS transactions
			(
			  id          		bigserial 		PRIMARY KEY,
			  budget_id			bigserial   	NOT NULL,
			  from_budget_id	int,
			  to_budget_id		int,
			  amount			int				NOT NULL,
			  created_at		timestamp		NOT NULL,
			  purpose			text,
			  payment_id		int,
			  CONSTRAINT    	fk_transactions_budget_id		FOREIGN KEY (budget_id) REFERENCES budgets (id) MATCH SIMPLE ON UPDATE CASCADE ON DELETE RESTRICT,
			  CONSTRAINT    	fk_transactions_from_budget_id	FOREIGN KEY (from_budget_id) REFERENCES budgets (id) MATCH SIMPLE ON UPDATE CASCADE ON DELETE RESTRICT,
			  CONSTRAINT    	fk_transactions_to_budget_id	FOREIGN KEY (to_budget_id) REFERENCES budgets (id) MATCH SIMPLE ON UPDATE CASCADE ON DELETE RESTRICT,
			  CONSTRAINT    	fk_transactions_payment_id		FOREIGN KEY (payment_id) REFERENCES payments (id) MATCH SIMPLE ON UPDATE CASCADE ON DELETE RESTRICT
			)`,
		`CREATE TABLE IF NOT EXISTS codes
			(
			  id			bigserial 		PRIMARY KEY,
			  code			text      		NOT NULL,
			  budget_ids   	int[]			NOT NULL,
			  ratios		int[]			NOT NULL,
			  user_id   	int,
			  CONSTRAINT    uk_codes_code  		UNIQUE (code),
			  CONSTRAINT    uk_codes_budget_ids	UNIQUE (budget_ids, ratios, user_id),
			  CONSTRAINT    fk_codes_user_id	FOREIGN KEY (user_id) REFERENCES users (id) MATCH SIMPLE ON UPDATE CASCADE ON DELETE CASCADE
			)`,
		`CREATE TABLE IF NOT EXISTS contributors
			(
			  id			bigserial 		PRIMARY KEY,
			  user_id   	int,
			  project_id   	int,
			  CONSTRAINT    uk_contributors_user_project	UNIQUE (user_id, project_id),
			  CONSTRAINT    fk_contributors_user_id		FOREIGN KEY (user_id) REFERENCES users (id) MATCH SIMPLE ON UPDATE CASCADE ON DELETE CASCADE,
			  CONSTRAINT    fk_contributors_project_id	FOREIGN KEY (project_id) REFERENCES projects (id) MATCH SIMPLE ON UPDATE CASCADE ON DELETE CASCADE
			)`,
	}

	// FIXME: add IF NOT EXISTS to CREATE INDEX statements (coming in v9.5)
	// See: http://www.postgresql.org/docs/devel/static/sql-createindex.html
	indexes := []string{
		`CREATE INDEX idx_users_uuid ON users(uuid)`,
		`CREATE INDEX idx_users_email ON users(email)`,
		`CREATE INDEX idx_users_authtoken ON users(authtoken)`,
		`CREATE INDEX idx_projects_uuid ON projects(uuid)`,
		`CREATE INDEX idx_projects_slug ON projects(slug)`,
		`CREATE INDEX idx_projects_name ON projects(name)`,
		`CREATE INDEX idx_budgets_uuid ON budgets(uuid)`,
		`CREATE INDEX idx_budgets_name ON budgets(name)`,
		`CREATE INDEX idx_budgets_project_id ON budgets(project_id)`,
		`CREATE INDEX idx_codes_code ON codes(code)`,
		`CREATE INDEX idx_payments_budget_id ON payments(budget_id)`,
		`CREATE INDEX idx_payments_created_at ON payments(created_at)`,
		`CREATE INDEX idx_transactions_budget_id ON transactions(budget_id)`,
		`CREATE INDEX idx_transactions_from_budget_id ON transactions(from_budget_id)`,
		`CREATE INDEX idx_transactions_created_at ON transactions(created_at)`,
		`CREATE INDEX idx_contributors_project_id ON contributors(project_id)`,
	}

	for _, v := range tables {
		fmt.Println("Creating table:", v)
		_, err := pgDB.Exec(v)
		if err != nil {
			panic(err)
		}
	}
	for _, v := range indexes {
		fmt.Println("Creating index:", v)
		_, err := pgDB.Exec(v)
		if err != nil && strings.Index(err.Error(), "already exists") < 0 {
			fmt.Println("Error:", err)
		}
	}
}

// WipeDatabase drops all database tables - use carefully!
func WipeDatabase() {
	drops := []string{
		`DROP TABLE codes`,
		`DROP TABLE contributors`,
		`DROP TABLE payments`,
		`DROP TABLE transactions`,
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
}

// FIXME
func migrateDatabase(from, to int) error {
	return nil
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
				user, err := context.LoadUserByUUID(key.(string))
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
				project, err := context.LoadProjectByUUID(key.(string))
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
				budget, err := context.LoadBudgetByUUID(key.(string))
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

	codesCache.SetAddedItemCallback(func(item *cache2go.CacheItem) {
		// fmt.Println("Now in codes-cache:", item.Key().(string), item.Data().(*DbProject).Name)
	})
	codesCache.SetAboutToDeleteItemCallback(func(item *cache2go.CacheItem) {
		// fmt.Println("Deleting from codes-cache:", item.Key().(string), item.Data().(*DbProject).Name, item.CreatedOn())
	})
	codesCache.SetDataLoader(func(key interface{}, args ...interface{}) *cache2go.CacheItem {
		if len(args) == 1 {
			if context, ok := args[0].(*APIContext); ok {
				code, err := context.LoadCodeByID(key.(int64))
				if err != nil {
					fmt.Println("codesCache ERROR for key", key, ":", err)
					return nil
				}

				entry := cache2go.NewCacheItem(key, 10*time.Minute, &code)
				return entry
			}
		}
		fmt.Println("Got no APIContext passed in")
		return nil
	})
}

func init() {
	initCaches()

	negativeInf := time.Time{}
	positiveInf, _ := time.Parse("2006", "3000")

	pq.EnableInfinityTs(negativeInf, positiveInf)
}
