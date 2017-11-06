--
-- Projects
--
INSERT INTO projects
    (id, uuid, slug, name, summary, about, website, license, repository, logo, created_at, private, private_balance, activated)
    VALUES (
        1,
        'abcd',
        'toktok',
        'toktok',
        'Typo/error resilient, human-readable token generator',
        'Typo/error resilient, human-readable token generator',
        'https://github.com/muesli/toktok',
        'MIT',
        'https://github.com/muesli/toktok.git',
        '',
        now(),
        false,
        false,
        true
    );
INSERT INTO budgets (uuid, project_id, user_id, parent, name, private, private_balance) VALUES ('b_abcd', 1, null, 0, 'toktok', false, false);

INSERT INTO projects
    (id, uuid, slug, name, summary, about, website, license, repository, logo, created_at, private, private_balance, activated)
    VALUES (
        2,
        'bcde',
        'cache2go',
        'cache2go',
        'Concurrency-safe Go caching library with expiration capabilities and access counters',
        'Concurrency-safe Go caching library with expiration capabilities and access counters',
        'https://github.com/muesli/cache2go',
        'FOSS',
        'https://github.com/muesli/cache2go.git',
        '',
        now(),
        false,
        false,
        true
    );
INSERT INTO budgets (uuid, project_id, user_id, parent, name, private, private_balance) VALUES ('b_bcde', 2, null, 0, 'cache2go', false, false);

INSERT INTO projects
    (id, uuid, slug, name, summary, about, website, license, repository, logo, created_at, private, private_balance, activated)
    VALUES (
        3,
        'cdef',
        'smolder',
        'smolder',
        'smolder makes it easy to write restful Golang JSON APIs',
        'smolder makes it easy to write restful Golang JSON APIs',
        'https://github.com/muesli/smolder',
        'AGPL',
        'https://github.com/muesli/smolder.git',
        '',
        now(),
        false,
        false,
        true
    );
INSERT INTO budgets (uuid, project_id, user_id, parent, name, private, private_balance) VALUES ('b_cdef', 3, null, 0, 'smolder', false, false);

INSERT INTO projects
    (id, uuid, slug, name, summary, about, website, license, repository, logo, created_at, private, private_balance, activated)
    VALUES (
        4,
        'defg',
        'beehive',
        'beehive',
        'A flexible event/agent & automation system with lots of bees 🐝',
        'A flexible event/agent & automation system with lots of bees 🐝',
        'https://github.com/muesli/beehive',
        'AGPL',
        'https://github.com/muesli/beehive.git',
        '',
        now(),
        false,
        false,
        true
    );
INSERT INTO budgets (uuid, project_id, user_id, parent, name, private, private_balance) VALUES ('b_defg', 4, null, 0, 'beehive', false, false);

INSERT INTO projects
    (id, uuid, slug, name, summary, about, website, license, repository, logo, created_at, private, private_balance, activated)
    VALUES (
        5,
        'efgh',
        'mirageos',
        'mirageos',
        'MirageOS is a library operating system that constructs unikernels',
        'MirageOS is a library operating system that constructs unikernels',
        'https://github.com/mirage/mirage',
        'ISC',
        'https://github.com/mirage/mirage.git',
        '',
        now(),
        false,
        false,
        true
    );
INSERT INTO budgets (uuid, project_id, user_id, parent, name, private, private_balance) VALUES ('b_efgh', 5, null, 0, 'mirage', false, false);

--
-- Codes
--
INSERT INTO codes (code, budget_ids, ratios) VALUES ('ABCDEFGH', '{1,2}', '{66,34}');

--
-- Users
--
INSERT INTO users
    (id, uuid, email, nickname, password, authtoken)
    VALUES (
        1,
        'mnop',
        'muesli@gmail.com',
        'muesli',
        '',
        '{9fec2b9fb02e2ec6e9c68351a3bb0c51}'
    );

--
-- Contributors
--
INSERT INTO contributors (user_id, project_id) VALUES (1, 1);
INSERT INTO contributors (user_id, project_id) VALUES (1, 2);
INSERT INTO contributors (user_id, project_id) VALUES (1, 3);
INSERT INTO contributors (user_id, project_id) VALUES (1, 4);
