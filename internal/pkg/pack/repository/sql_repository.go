package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/pack"
	"github.com/labstack/gommon/log"

	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
)

type sqlPackRepository struct {
	db *sql.DB
}

func NewSqlPackRepository(db *sql.DB) pack.Repository {
	return &sqlPackRepository{
		db: db,
	}
}

func scanPackRow(row *sql.Row) (*models.Pack, error) {
	var p models.Pack
	var questions []byte

	err := row.Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.Rating,
		&p.Author,
		&p.Tags,
		&questions,
	)

	if err = json.Unmarshal(questions, &p.Questions); err != nil {
		return nil, err
	}

	return &p, err
}

func scanPackRows(rows *sql.Rows) (*models.Pack, error) {
	var p models.Pack
	err := rows.Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.Rating,
		&p.Author,
		&p.Tags,
	)

	return &p, err
}

func (repo *sqlPackRepository) Create(pack *models.Pack) error {
	pPack, err := json.Marshal(pack.Questions)
	if err != nil {
		return nil
	}
	idRow := repo.db.QueryRow(
		`
			INSERT INTO "svoyak"."Pack" (id, name, description, rating, author, tags, pack)
			VALUES (DEFAULT, $1::varchar, $2::text, $3::integer, $4::integer, $5::varchar, $6::json)
			RETURNING id;
		`,
		pack.Name, pack.Description, pack.Rating, pack.Author, pack.Tags, pPack,
	)

	return idRow.Scan(&pack.ID)
}

func (repo *sqlPackRepository) Update(pack *models.Pack) error {
	res, err := repo.db.Exec(
		`
			UPDATE "svoyak"."Pack"
			SET name        = $1::varchar,
    			description = $2::text,
    			rating      = $3::integer,
    			author      = $4::integer,
    			tags        = $5::varchar,
    			pack        = $6::json
			WHERE id = $7::integer;
		`,
		pack.Name, pack.Description, pack.Rating, pack.Author, pack.Tags,pack.Questions, pack.ID,
	)
	if err != nil {
		return nil
	}

	c, err := res.RowsAffected()
	if c != 1 {
		return errors.New("unable to update pack: no pack found")
	}

	return err
}

func (repo *sqlPackRepository) Delete(packID int) error {
	res, err := repo.db.Exec(
		`
			DELETE FROM "svoyak"."Pack"
			WHERE id = $1::integer;
		`,
		packID,
	)
	if err != nil {
		return err
	}

	c, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if c != 1 {
		return errors.New("unable to delete pack: no pack found")
	}

	return nil
}


func (repo sqlPackRepository) GetByID(packID int) (*models.Pack, error) {
	row := repo.db.QueryRow(
		`
			SELECT id, name, description, rating, author, tags, pack
			FROM "svoyak"."Pack"
			WHERE id = $1::integer;
		`,
		packID,
	)

	return scanPackRow(row)
}

func (repo sqlPackRepository) FetchOfflinePublic() ([]int, error) {
	rows, err := repo.db.Query(
	`
		SELECT id
		FROM "svoyak"."Pack"
		WHERE offline = TRUE
	`)
	if err != nil {
		return nil, err
	}

	var pid int
	var pids []int
	for rows.Next() {
		if err := rows.Scan(&pid); err != nil {
			return nil, err
		}
		pids = append(pids, pid)
	}
	return pids, nil
}

func (repo sqlPackRepository) FetchOffline(caller models.User) ([]int, error) {
	rows, err := repo.db.Query(
	`
		SELECT DISTINCT "Pack_id"
		FROM "svoyak"."GameUserHist"
		WHERE "User_id" = $1::int
	`,
	caller.ID,
	)

	if err != nil {
		return nil, err
	}

	var pid int
	var pids = make([]int, 0, 64)
	for rows.Next() {
		if err := rows.Scan(&pid); err != nil {
			return nil, err
		}
		pids = append(pids, pid)
	}
	return pids, nil
}

func (repo sqlPackRepository) FetchOrderedByRating(desc bool, page, pageSize int) ([]models.Pack, error) {
	var order string
	if desc {
		order = "DESC"
	} else {
		order = "ASC"
	}

	rows, err := repo.db.Query(
		fmt.Sprintf(
		`
			SELECT id, name, description, rating, author, tags
			FROM "svoyak"."Pack"
			ORDER BY rating %s;
		`, order,
		),
	)

	if err != nil {
		return nil, err
	}

	var packs = make([]models.Pack, 0, pageSize)

	for rows.Next() {
		p, err := scanPackRows(rows)

		if err != nil {
			return nil, err
		}

		packs = append(packs, *p)
	}

	defer rows.Close()

	return packs, nil
}

func (repo sqlPackRepository) FetchByAuthor(u models.User) ([]models.Pack, error) {
	return nil, nil
}

func (repo sqlPackRepository) FetchByTags(tags string, page, pageSize int) ([]models.Pack, error) {
	rows, err := repo.db.Query(
		`
			SELECT id, name, description, rating, author, tags
			FROM "svoyak"."Pack"
			WHERE tags = $1::varchar
			OFFSET $2::integer LIMIT $3::integer;
		`,
		tags, page * pageSize, pageSize,
	)

	if err != nil {
		return nil, err
	}

	defer log.Error(rows.Close())

	var packs = make([]models.Pack, 0, pageSize)

	for rows.Next() {
		p, err := scanPackRows(rows)

		if err != nil {
			return nil, err
		}

		packs = append(packs, *p)
	}

	return packs, rows.Err()
}
