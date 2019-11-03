package question

import (
	"context"
	"errors"

	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/question"
	"github.com/jackc/pgx/v4"
)

type sqlQuestionRepository struct {
	conn *pgx.Conn
}

func NewSqlQuestionRepository(conn *pgx.Conn) question.Repository {
	return &sqlQuestionRepository{conn}
}

func (repo sqlQuestionRepository) GetByID(questionID int) (*models.Question, error) {
	row := repo.conn.QueryRow(context.Background(), "SELECT id, text, media, answer, rating, author, tags FROM svoyak.Question WHERE id = $1", questionID)

	var question models.Question
	err := row.Scan(&question.ID, &question.Text, &question.Media, &question.Answer, &question.Rating, &question.Author, &question.Tags)

	if err != nil {
		return nil, err
	}

	return &question, nil
}

func (repo sqlQuestionRepository) FetchByTags(tags string, pageSize, page int) (*[]models.Question, error) {
	rows, err := repo.conn.Query(context.Background(), "SELECT id, text, media, answer, rating, author WHERE tags = '$1' OFFSET $2 LIMIT $3", tags, (page * pageSize), pageSize)

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

	rows, err := repo.conn.Query(context.Background(), "SELECT id, text, media, answer, rating, author , tags ORDER BY rating $1", order)

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
	idRow := repo.conn.QueryRow(context.Background(), "INSERT INTO svoyak.Question VALUES (DEFAULT, '$1', '$2', '$3', $4, $5, '$6') RETURNING id", question.Text, question.Media, question.Answer, question.Rating, question.Author, question.Tags)

	err := idRow.Scan(&question.ID)

	if err != nil {
		return nil, err
	}

	return &question, nil
}

func (repo *sqlQuestionRepository) Update(question models.Question) error {
	commandTag, err := repo.conn.Exec(context.Background(), "UPDATE svoyak.Question SET text = '$1', media = '$2', answer = '$3', rating = $4, author = $5, tags = '$6' WHERE id = $7")

	if err != nil {
		return err
	}

	if commandTag.RowsAffected() != 1 {
		return errors.New("No question found to update")
	}

	return nil
}

func (repo *sqlQuestionRepository) Delete(questionID int) error {
	commandTag, err := repo.conn.Exec(context.Background(), "DELETE FROM svoyak.Question WHERE id = '$1'", questionID)

	if err != nil {
		return err
	}

	if commandTag.RowsAffected() != 1 {
		return errors.New("Unable to delete session: No session found")
	}

	return nil
}
