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
--   INDEX "username_idx" ("username" ASC) VISIBLE,
--   INDEX "email_idx" ("email" ASC) VISIBLE);


-----------------------------------------------------
-- Table "svoyak"."Session"
-----------------------------------------------------
CREATE TABLE IF NOT EXISTS "svoyak"."Session" (
  "UUID" INT NOT NULL,
  "User_id" INT NOT NULL,
  PRIMARY KEY ("UUID"),
--   UNIQUE INDEX "UUID_UNIQUE" ("UUID" ASC) VISIBLE,
--   INDEX "fk_Session_User_idx" ("User_id" ASC) VISIBLE,
  CONSTRAINT "fk_Session_User"
    FOREIGN KEY ("User_id")
    REFERENCES "svoyak"."User" ("id")
    ON DELETE NO ACTION
    ON UPDATE NO ACTION);


-- -----------------------------------------------------
-- Table "svoyak"."QuestionType"
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS "svoyak"."QuestionType" (
  "id" INT NOT NULL,
  "type" VARCHAR(45) NULL,
  PRIMARY KEY ("id"));


-- -----------------------------------------------------
-- Table "svoyak"."Question"
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS "svoyak"."Question" (
  "id" INT NOT NULL,
  "text" TEXT NOT NULL,
  "img" VARCHAR(45) NULL,
  "answer" VARCHAR(45) NOT NULL,
  "type_id" INT NOT NULL,
  "rating" INT NOT NULL,
  "author" INT NOT NULL,
  "tags" VARCHAR(128) NULL,
  PRIMARY KEY ("id"),
--   INDEX "fk_Question_QuestionType1_idx" ("type_id" ASC) VISIBLE,
--   INDEX "fk_Question_User1_idx" ("author" ASC) VISIBLE,
--   FULLTEXT INDEX "tags_fulltext_idx" ("tags") VISIBLE,
  CONSTRAINT "fk_Question_QuestionType"
    FOREIGN KEY ("type_id")
    REFERENCES "svoyak"."QuestionType" ("id")
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
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
  "ctime" DATE NOT NULL,
  "mtime" DATE NOT NULL,
  "author" INT NOT NULL,
  "rating" INT NOT NULL,
  "private" BOOLEAN NOT NULL,
  "tags" VARCHAR(256) NULL,
  PRIMARY KEY ("id"),
--   INDEX "fk_QuestionPack_User1_idx" ("author" ASC) VISIBLE,
--   INDEX "name_idx" ("name" ASC) VISIBLE,
--   INDEX "author_id" ("author" ASC) VISIBLE,
--   FULLTEXT INDEX "tags_fulltext_idx" () VISIBLE,
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
  PRIMARY KEY ("QuestionPack_id", "Question_id"),
--   INDEX "fk_QuestionPack_has_Question_Question1_idx" ("Question_id" ASC) VISIBLE,
--   INDEX "fk_QuestionPack_has_Question_QuestionPack1_idx" ("QuestionPack_id" ASC) VISIBLE,
  CONSTRAINT "fk_QuestionPack_has_Question_QuestionPack1"
    FOREIGN KEY ("QuestionPack_id")
    REFERENCES "svoyak"."QuestionPack" ("id")
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT "fk_QuestionPack_Question"
    FOREIGN KEY ("Question_id")
    REFERENCES "svoyak"."Question" ("id")
    ON DELETE NO ACTION
    ON UPDATE NO ACTION);


-- -----------------------------------------------------
-- Table "svoyak"."GameState"
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS "svoyak"."GameState" (
  "id" INT NOT NULL,
  "state" VARCHAR(45) NOT NULL,
  PRIMARY KEY ("id"));


-- -----------------------------------------------------
-- Table "svoyak"."Game"
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS "svoyak"."Game" (
  "UUID" VARCHAR(45) NOT NULL,
  "name" VARCHAR(45) NOT NULL,
  "size" SMALLINT NOT NULL,
  "GameState_id" INT NOT NULL,
  "creator" INT NOT NULL,
  "QuestionPack_id" INT NOT NULL,
  PRIMARY KEY ("UUID"),
--   INDEX "fk_Game_GameState1_idx" ("GameState_id" ASC) VISIBLE,
--   INDEX "fk_Game_User1_idx" ("creator" ASC) VISIBLE,
--   INDEX "fk_Game_PackQuestion1_idx" ("PackQuestion_QuestionPack_id" ASC, "PackQuestion_Question_id" ASC) VISIBLE,
--   INDEX "name_idx" ("name" ASC) VISIBLE,
  CONSTRAINT "fk_Game_GameState"
    FOREIGN KEY ("GameState_id")
    REFERENCES "svoyak"."GameState" ("id")
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
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
--   INDEX "fk_User_has_Game_Game1_idx" ("Game_UUID" ASC) VISIBLE,
--   INDEX "fk_User_has_Game_User1_idx" ("User_id" ASC) VISIBLE,
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
