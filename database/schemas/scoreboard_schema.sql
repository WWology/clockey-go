CREATE TABLE
    IF NOT EXISTS scoreboards (
        id INTEGER PRIMARY KEY,
        member INTEGER NOT NULL,
        score INTEGER NOT NULL,
        game TEXT NOT NULL
    );
