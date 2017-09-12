INSERT INTO projects
    (id, slug, name, about, website, license, repository, logo)
    VALUES (
        1,
        'toktok',
        'toktok',
        'Typo/error resilient, human-readable token generator',
        'https://github.com/muesli/toktok',
        'MIT',
        'https://github.com/muesli/toktok.git',
        ''
    );
INSERT INTO budgets (project_id, parent, name) VALUES (1, 0, 'toktok');

INSERT INTO projects
    (id, slug, name, about, website, license, repository, logo)
    VALUES (
        2,
        'cache2go',
        'cache2go',
        'Concurrency-safe Go caching library with expiration capabilities and access counters',
        'https://github.com/muesli/cache2go',
        'FOSS',
        'https://github.com/muesli/cache2go.git',
        ''
    );
INSERT INTO budgets (project_id, parent, name) VALUES (2, 0, 'cache2go');

INSERT INTO projects
    (id, slug, name, about, website, license, repository, logo)
    VALUES (
        3,
        'smolder',
        'smolder',
        'smolder makes it easy to write restful Golang JSON APIs',
        'https://github.com/muesli/smolder',
        'AGPL',
        'https://github.com/muesli/smolder.git',
        ''
    );
INSERT INTO budgets (project_id, parent, name) VALUES (3, 0, 'smolder');
