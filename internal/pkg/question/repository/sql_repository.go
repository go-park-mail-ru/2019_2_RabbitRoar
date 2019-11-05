package question

import (
	"context"
	"errors"

	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/question"
	"github.com/jackc/pgx/v4/pgxpool"
)

type sqlQuestionRepository struct {
	conn *pgxpool.Pool
}

func NewSqlQuestionRepository(conn *pgxpool.Pool) question.Repository {
	return &sqlQuestionRepository{
		conn: conn,
	}
}

func (repo sqlQuestionRepository) GetByID(questionID int) (*models.Question, error) {
	row := repo.conn.QueryRow(
		context.Background(),
		"SELECT id, text, media, answer, rating, author, tags"+
			"FROM svoyak.\"Question\""+
			"WHERE id = $1::integer;",
		questionID,
	)

	var question models.Question
	err := row.Scan(&question.ID, &question.Text, &question.Media, &question.Answer, &question.Rating, &question.Author, &question.Tags)

	return &question, err
}

func (repo sqlQuestionRepository) FetchByTags(tags string, pageSize, page int) (*[]models.Question, error) {
	rows, err := repo.conn.Query(
		context.Background(),
		"SELECT id, text, media, answer, rating, author"+
			"FROM svoyak.\"Question\""+
			"WHERE tags = $1::varchar"+
			"OFFSET $2::integer LIMIT $3::integer;",
		tags, (page * pageSize), pageSize,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var questions []models.Question

	for rows.Next() {
		question := models.Question{
			Tags: tags,
		}

		err := rows.Scan(&question.ID, &question.Text, &question.Media, &question.Answer, &question.Rating, &question.Author)

		if err != nil {
			return nil, err
		}

		questions = append(questions, question)
	}

	return &questions, rows.Err()
}

func (repo sqlQuestionRepository) FetchOrderedByRating(desc bool, pageSize, page int) (*[]models.Question, error) {
	var order string
	if desc {
		order = "DESC"
	} else {
		order = "ASC"
	}

	rows, err := repo.conn.Query(
		context.Background(),
		"SELECT id, text, media, answer, rating, author, tags"+
			"FROM svoyak.\"Question\""+
			"ORDER BY rating $1::text;",
		order,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var questions []models.Question

	for rows.Next() {
		var question models.Question

		err := rows.Scan(&question.ID, &question.Text, &question.Media, &question.Answer, &question.Rating, &question.Author, &question.Tags)

		if err != nil {
			return nil, err
		}

		questions = append(questions, question)
	}

	return &questions, rows.Err()
}

func (repo *sqlQuestionRepository) Create(question models.Question) (*models.Question, error) {
	idRow := repo.conn.QueryRow(
		context.Background(),
		"INSERT INTO svoyak.\"Question\" (id, text, media, answer, rating, author, tags)"+
			"VALUES (DEFAULT, $1::text, $2::varchar, $3::varchar, $4::integer, $5::integer, $6::varchar)"+
			"RETURNING id;",
		question.Text, question.Media, question.Answer, question.Rating, question.Author, question.Tags,
	)

	err := idRow.Scan(&question.ID)

	return &question, err
}

func (repo *sqlQuestionRepository) Update(question models.Question) error {
	commandTag, err := repo.conn.Exec(
		context.Background(),
		"UPDATE svoyak.\"Question\""+
			"SET text = $1::text, media = $2::varchar, answer = $3::varchar, rating = $4::integer, author = $5::integer, tags = $6::varchar"+
			"WHERE id = $7::integer;",
		question.Text, question.Media, question.Answer, question.Rating, question.Author, question.Tags, question.ID,
	)

	if commandTag.RowsAffected() != 1 {
		return errors.New("Unable to update question: No question found")
	}

	return err
}

func (repo *sqlQuestionRepository) Delete(questionID int) error {
	commandTag, err := repo.conn.Exec(
		context.Background(),
		"DELETE FROM svoyak.\"Question\""+
			"WHERE id = $1::integer;",
		questionID,
	)

	if commandTag.RowsAffected() != 1 {
		return errors.New("Unable to delete question: No question found")
	}

	return err
}
