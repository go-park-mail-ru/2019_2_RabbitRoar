CREATE SCHEMA IF NOT EXISTS "svoyak";


CREATE TABLE IF NOT EXISTS "svoyak"."User"
(
    "id"       SERIAL,
    "username" VARCHAR(45)  NOT NULL,
    "password" bytea        NOT NULL,
    "email"    VARCHAR(45)  NOT NULL,
    "rating"   INT          NOT NULL,
    "avatar"   VARCHAR(128) NULL,
    "Game_id"  INT,
    PRIMARY KEY ("id"),
    CONSTRAINT "fk_User_Game"
        FOREIGN KEY ("Game_id")
            REFERENCES "svoyak"."Game" ("UUID")
            ON DELETE SET NULL
            ON UPDATE NO ACTION
);


CREATE TABLE IF NOT EXISTS "svoyak"."Session"
(
    "UUID"    VARCHAR(45) NOT NULL,
    "User_id" INT         NOT NULL,
    PRIMARY KEY ("UUID"),
    CONSTRAINT "fk_Session_User"
        FOREIGN KEY ("User_id")
            REFERENCES "svoyak"."User" ("id")
            ON DELETE NO ACTION
            ON UPDATE NO ACTION
);


CREATE TABLE IF NOT EXISTS "svoyak"."Pack"
(
    "id"          INT          NOT NULL,
    "name"        VARCHAR(45)  NOT NULL,
    "author"      INT          NOT NULL,
    "rating"      INT          NOT NULL,
    "description" text         NOT NULL,
    "tags"        VARCHAR(256) NOT NULL,
    "pack"        json         NOT NULL,
    PRIMARY KEY ("id"),
    CONSTRAINT "fk_QuestionPack_User"
        FOREIGN KEY ("author")
            REFERENCES "svoyak"."User" ("id")
            ON DELETE NO ACTION
            ON UPDATE NO ACTION
);


CREATE TABLE IF NOT EXISTS "svoyak"."Game"
(
    "UUID"           VARCHAR(45) NOT NULL,
    "name"           VARCHAR(45) NOT NULL,
    "players_cap"    SMALLINT    NOT NULL,
    "players_joined" SMALLINT    NOT NULL,
    "creator"        INT         NOT NULL,
    "Pack_id"        INT         NOT NULL,
    PRIMARY KEY ("UUID"),
    CONSTRAINT "fk_Game_User"
        FOREIGN KEY ("creator")
            REFERENCES "svoyak"."User" ("id")
            ON DELETE NO ACTION
            ON UPDATE NO ACTION,
    CONSTRAINT "fk_Game_PackQuestion"
        FOREIGN KEY ("Pack_id")
            REFERENCES "svoyak"."Pack" ("id")
            ON DELETE NO ACTION
            ON UPDATE NO ACTION
);


CREATE TABLE IF NOT EXISTS "svoyak"."UserPack"
(
    "User_id" INT NOT NULL,
    "Pack_id" INT NOT NULL,
    PRIMARY KEY ("User_id", "Pack_id"),
    CONSTRAINT "fk_UserPack_User"
        FOREIGN KEY ("User_id")
            REFERENCES "svoyak"."User" ("id")
            ON DELETE NO ACTION
            ON UPDATE NO ACTION,
    CONSTRAINT "fk_UserPack_Pack",
    FOREIGN KEY ("Pack_id")
        REFERENCES svoyak."Pack" ("id")
        ON DELETE NO ACTION
        ON UPDATE NO ACTION
);


-- INSERT INTO "svoyak"."User"
