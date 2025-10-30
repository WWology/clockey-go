CREATE TABLE
    IF NOT EXISTS events (
        id INTEGER PRIMARY KEY,
        name TEXT NOT NULL,
        time INTEGER NOT NULL,
        type TEXT NOT NULL,
        gardener INTEGER NOT NULL CHECK (gardener in (293360731867316225, 204923365205475329, 754724309276164159, 172360818715918337, 332438787588227072)),
        hours INTEGER NOT NULL
    );

CREATE TABLE
    IF NOT EXISTS dota_scoreboard (
        id INTEGER PRIMARY KEY,
        score INTEGER NOT NULL
    );

CREATE TABLE
    IF NOT EXISTS cs_scoreboard (
        id INTEGER PRIMARY KEY,
        score INTEGER NOT NULL
    );

CREATE TABLE
    IF NOT EXISTS mlbb_scoreboard (
        id INTEGER PRIMARY KEY,
        score INTEGER NOT NULL
    );

CREATE TABLE
    IF NOT EXISTS hok_scoreboard (
        id INTEGER PRIMARY KEY,
        score INTEGER NOT NULL
    );
