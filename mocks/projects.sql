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
        'robur',
        'robur',
        'At robur, we build performant bespoke minimal operating systems for high-assurance services',
        'At robur, we build performant bespoke minimal operating systems for high-assurance services',
        'http://robur.io',
        'ISC',
        'https://github.com/mirage/mirage.git',
        '',
        now(),
        false,
        false,
        true
    );
INSERT INTO budgets (uuid, project_id, user_id, parent, name, private, private_balance) VALUES ('b_efgh', 5, null, 0, 'robur', false, false);

INSERT INTO projects
    (id, uuid, slug, name, summary, about, website, license, repository, logo, created_at, private, private_balance, activated)
    VALUES (
        6,
        'fghi',
        'cct',
        'CCT',
        'The Center for the Cultivation of Technology',
        'The Center for the Cultivation of Technology is a charitable non-profit host organization for international Free Software projects',
        'https://techcultivation.org',
        'AGPL',
        'https://gitlab.techcultivation.org',
        '',
        now(),
        false,
        false,
        true
    );
INSERT INTO budgets (uuid, project_id, user_id, parent, name, private, private_balance) VALUES ('b_fghi', 6, null, 0, 'cct', false, false);

INSERT INTO projects
    (id, uuid, slug, name, summary, about, website, license, repository, logo, created_at, private, private_balance, activated)
    VALUES (
        7,
        'ghij',
        'sangha',
        'sangha',
        'sangha',
        'The Center for the Cultivation of Technology is a charitable non-profit host organization for international Free Software projects',
        'https://gitlab.techcultivation.org/sangha/sangha',
        'AGPL',
        'https://gitlab.techcultivation.org/sangha/sangha',
        '',
        now(),
        false,
        false,
        true
    );
INSERT INTO budgets (uuid, project_id, user_id, parent, name, private, private_balance) VALUES ('b_ghij', 7, null, 0, 'sangha', false, false);

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
INSERT INTO users
    (id, uuid, email, nickname, password, authtoken)
    VALUES (
        2,
        'nopq',
        'user2@gmail.com',
        'user2',
        '',
        '{}'
    );
INSERT INTO users
    (id, uuid, email, nickname, password, authtoken)
    VALUES (
        3,
        'opqr',
        'user3@gmail.com',
        'user3',
        '',
        '{}'
    );
INSERT INTO users
    (id, uuid, email, nickname, password, authtoken)
    VALUES (
        4,
        'pqrs',
        'user4@gmail.com',
        'user4',
        '',
        '{}'
    );
INSERT INTO users
    (id, uuid, email, nickname, password, authtoken)
    VALUES (
        5,
        'qrst',
        'user5@gmail.com',
        'user5',
        '',
        '{}'
    );
INSERT INTO users
    (id, uuid, email, nickname, password, authtoken)
    VALUES (
        6,
        'rstu',
        'user6@gmail.com',
        'user6',
        '',
        '{}'
    );
INSERT INTO users
    (id, uuid, email, nickname, password, authtoken)
    VALUES (
        7,
        'stuv',
        'user7@gmail.com',
        'user7',
        '',
        '{}'
    );
INSERT INTO users
    (id, uuid, email, nickname, password, authtoken)
    VALUES (
        8,
        'tuvw',
        'user8@gmail.com',
        'user8',
        '',
        '{}'
    );
INSERT INTO users
    (id, uuid, email, nickname, password, authtoken)
    VALUES (
        9,
        'uvwx',
        'user9@gmail.com',
        'user9',
        '',
        '{}'
    );
INSERT INTO users
    (id, uuid, email, nickname, password, authtoken)
    VALUES (
        10,
        'vwxy',
        'user10@gmail.com',
        'user10',
        '',
        '{}'
    );
INSERT INTO users
    (id, uuid, email, nickname, password, authtoken)
    VALUES (
        11,
        'wxyz',
        'user11@gmail.com',
        'user11',
        '',
        '{}'
    );
INSERT INTO users
    (id, uuid, email, nickname, password, authtoken)
    VALUES (
        12,
        'xyza',
        'user12@gmail.com',
        'user12',
        '',
        '{}'
    );

--
-- Contributors
--
INSERT INTO contributors (user_id, project_id) VALUES (1, 1);
INSERT INTO contributors (user_id, project_id) VALUES (2, 1);
INSERT INTO contributors (user_id, project_id) VALUES (3, 1);
INSERT INTO contributors (user_id, project_id) VALUES (4, 1);
INSERT INTO contributors (user_id, project_id) VALUES (5, 1);
INSERT INTO contributors (user_id, project_id) VALUES (6, 1);
INSERT INTO contributors (user_id, project_id) VALUES (7, 1);
INSERT INTO contributors (user_id, project_id) VALUES (8, 1);
INSERT INTO contributors (user_id, project_id) VALUES (9, 1);
INSERT INTO contributors (user_id, project_id) VALUES (10, 1);
INSERT INTO contributors (user_id, project_id) VALUES (11, 1);
INSERT INTO contributors (user_id, project_id) VALUES (12, 1);
INSERT INTO contributors (user_id, project_id) VALUES (1, 2);
INSERT INTO contributors (user_id, project_id) VALUES (1, 3);
INSERT INTO contributors (user_id, project_id) VALUES (1, 4);
