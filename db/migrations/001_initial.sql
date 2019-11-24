CREATE SCHEMA IF NOT EXISTS "svoyak";


CREATE TABLE IF NOT EXISTS "svoyak"."User"
(
    "id"       SERIAL       NOT NULL UNIQUE,
    "username" VARCHAR(45)  NOT NULL UNIQUE,
    "password" bytea        NOT NULL,
    "email"    VARCHAR(45)  NOT NULL UNIQUE,
    "rating"   INT          NOT NULL,
    "avatar"   VARCHAR(128) NULL,
    PRIMARY KEY ("id")
);


CREATE TABLE IF NOT EXISTS "svoyak"."Session"
(
    "UUID"    VARCHAR(45) NOT NULL,
    "User_id" INT         NOT NULL,
    PRIMARY KEY ("UUID"),
    CONSTRAINT "fk_Session_User"
        FOREIGN KEY ("User_id")
            REFERENCES "svoyak"."User" ("id")
            ON DELETE CASCADE
            ON UPDATE CASCADE
);


CREATE TABLE IF NOT EXISTS "svoyak"."Pack"
(
    "id"          SERIAL             NOT NULL,
    "name"        VARCHAR(45)        NOT NULL,
    "author"      INT                NOT NULL,
    "rating"      INT  DEFAULT 0     NOT NULL,
    "description" text               NOT NULL,
    "tags"        VARCHAR(256)       NOT NULL,
    "pack"        json               NOT NULL,
    "offline"     bool DEFAULT FALSE NOT NULL,
    PRIMARY KEY ("id"),
    CONSTRAINT "fk_QuestionPack_User"
        FOREIGN KEY ("author")
            REFERENCES "svoyak"."User" ("id")
            ON DELETE CASCADE
            ON UPDATE CASCADE
);


CREATE TABLE IF NOT EXISTS "svoyak"."GameUser"
(
    "User_id"   INT         NOT NULL UNIQUE,
    "Game_UUID" VARCHAR(45) NOT NULL,
    PRIMARY KEY ("User_id", "Game_UUID"),
    CONSTRAINT "fk_GameUser_User_id"
        FOREIGN KEY ("User_id")
            REFERENCES "svoyak"."User" ("id")
            ON DELETE CASCADE
            ON UPDATE CASCADE
);


CREATE TABLE IF NOT EXISTS "svoyak"."GameUserHist"
(
    "User_id" INT NOT NULL,
    "Pack_id" INT NOT NULL,
    PRIMARY KEY ("User_id", "Pack_id"),
    CONSTRAINT "fk_UserPack_User"
        FOREIGN KEY ("User_id")
        REFERENCES "svoyak"."User" ("id")
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT "fk_UserPack_Pack"
        FOREIGN KEY ("Pack_id")
            REFERENCES svoyak."Pack" ("id")
            ON DELETE CASCADE
            ON UPDATE CASCADE
);
