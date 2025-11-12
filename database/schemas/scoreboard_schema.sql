CREATE TABLE
    IF NOT EXISTS scoreboards (
        member INTEGER PRIMARY KEY,
        score INTEGER NOT NULL,
        game TEXT NOT NULL
    );
