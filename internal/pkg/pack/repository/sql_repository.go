package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/pack"
)

//TODO: move order by generation to function
//TODO: move offset and page generation to function

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
		&p.Offline,
		&questions,
	)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(questions, &p.Questions); err != nil {
		return nil, err
	}

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

func (repo sqlPackRepository) Played(packID, userID int) (bool, error) {
	row := repo.db.QueryRow(
		`
			SELECT 1
			FROM "svoyak"."GameUserHist"
			WHERE "Pack_id" = $1::integer AND "User_id" = $2::integer;
		`,
		packID, userID,
	)
	var played int
	err := row.Scan(&played)
	if err != nil {
		return false, nil
	}
	if played == 1 {
		return true, nil
	}
	return false, nil
}

func (repo sqlPackRepository) GetByID(packID int) (*models.Pack, error) {
	row := repo.db.QueryRow(
		`
			SELECT id, name, description, rating, author, tags, offline, pack
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
	var pids = make([]int, 0, 20)
	for rows.Next() {
		if err := rows.Scan(&pid); err != nil {
			return nil, err
		}
		pids = append(pids, pid)
	}
	return pids, nil
}

func (repo sqlPackRepository) FetchOfflineAuthor(caller models.User) ([]int, error) {
	rows, err := repo.db.Query(
		`
		SELECT id
		FROM "svoyak"."Pack"
		WHERE author = $1::integer
	`, caller.ID)
	if err != nil {
		return nil, err
	}

	var id int
	var ids = make([]int, 0, 20)
	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
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

func scanPackRows(rows *sql.Rows, pageSize int) ([]models.Pack, error) {
	var packs = make([]models.Pack, 0, pageSize)

	for rows.Next() {
		var p models.Pack
		err :=  rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.Rating,
			&p.Author,
			&p.Tags,
		)

		if err != nil {
			return nil, err
		}

		packs = append(packs, p)
	}

	return packs, nil
}

func (repo sqlPackRepository) FetchOrderedByRating(desc bool, page, pageSize int) ([]models.Pack, error) {
	var order string
	if desc {
		order = "DESC"
	} else {
		order = "ASC"
	}

	query := fmt.Sprintf(
		`
			SELECT id, name, description, rating, author, tags
			FROM "svoyak"."Pack"
			ORDER BY rating %s
			OFFSET $1::integer LIMIT $2::integer;
		`, order,
	)

	rows, err := repo.db.Query(query, pageSize * page, pageSize)

	defer rows.Close()

	if err != nil {
		return nil, err
	}

	return scanPackRows(rows, pageSize)
}

func (repo sqlPackRepository) FetchByAuthor(u models.User, desc bool, page, pageSize int) ([]models.Pack, error) {
	var order string
	if desc {
		order = "DESC"
	} else {
		order = "ASC"
	}

	query := fmt.Sprintf(
	`
		SELECT id, name, description, rating, author, tags
		FROM "svoyak"."Pack"
		WHERE author = $1
		ORDER BY rating %s
		OFFSET $2::integer LIMIT $3::integer;
	`, order)
	rows, err := repo.db.Query(query, u.ID, page * pageSize, pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanPackRows(rows, pageSize)
}

func (repo sqlPackRepository) FetchByTags(tags string, desc bool, page, pageSize int) ([]models.Pack, error) {
	var order string
	if desc {
		order = "DESC"
	} else {
		order = "ASC"
	}

	query := fmt.Sprintf(
		`
			SELECT id, name, description, rating, author, tags
			FROM "svoyak"."Pack"
			WHERE tags = $1::varchar
			ORDER BY rating %s
			OFFSET $2::integer LIMIT $3::integer;
		`,
		order,
	)

	rows, err := repo.db.Query(query, tags, page * pageSize, pageSize)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return scanPackRows(rows, pageSize)
}
