package db

import (
	"errors"
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
	ProcessingCut  int64
	Activated      bool
	UserID         *int64
}

// LoadProjectByUUID loads a project by UUID from the database
func (context *APIContext) LoadProjectByUUID(uuid string) (Project, error) {
	project := Project{}
	if len(uuid) == 0 {
		return project, ErrInvalidID
	}

	err := context.QueryRow("SELECT id, uuid, slug, name, summary, about, website, license, repository, logo, created_at, private, private_balance, processing_cut, activated, user_id FROM projects WHERE uuid = $1", uuid).
		Scan(&project.ID, &project.UUID, &project.Slug, &project.Name, &project.Summary, &project.About, &project.Website, &project.License, &project.Repository, &project.Logo, &project.CreatedAt, &project.Private, &project.PrivateBalance, &project.ProcessingCut, &project.Activated, &project.UserID)

	if !project.HasAccess(context.Auth) {
		return Project{}, errors.New("No such project")
	}
	return project, err
}

// GetProjectByID loads a project by ID from the database
func (context *APIContext) GetProjectByID(id int64) (Project, error) {
	project := Project{}
	if id == 0 {
		return project, ErrInvalidID
	}

	err := context.QueryRow("SELECT id, uuid, slug, name, summary, about, website, license, repository, logo, created_at, private, private_balance, processing_cut, activated, user_id FROM projects WHERE id = $1", id).
		Scan(&project.ID, &project.UUID, &project.Slug, &project.Name, &project.Summary, &project.About, &project.Website, &project.License, &project.Repository, &project.Logo, &project.CreatedAt, &project.Private, &project.PrivateBalance, &project.ProcessingCut, &project.Activated, &project.UserID)

	if !project.HasAccess(context.Auth) {
		return Project{}, errors.New("No such project")
	}
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

	if !project.HasAccess(context.Auth) {
		return Project{}, errors.New("No such project")
	}
	return project, nil
}

// LoadProjectBySlug loads a project by ID from the database
func (context *APIContext) LoadProjectBySlug(slug string) (Project, error) {
	project := Project{}
	if slug == "" {
		return project, ErrInvalidID
	}

	err := context.QueryRow("SELECT id, uuid, slug, name, summary, about, website, license, repository, logo, created_at, private, private_balance, processing_cut, activated, user_id FROM projects WHERE slug = $1", slug).
		Scan(&project.ID, &project.UUID, &project.Slug, &project.Name, &project.Summary, &project.About, &project.Website, &project.License, &project.Repository, &project.Logo, &project.CreatedAt, &project.Private, &project.PrivateBalance, &project.ProcessingCut, &project.Activated, &project.UserID)

	if !project.HasAccess(context.Auth) {
		return Project{}, errors.New("No such project")
	}
	return project, err
}

// LoadAllProjects loads all projects from the database
func (context *APIContext) LoadAllProjects() ([]Project, error) {
	projects := []Project{}

	rows, err := context.Query("SELECT id, uuid, slug, name, summary, about, website, license, repository, logo, created_at, private, private_balance, processing_cut, activated, user_id FROM projects")
	if err != nil {
		return projects, err
	}

	defer rows.Close()
	for rows.Next() {
		project := Project{}
		err = rows.Scan(&project.ID, &project.UUID, &project.Slug, &project.Name, &project.Summary, &project.About, &project.Website, &project.License, &project.Repository, &project.Logo, &project.CreatedAt, &project.Private, &project.PrivateBalance, &project.ProcessingCut, &project.Activated, &project.UserID)
		if err != nil {
			return projects, err
		}

		if !project.HasAccess(context.Auth) {
			continue
		}
		projects = append(projects, project)
	}

	return projects, err
}

// Update a project in the database
func (project *Project) Update(context *APIContext) error {
	_, err := context.Exec("UPDATE projects SET about = $1, summary = $2, slug = $3, name = $4, website = $5, license = $6, repository = $7, private = $8, private_balance = $9, processing_cut = $10, activated = $11 WHERE id = $12",
		project.About, project.Summary, project.Slug, project.Name, project.Website, project.License, project.Repository, project.Private, project.PrivateBalance, project.ProcessingCut, project.Activated, project.ID)

	projectsCache.Delete(project.UUID)
	return err
}

// Save a project to the database
func (project *Project) Save(context *APIContext) error {
	project.UUID, _ = UUID()

	err := context.QueryRow("INSERT INTO projects (uuid, slug, name, summary, about, website, license, repository, logo, created_at, private, private_balance) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id",
		project.UUID, project.Slug, project.Name, project.Summary, project.About, project.Website, project.License, project.Repository, project.Logo, time.Now().UTC(), project.Private, project.PrivateBalance).Scan(&project.ID)

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

// Balance returns this project's total balance
func (project *Project) Balance(context *APIContext) (int64, error) {
	var b int64

	budgets, err := context.LoadBudgets(project)
	if err != nil {
		return 0, err
	}

	for _, budget := range budgets {
		bal, err := budget.Balance(context)
		if err != nil {
			return 0, err
		}
		b += bal
	}

	return b, nil
}

// BalanceStats returns this project's total balances for the past months
func (project *Project) BalanceStats(context *APIContext) ([]int64, error) {
	var b []int64

	budgets, err := context.LoadBudgets(project)
	if err != nil {
		return b, err
	}

	for _, budget := range budgets {
		bal, err := budget.BalanceStats(context)
		if err != nil {
			return b, err
		}

		for idx, v := range bal {
			if idx >= len(b) {
				b = append(b, 0)
			}
			b[idx] += v
		}
	}

	return b, nil
}

func (project *Project) HasAccess(user *User) bool {
	return !project.Private || user.ID == 1 || (project.UserID != nil && user.ID == *project.UserID)
}

func (project *Project) HasTransactionAccess(user *User) bool {
	return !project.PrivateBalance || user.ID == 1 || (project.UserID != nil && user.ID == *project.UserID)
}

// SearchProjects searches database for projects
func (context *APIContext) SearchProjects(term string) ([]Project, error) {
	projects := []Project{}

	rows, err := context.Query("SELECT DISTINCT id FROM projects "+
		"WHERE (LOWER(name) LIKE LOWER('%' || $1 || '%') OR "+
		"LOWER(summary) LIKE LOWER('%' || $1 || '%'))", term)
	if err != nil {
		return projects, err
	}

	defer rows.Close()
	for rows.Next() {
		var id int64
		err = rows.Scan(&id)
		if err != nil {
			return projects, err
		}

		p, err := context.GetProjectByID(id)
		if err != nil {
			return projects, err
		}

		projects = append(projects, p)
	}

	return projects, nil
}
