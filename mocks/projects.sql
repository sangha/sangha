INSERT INTO projects
    (id, slug, name, about, website, license, repository, logo, created_at)
    VALUES (
        1,
        'toktok',
        'toktok',
        'Typo/error resilient, human-readable token generator',
        'https://github.com/muesli/toktok',
        'MIT',
        'https://github.com/muesli/toktok.git',
        '',
        now()
    );
INSERT INTO budgets (project_id, user_id, parent, name) VALUES (1, null, 0, 'toktok');

INSERT INTO projects
    (id, slug, name, about, website, license, repository, logo, created_at)
    VALUES (
        2,
        'cache2go',
        'cache2go',
        'Concurrency-safe Go caching library with expiration capabilities and access counters',
        'https://github.com/muesli/cache2go',
        'FOSS',
        'https://github.com/muesli/cache2go.git',
        '',
        now()
    );
INSERT INTO budgets (project_id, user_id, parent, name) VALUES (2, null, 0, 'cache2go');

INSERT INTO projects
    (id, slug, name, about, website, license, repository, logo, created_at)
    VALUES (
        3,
        'smolder',
        'smolder',
        'smolder makes it easy to write restful Golang JSON APIs',
        'https://github.com/muesli/smolder',
        'AGPL',
        'https://github.com/muesli/smolder.git',
        '',
        now()
    );
INSERT INTO budgets (project_id, user_id, parent, name) VALUES (3, null, 0, 'smolder');

INSERT INTO projects
    (id, slug, name, about, website, license, repository, logo, created_at)
    VALUES (
        4,
        'beehive',
        'beehive',
        'A flexible event/agent & automation system with lots of bees 🐝',
        'https://github.com/muesli/beehive',
        'AGPL',
        'https://github.com/muesli/beehive.git',
        '08ca2b0038993af9c2f383f0914004eefd33a6df',
        now()
    );
INSERT INTO budgets (project_id, user_id, parent, name) VALUES (4, null, 0, 'beehive');

INSERT INTO codes (code, budget_ids, ratios) VALUES ('ABCDEFGH', '{1,2}', '{66,34}');
