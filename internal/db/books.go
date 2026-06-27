package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/joaovictornsv/books-cli/internal/models"
)

var ErrNotFound = errors.New("book not found")

type ListFilter struct {
	Status          *models.Status
	PriorityToBuy   *bool
	EligibleToSell  *bool
	IncludeArchived bool
	Pagination      *models.Pagination
}

type SearchFilter struct {
	Query           string
	Author          string
	IncludeArchived bool
	Pagination      *models.Pagination
}

type BooksResult struct {
	Books []models.Book
	Total int
}

type Repository struct {
	db *DB
}

func NewRepository(db *DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, book models.Book) (models.Book, error) {
	if book.AddedAt == "" {
		book.AddedAt = models.NowTimestamp()
	}
	if err := book.Validate(); err != nil {
		return models.Book{}, err
	}

	res, err := r.db.sql.ExecContext(ctx, `
		INSERT INTO books (
			title, author, status, priority_to_buy, eligible_to_sell, sold,
			notes, added_at, started_at, finished_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		book.Title,
		nullString(book.Author),
		book.Status.String(),
		book.PriorityToBuy,
		book.EligibleToSell,
		book.Sold,
		nullString(book.Notes),
		book.AddedAt,
		nullString(book.StartedAt),
		nullString(book.FinishedAt),
	)
	if err != nil {
		return models.Book{}, fmt.Errorf("insert book: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return models.Book{}, fmt.Errorf("last insert id: %w", err)
	}

	return r.GetByID(ctx, id)
}

func (r *Repository) GetByID(ctx context.Context, id int64) (models.Book, error) {
	row := r.db.sql.QueryRowContext(ctx, `
		SELECT id, title, author, status, priority_to_buy, eligible_to_sell, sold,
		       notes, added_at, started_at, finished_at
		FROM books WHERE id = ?`, id)

	book, err := scanBook(row)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Book{}, ErrNotFound
	}
	if err != nil {
		return models.Book{}, err
	}
	return book, nil
}

func (r *Repository) List(ctx context.Context, filter ListFilter) (BooksResult, error) {
	where, args := buildListWhere(filter)

	total, err := r.countBooks(ctx, where, args)
	if err != nil {
		return BooksResult{}, err
	}

	query := `
		SELECT id, title, author, status, priority_to_buy, eligible_to_sell, sold,
		       notes, added_at, started_at, finished_at
		FROM books WHERE 1=1` + where + ` ORDER BY id ASC`

	queryArgs := append([]any{}, args...)
	if filter.Pagination != nil && filter.Pagination.Enabled() {
		query += ` LIMIT ? OFFSET ?`
		queryArgs = append(queryArgs, filter.Pagination.Limit, filter.Pagination.Offset())
	}

	books, err := r.queryBooks(ctx, query, queryArgs...)
	if err != nil {
		return BooksResult{}, err
	}
	return BooksResult{Books: books, Total: total}, nil
}

func (r *Repository) Search(ctx context.Context, filter SearchFilter) (BooksResult, error) {
	where, args := buildSearchWhere(filter)

	total, err := r.countBooks(ctx, where, args)
	if err != nil {
		return BooksResult{}, err
	}

	query := `
		SELECT id, title, author, status, priority_to_buy, eligible_to_sell, sold,
		       notes, added_at, started_at, finished_at
		FROM books WHERE 1=1` + where + ` ORDER BY id ASC`

	queryArgs := append([]any{}, args...)
	if filter.Pagination != nil && filter.Pagination.Enabled() {
		query += ` LIMIT ? OFFSET ?`
		queryArgs = append(queryArgs, filter.Pagination.Limit, filter.Pagination.Offset())
	}

	books, err := r.queryBooks(ctx, query, queryArgs...)
	if err != nil {
		return BooksResult{}, err
	}
	return BooksResult{Books: books, Total: total}, nil
}

func buildListWhere(filter ListFilter) (string, []any) {
	query := ""
	args := []any{}

	if !filter.IncludeArchived {
		query += ` AND status != ?`
		args = append(args, models.StatusArchived.String())
	}
	if filter.Status != nil {
		query += ` AND status = ?`
		args = append(args, filter.Status.String())
	}
	if filter.PriorityToBuy != nil {
		query += ` AND priority_to_buy = ?`
		args = append(args, models.ToBool01(*filter.PriorityToBuy))
	}
	if filter.EligibleToSell != nil {
		query += ` AND eligible_to_sell = ?`
		args = append(args, models.ToBool01(*filter.EligibleToSell))
	}
	return query, args
}

func buildSearchWhere(filter SearchFilter) (string, []any) {
	query := ""
	args := []any{}

	if !filter.IncludeArchived {
		query += ` AND status != ?`
		args = append(args, models.StatusArchived.String())
	}
	if q := strings.TrimSpace(filter.Query); q != "" {
		query += ` AND LOWER(title) LIKE ? ESCAPE '\'`
		args = append(args, "%"+escapeLike(strings.ToLower(q))+"%")
	}
	if a := strings.TrimSpace(filter.Author); a != "" {
		query += ` AND LOWER(COALESCE(author, '')) LIKE ? ESCAPE '\'`
		args = append(args, "%"+escapeLike(strings.ToLower(a))+"%")
	}
	return query, args
}

