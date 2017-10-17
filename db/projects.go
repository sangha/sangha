package db

import (
	"time"
)

// Project represents the db schema of a project
type Project struct {
	ID             int64
	UUID           string
	Slug           string
	Name           string
	Summary        string
	About          string
	Website        string
	License        string
	Repository     string
	Logo           string
	CreatedAt      time.Time
	Private        bool
	PrivateBalance bool
	Activated      bool
}

// LoadProjectByUUID loads a project by UUID from the database
func (context *APIContext) LoadProjectByUUID(uuid string) (Project, error) {
	project := Project{}
	if len(uuid) == 0 {
		return project, ErrInvalidID
	}

	err := context.QueryRow("SELECT id, uuid, slug, name, summary, about, website, license, repository, logo, created_at, private, private_balance, activated FROM projects WHERE uuid = $1", uuid).
		Scan(&project.ID, &project.UUID, &project.Slug, &project.Name, &project.Summary, &project.About, &project.Website, &project.License, &project.Repository, &project.Logo, &project.CreatedAt, &project.Private, &project.PrivateBalance, &project.Activated)
	return project, err
}

// GetProjectByID loads a project by ID from the database
func (context *APIContext) GetProjectByID(id int64) (Project, error) {
	project := Project{}
	if id == 0 {
		return project, ErrInvalidID
	}

	err := context.QueryRow("SELECT id, uuid, slug, name, summary, about, website, license, repository, logo, created_at, private, private_balance, activated FROM projects WHERE id = $1", id).
		Scan(&project.ID, &project.UUID, &project.Slug, &project.Name, &project.Summary, &project.About, &project.Website, &project.License, &project.Repository, &project.Logo, &project.CreatedAt, &project.Private, &project.PrivateBalance, &project.Activated)
	return project, err
}

// GetProjectByUUID returns a project by UUID from the cache
func (context *APIContext) GetProjectByUUID(uuid string) (Project, error) {
	project := Project{}
	projectsCache, err := projectsCache.Value(uuid, context)
	if err != nil {
		return project, err
	}

	project = *projectsCache.Data().(*Project)
	return project, nil
}

// LoadProjectBySlug loads a project by ID from the database
func (context *APIContext) LoadProjectBySlug(slug string) (Project, error) {
	project := Project{}
	if slug == "" {
		return project, ErrInvalidID
	}

	err := context.QueryRow("SELECT id, uuid, slug, name, summary, about, website, license, repository, logo, created_at, private, private_balance, activated FROM projects WHERE slug = $1", slug).
		Scan(&project.ID, &project.UUID, &project.Slug, &project.Name, &project.Summary, &project.About, &project.Website, &project.License, &project.Repository, &project.Logo, &project.CreatedAt, &project.Private, &project.PrivateBalance, &project.Activated)
	return project, err
}

// LoadAllProjects loads all projects from the database
func (context *APIContext) LoadAllProjects() ([]Project, error) {
	projects := []Project{}

	rows, err := context.Query("SELECT id, uuid, slug, name, summary, about, website, license, repository, logo, created_at, private, private_balance, activated FROM projects")
	if err != nil {
		return projects, err
	}

	defer rows.Close()
	for rows.Next() {
		project := Project{}
		err = rows.Scan(&project.ID, &project.UUID, &project.Slug, &project.Name, &project.Summary, &project.About, &project.Website, &project.License, &project.Repository, &project.Logo, &project.CreatedAt, &project.Private, &project.PrivateBalance, &project.Activated)
		if err != nil {
			return projects, err
		}

		projects = append(projects, project)
	}

	return projects, err
}

// Update a project in the database
func (project *Project) Update(context *APIContext) error {
	_, err := context.Exec("UPDATE projects SET about = $1, summary = $2, slug = $3, name = $4, website = $5, license = $6, repository = $7, private = $8, private_balance = $9 WHERE id = $10",
		project.About, project.Summary, project.Slug, project.Name, project.Website, project.License, project.Repository, project.Private, project.PrivateBalance, project.ID)
	if err != nil {
		panic(err)
	}

	projectsCache.Delete(project.UUID)
	return err
}

// Save a project to the database
func (project *Project) Save(context *APIContext) error {
	project.UUID, _ = UUID()

	err := context.QueryRow("INSERT INTO projects (uuid, slug, name, summary, about, website, license, repository, logo, created_at, private, private_balance) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id",
		project.UUID, project.Slug, project.Name, project.Summary, project.About, project.Website, project.License, project.Repository, project.Logo, project.CreatedAt, project.Private, project.PrivateBalance).Scan(&project.ID)
	projectsCache.Delete(project.UUID)
	return err
}

// Contributors loads all contributors from the database
func (project *Project) Contributors(context *APIContext) ([]User, error) {
	users := []User{}

	rows, err := context.Query("SELECT user_id FROM contributors WHERE project_id = $1", project.ID)
	if err != nil {
		return users, err
	}

	defer rows.Close()
	for rows.Next() {
		var uid int64
		err = rows.Scan(&uid)
		if err != nil {
			return users, err
		}

		user, err := context.LoadUserByID(uid)
		if err != nil {
			return users, err
		}

		users = append(users, user)
	}

	return users, err
}
