CREATE TYPE public.scoreboard_game AS ENUM ('Dota', 'CS', 'MLBB', 'HoK');

CREATE TABLE public.scoreboards (
  member BIGINT NOT NULL,
  score SMALLINT NOT NULL,
  game public.scoreboard_game NOT NULL,
  CONSTRAINT scoreboards_pkey PRIMARY KEY (member, game)
) TABLESPACE pg_default;