func (r *Repository) countBooks(ctx context.Context, where string, args []any) (int, error) {
	var total int
	err := r.db.sql.QueryRowContext(ctx, `SELECT COUNT(1) FROM books WHERE 1=1`+where, args...).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("count books: %w", err)
	}
	return total, nil
}

func (r *Repository) Update(ctx context.Context, id int64, patch models.BookPatch) (models.Book, error) {
	if err := patch.Validate(); err != nil {
		return models.Book{}, err
	}

	current, err := r.GetByID(ctx, id)
	if err != nil {
		return models.Book{}, err
	}

	updated := current
	if patch.Title != nil {
		updated.Title = *patch.Title
	}
	if patch.Author != nil {
		updated.Author = patch.Author
	}
	if patch.Status != nil {
		updated.Status = *patch.Status
	}
	if patch.PriorityToBuy != nil {
		updated.PriorityToBuy = *patch.PriorityToBuy
	}
	if patch.EligibleToSell != nil {
		updated.EligibleToSell = *patch.EligibleToSell
	}
	if patch.Sold != nil {
		updated.Sold = *patch.Sold
	}
	if patch.Notes != nil {
		updated.Notes = patch.Notes
	}

	applyStatusSideEffects(&current, &updated)

	if err := updated.Validate(); err != nil {
		return models.Book{}, err
	}

	_, err = r.db.sql.ExecContext(ctx, `
		UPDATE books SET
			title = ?, author = ?, status = ?, priority_to_buy = ?, eligible_to_sell = ?,
			sold = ?, notes = ?, started_at = ?, finished_at = ?
		WHERE id = ?`,
		updated.Title,
		nullString(updated.Author),
		updated.Status.String(),
		updated.PriorityToBuy,
		updated.EligibleToSell,
		updated.Sold,
		nullString(updated.Notes),
		nullString(updated.StartedAt),
		nullString(updated.FinishedAt),
		id,
	)
	if err != nil {
		return models.Book{}, fmt.Errorf("update book: %w", err)
	}

	return r.GetByID(ctx, id)
}

func (r *Repository) Archive(ctx context.Context, id int64) (models.Book, error) {
	status := models.StatusArchived
	return r.Update(ctx, id, models.BookPatch{Status: &status})
}

func applyStatusSideEffects(before, after *models.Book) {
	if before.Status != models.StatusReading && after.Status == models.StatusReading && after.StartedAt == nil {
		now := models.NowTimestamp()
		after.StartedAt = &now
	}
	if after.Status == models.StatusRead && after.FinishedAt == nil {
		now := models.NowTimestamp()
		after.FinishedAt = &now
	}
}

func (r *Repository) queryBooks(ctx context.Context, query string, args ...any) ([]models.Book, error) {
	rows, err := r.db.sql.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query books: %w", err)
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		book, err := scanBook(rows)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate books: %w", err)
	}
	return books, nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanBook(row rowScanner) (models.Book, error) {
	var book models.Book
	var author, notes, startedAt, finishedAt sql.NullString
	var status string

	err := row.Scan(
		&book.ID,
		&book.Title,
		&author,
		&status,
		&book.PriorityToBuy,
		&book.EligibleToSell,
		&book.Sold,
		&notes,
		&book.AddedAt,
		&startedAt,
		&finishedAt,
	)
	if err != nil {
		return models.Book{}, fmt.Errorf("scan book: %w", err)
	}

	book.Author = nullToPtr(author)
	book.Notes = nullToPtr(notes)
	book.StartedAt = nullToPtr(startedAt)
	book.FinishedAt = nullToPtr(finishedAt)
	book.Status = models.Status(status)

	return book, nil
}

func nullString(v *string) sql.NullString {
	if v == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *v, Valid: true}
}

func nullToPtr(v sql.NullString) *string {
	if !v.Valid {
		return nil
	}
	s := v.String
	return &s
}

func escapeLike(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `%`, `\%`)
	s = strings.ReplaceAll(s, `_`, `\_`)
	return s
}
