CREATE TABLE
    IF NOT EXISTS scoreboards (
        member INTEGER NOT NULL,
        score INTEGER NOT NULL,
        game TEXT NOT NULL,
        PRIMARY KEY (member, game)
    );
