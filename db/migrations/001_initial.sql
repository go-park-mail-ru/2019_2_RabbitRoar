-- -----------------------------------------------------
-- Schema svoyak
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS "svoyak";

-- -----------------------------------------------------
-- Table "svoyak"."User"
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS "svoyak"."User" (
  "id" INT NOT NULL,
  "username" VARCHAR(45) NOT NULL,
  "password" VARCHAR(45) NOT NULL,
  "email" VARCHAR(45) NOT NULL,
  "rating" INT NOT NULL,
  "avatar" VARCHAR(45) NULL,
  PRIMARY KEY ("id"));


-----------------------------------------------------
-- Table "svoyak"."Session"
-----------------------------------------------------
CREATE TABLE IF NOT EXISTS "svoyak"."Session" (
  "UUID" VARCHAR(45) NOT NULL,
  "User_id" INT NOT NULL,
  PRIMARY KEY ("UUID"),
  CONSTRAINT "fk_Session_User"
    FOREIGN KEY ("User_id")
    REFERENCES "svoyak"."User" ("id")
    ON DELETE NO ACTION
    ON UPDATE NO ACTION);


-- -----------------------------------------------------
-- Table "svoyak"."Question"
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS "svoyak"."Question" (
  "id" INT NOT NULL,
  "text" TEXT NOT NULL,
  "media" VARCHAR(45) NULL,
  "answer" VARCHAR(45) NOT NULL,
  "type_id" INT NOT NULL,
  "rating" INT NOT NULL,
  "author" INT NOT NULL,
  "tags" VARCHAR(128) NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_Question_User"
    FOREIGN KEY ("author")
    REFERENCES "svoyak"."User" ("id")
    ON DELETE NO ACTION
    ON UPDATE NO ACTION);


-- -----------------------------------------------------
-- Table "svoyak"."QuestionPack"
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS "svoyak"."QuestionPack" (
  "id" INT NOT NULL,
  "name" VARCHAR(45) NOT NULL,
  "img" VARCHAR(45) NULL,
  "ctime" DATE NOT NULL,
  "mtime" DATE NOT NULL,
  "author" INT NOT NULL,
  "rating" INT NOT NULL,
  "private" BOOLEAN NOT NULL,
  "tags" VARCHAR(256) NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_QuestionPack_User"
    FOREIGN KEY ("author")
    REFERENCES "svoyak"."User" ("id")
    ON DELETE NO ACTION
    ON UPDATE NO ACTION);


-- -----------------------------------------------------
-- Table "svoyak"."PackQuestion"
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS "svoyak"."PackQuestion" (
  "QuestionPack_id" INT NOT NULL,
  "Question_id" INT NOT NULL,
  "theme" VARCHAR(45) NOT NULL ,
  "cost" INT NOT NULL,
  PRIMARY KEY ("QuestionPack_id", "Question_id"),
  CONSTRAINT "fk_QuestionPack_Question"
    FOREIGN KEY ("QuestionPack_id")
    REFERENCES "svoyak"."QuestionPack" ("id")
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT "fk_Question_QuestionPack"
    FOREIGN KEY ("Question_id")
    REFERENCES "svoyak"."Question" ("id")
    ON DELETE NO ACTION
    ON UPDATE NO ACTION);


-- -----------------------------------------------------
-- Table "svoyak"."Game"
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS "svoyak"."Game" (
  "UUID" VARCHAR(45) NOT NULL,
  "name" VARCHAR(45) NOT NULL,
  "players_cap" SMALLINT NOT NULL,
--   "players_joined" SMALLINT NOT NULL, <- optimize me with trigger
  "state" SMALLINT NOT NULL,
  "creator" INT NOT NULL,
  "QuestionPack_id" INT NOT NULL,
  PRIMARY KEY ("UUID"),
  CONSTRAINT "fk_Game_User"
    FOREIGN KEY ("creator")
    REFERENCES "svoyak"."User" ("id")
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT "fk_Game_PackQuestion"
    FOREIGN KEY ("QuestionPack_id")
    REFERENCES "svoyak"."QuestionPack" ("id")
    ON DELETE NO ACTION
    ON UPDATE NO ACTION);


-- -----------------------------------------------------
-- Table "svoyak"."GameUser"
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS "svoyak"."GameUser" (
  "User_id" INT NOT NULL,
  "Game_UUID" VARCHAR(45) NOT NULL,
  PRIMARY KEY ("User_id", "Game_UUID"),
  CONSTRAINT "fk_User_Game"
    FOREIGN KEY ("User_id")
    REFERENCES "svoyak"."User" ("id")
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT "fk_Game_User"
    FOREIGN KEY ("Game_UUID")
    REFERENCES "svoyak"."Game" ("UUID")
    ON DELETE NO ACTION
    ON UPDATE NO ACTION);
