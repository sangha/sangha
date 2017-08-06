package db

// Project represents the db schema of a project
type Project struct {
	ID         int64
	Slug       string
	Name       string
	About      string
	Website    string
	License    string
	Repository string
	Logo       string
	Activated  bool
}

// LoadProjectByID loads a project by ID from the database
func (context *APIContext) LoadProjectByID(id int64) (Project, error) {
	project := Project{}
	if id < 1 {
		return project, ErrInvalidID
	}

	err := context.QueryRow("SELECT id, slug, name, about, website, license, repository, logo, activated FROM projects WHERE id = $1", id).
		Scan(&project.ID, &project.Slug, &project.Name, &project.About, &project.Website, &project.License, &project.Repository, &project.Logo, &project.Activated)
	return project, err
}

// GetProjectByID returns a project by ID from the cache
func (context *APIContext) GetProjectByID(id int64) (Project, error) {
	project := Project{}
	projectsCache, err := projectsCache.Value(id, context)
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

	err := context.QueryRow("SELECT id, slug, name, about, website, license, repository, logo, activated FROM projects WHERE slug = $1", slug).
		Scan(&project.ID, &project.Slug, &project.Name, &project.About, &project.Website, &project.License, &project.Repository, &project.Logo, &project.Activated)
	return project, err
}

// LoadAllProjects loads all projects from the database
func (context *APIContext) LoadAllProjects() ([]Project, error) {
	projects := []Project{}

	rows, err := context.Query("SELECT id, slug, name, about, website, license, repository, logo, activated FROM projects")
	if err != nil {
		return projects, err
	}

	defer rows.Close()
	for rows.Next() {
		project := Project{}
		err = rows.Scan(&project.ID, &project.Slug, &project.Name, &project.About, &project.Website, &project.License, &project.Repository, &project.Logo, &project.Activated)
		if err != nil {
			return projects, err
		}

		projects = append(projects, project)
	}

	return projects, err
}

// Update a project in the database
func (project *Project) Update(context *APIContext) error {
	_, err := context.Exec("UPDATE projects SET about = $1, slug = $2, name = $3, website = $4, license = $5, repository = $6 WHERE id = $7",
		project.About, project.Slug, project.Name, project.Website, project.License, project.Repository, project.ID)
	if err != nil {
		panic(err)
	}

	projectsCache.Delete(project.ID)
	return err
}

// Save a project to the database
func (project *Project) Save(context *APIContext) error {
	err := context.QueryRow("INSERT INTO projects (slug, name, about, website, license, repository, logo) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		project.Slug, project.Name, project.About, project.Website, project.License, project.Repository, project.Logo).Scan(&project.ID)
	projectsCache.Delete(project.ID)
	return err
}
