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

var ErrAmbiguousTitle = errors.New("ambiguous title")

type ListFilter struct {
	Status          *models.Status
	Category        *models.Category
	PriorityToBuy   *bool
	EligibleToSell  *bool
	IncludeArchived bool
	Pagination      *models.Pagination
	Sort            models.Sort
}

type SearchFilter struct {
	Terms           []string
	Author          string
	Category        *models.Category
	IncludeArchived bool
	Pagination      *models.Pagination
	Sort            models.Sort
}

type CheckFilter struct {
	Title           string
	Author          string
	Exact           bool
	IncludeArchived bool
}

type TitleFilter struct {
	Title           string
	Author          string
	Exact           bool
	IncludeArchived bool
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
			title, author, category, status, priority_to_buy, eligible_to_sell, sold,
			notes, description, added_at, started_at, finished_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		book.Title,
		nullString(book.Author),
		nullCategory(book.Category),
		book.Status.String(),
		book.PriorityToBuy,
		book.EligibleToSell,
		book.Sold,
		nullString(book.Notes),
		nullString(book.Description),
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
		SELECT id, title, author, category, status, priority_to_buy, eligible_to_sell, sold,
		       notes, description, added_at, started_at, finished_at
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

func (r *Repository) GetByTitle(ctx context.Context, filter TitleFilter) (models.Book, error) {
	where, args := buildCheckWhere(CheckFilter{
		Title:           filter.Title,
		Author:          filter.Author,
		Exact:           filter.Exact,
		IncludeArchived: filter.IncludeArchived,
	})
	total, err := r.countBooks(ctx, where, args)
	if err != nil {
		return models.Book{}, err
	}
	switch total {
	case 0:
		return models.Book{}, ErrNotFound
	case 1:
		orderBy, err := models.DefaultSort().OrderByClause()
		if err != nil {
			return models.Book{}, err
		}
		query := `
		SELECT id, title, author, category, status, priority_to_buy, eligible_to_sell, sold,
		       notes, description, added_at, started_at, finished_at
		FROM books WHERE 1=1` + where + orderBy + ` LIMIT 1`
		books, err := r.queryBooks(ctx, query, args...)
		if err != nil {
			return models.Book{}, err
		}
		if len(books) == 0 {
			return models.Book{}, ErrNotFound
		}
		return books[0], nil
	default:
		return models.Book{}, fmt.Errorf("%w: %d matches (use --author or id)", ErrAmbiguousTitle, total)
	}
}

func (r *Repository) List(ctx context.Context, filter ListFilter) (BooksResult, error) {
	where, args := buildListWhere(filter)
	return r.queryBooksPage(ctx, where, args, filter.Pagination, filter.Sort)
}

func (r *Repository) Count(ctx context.Context, filter ListFilter) (int, error) {
	where, args := buildListWhere(filter)
	return r.countBooks(ctx, where, args)
}

func (r *Repository) Stats(ctx context.Context, year int) (models.LibraryStats, error) {
	stats := models.LibraryStats{
		Year:       year,
		ByStatus:   make(map[string]int),
		ByCategory: make(map[string]int),
	}

	archived := models.StatusArchived.String()

	byStatus, err := scanGroupedCounts(ctx, r.db.sql, `
		SELECT status, COUNT(1) FROM books WHERE status != ? GROUP BY status`, archived)
	if err != nil {
		return models.LibraryStats{}, fmt.Errorf("stats by status: %w", err)
	}
	stats.ByStatus = byStatus

	byCategory, err := scanGroupedCounts(ctx, r.db.sql, `
		SELECT category, COUNT(1) FROM books
		WHERE status != ? AND category IS NOT NULL
		GROUP BY category`, archived)
	if err != nil {
		return models.LibraryStats{}, fmt.Errorf("stats by category: %w", err)
	}
	stats.ByCategory = byCategory

	yearStart := fmt.Sprintf("%04d-01-01T00:00:00Z", year)
	yearEnd := fmt.Sprintf("%04d-01-01T00:00:00Z", year+1)
	err = r.db.sql.QueryRowContext(ctx, `
		SELECT COUNT(1) FROM books
		WHERE status != ? AND finished_at IS NOT NULL
		  AND finished_at >= ? AND finished_at < ?`,
		archived, yearStart, yearEnd,
	).Scan(&stats.FinishedThisYear)
	if err != nil {
		return models.LibraryStats{}, fmt.Errorf("finished this year: %w", err)
	}

	toBuy := models.StatusToBuy.String()
	err = r.db.sql.QueryRowContext(ctx, `
		SELECT COUNT(1) FROM books
		WHERE status = ? AND priority_to_buy = 1`, toBuy,
	).Scan(&stats.PriorityWishlist)
	if err != nil {
		return models.LibraryStats{}, fmt.Errorf("priority wishlist: %w", err)
	}

	return stats, nil
}

func scanGroupedCounts(ctx context.Context, db *sql.DB, query string, args ...any) (map[string]int, error) {
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := make(map[string]int)
	for rows.Next() {
		var key string
		var count int
		if err := rows.Scan(&key, &count); err != nil {
			return nil, fmt.Errorf("scan grouped count: %w", err)
		}
		counts[key] = count
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate grouped counts: %w", err)
	}
	return counts, nil
}

func (r *Repository) Search(ctx context.Context, filter SearchFilter) (BooksResult, error) {
	where, args := buildSearchWhere(filter)
	return r.queryBooksPage(ctx, where, args, filter.Pagination, filter.Sort)
}

func (r *Repository) Check(ctx context.Context, filter CheckFilter) (BooksResult, error) {
	where, args := buildCheckWhere(filter)
	return r.queryBooksPage(ctx, where, args, nil, models.DefaultSort())
}

func (r *Repository) queryBooksPage(ctx context.Context, where string, args []any, pagination *models.Pagination, sort models.Sort) (BooksResult, error) {
	total, err := r.countBooks(ctx, where, args)
	if err != nil {
		return BooksResult{}, err
	}

	orderBy, err := sort.OrderByClause()
	if err != nil {
		return BooksResult{}, err
	}

	query := `
		SELECT id, title, author, category, status, priority_to_buy, eligible_to_sell, sold,
		       notes, description, added_at, started_at, finished_at
		FROM books WHERE 1=1` + where + orderBy

	queryArgs := append([]any{}, args...)
	if pagination != nil && pagination.Enabled() {
		query += ` LIMIT ? OFFSET ?`
		queryArgs = append(queryArgs, pagination.Limit, pagination.Offset())
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
	if filter.Category != nil {
		query += ` AND category = ?`
		args = append(args, filter.Category.String())
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
	if terms := normalizeSearchTerms(filter.Terms); len(terms) > 0 {
		clauses := make([]string, len(terms))
		for i, term := range terms {
			pattern := "%" + escapeLike(strings.ToLower(term)) + "%"
			clauses[i] = `(LOWER(title) LIKE ? ESCAPE '\' OR LOWER(COALESCE(description, '')) LIKE ? ESCAPE '\' OR LOWER(COALESCE(author, '')) LIKE ? ESCAPE '\')`
			args = append(args, pattern, pattern, pattern)
		}
		query += ` AND (` + strings.Join(clauses, ` OR `) + `)`
	}
	if a := strings.TrimSpace(filter.Author); a != "" {
		query += ` AND LOWER(COALESCE(author, '')) LIKE ? ESCAPE '\'`
		args = append(args, "%"+escapeLike(strings.ToLower(a))+"%")
	}
	if filter.Category != nil {
		query += ` AND category = ?`
		args = append(args, filter.Category.String())
	}
	return query, args
}

func buildCheckWhere(filter CheckFilter) (string, []any) {
	query := ""
	args := []any{}

	if !filter.IncludeArchived {
		query += ` AND status != ?`
		args = append(args, models.StatusArchived.String())
	}
	title := strings.TrimSpace(filter.Title)
	if title != "" {
		if filter.Exact {
			query += ` AND LOWER(title) = LOWER(?)`
			args = append(args, title)
		} else {
			pattern := "%" + escapeLike(strings.ToLower(title)) + "%"
			query += ` AND LOWER(title) LIKE ? ESCAPE '\'`
			args = append(args, pattern)
		}
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

type querier interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func (r *Repository) Update(ctx context.Context, id int64, patch models.BookPatch) (models.Book, error) {
	return updateWithQuerier(ctx, r.db.sql, id, patch)
}

func (r *Repository) UpdateMany(ctx context.Context, ids []int64, patch models.BookPatch) ([]models.Book, error) {
	if err := patch.Validate(); err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, fmt.Errorf("at least one id is required")
	}

	tx, err := r.db.sql.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	updated := make([]models.Book, 0, len(ids))
	for _, id := range ids {
		book, err := updateWithQuerier(ctx, tx, id, patch)
		if err != nil {
			return nil, err
		}
		updated = append(updated, book)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}
	return updated, nil
}

func updateWithQuerier(ctx context.Context, q querier, id int64, patch models.BookPatch) (models.Book, error) {
	if err := patch.Validate(); err != nil {
		return models.Book{}, err
	}

	current, err := getByIDWithQuerier(ctx, q, id)
	if err != nil {
		return models.Book{}, err
	}

	updated := applyPatch(current, patch)
	if err := updated.Validate(); err != nil {
		return models.Book{}, err
	}

	_, err = q.ExecContext(ctx, `
		UPDATE books SET
			title = ?, author = ?, category = ?, status = ?, priority_to_buy = ?, eligible_to_sell = ?,
			sold = ?, notes = ?, description = ?, started_at = ?, finished_at = ?
		WHERE id = ?`,
		updated.Title,
		nullString(updated.Author),
		nullCategory(updated.Category),
		updated.Status.String(),
		updated.PriorityToBuy,
		updated.EligibleToSell,
		updated.Sold,
		nullString(updated.Notes),
		nullString(updated.Description),
		nullString(updated.StartedAt),
		nullString(updated.FinishedAt),
		id,
	)
	if err != nil {
		return models.Book{}, fmt.Errorf("update book: %w", err)
	}

	return getByIDWithQuerier(ctx, q, id)
}

func getByIDWithQuerier(ctx context.Context, q querier, id int64) (models.Book, error) {
	row := q.QueryRowContext(ctx, `
		SELECT id, title, author, category, status, priority_to_buy, eligible_to_sell, sold,
		       notes, description, added_at, started_at, finished_at
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

func applyPatch(current models.Book, patch models.BookPatch) models.Book {
	updated := current
	if patch.Title != nil {
		updated.Title = *patch.Title
	}
	if patch.Author != nil {
		updated.Author = patch.Author
	} else if patch.ClearAuthor {
		updated.Author = nil
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
	if patch.Description != nil {
		updated.Description = patch.Description
	}
	if patch.Category != nil {
		updated.Category = patch.Category
	} else if patch.ClearCategory {
		updated.Category = nil
	}
	if patch.StartedAt != nil {
		updated.StartedAt = patch.StartedAt
	} else if patch.ClearStartedAt {
		updated.StartedAt = nil
	}
	if patch.FinishedAt != nil {
		updated.FinishedAt = patch.FinishedAt
	} else if patch.ClearFinishedAt {
		updated.FinishedAt = nil
	}
	return updated
}

func (r *Repository) Delete(ctx context.Context, id int64) (models.Book, error) {
	book, err := r.GetByID(ctx, id)
	if err != nil {
		return models.Book{}, err
	}

	_, err = r.db.sql.ExecContext(ctx, `DELETE FROM books WHERE id = ?`, id)
	if err != nil {
		return models.Book{}, fmt.Errorf("delete book: %w", err)
	}

	return book, nil
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
	var author, category, notes, description, startedAt, finishedAt sql.NullString
	var status string

	err := row.Scan(
		&book.ID,
		&book.Title,
		&author,
		&category,
		&status,
		&book.PriorityToBuy,
		&book.EligibleToSell,
		&book.Sold,
		&notes,
		&description,
		&book.AddedAt,
		&startedAt,
		&finishedAt,
	)
	if err != nil {
		return models.Book{}, fmt.Errorf("scan book: %w", err)
	}

	book.Author = nullToPtr(author)
	book.Category = nullToCategory(category)
	book.Notes = nullToPtr(notes)
	book.Description = nullToPtr(description)
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

func nullCategory(v *models.Category) sql.NullString {
	if v == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: v.String(), Valid: true}
}

func nullToCategory(v sql.NullString) *models.Category {
	if !v.Valid {
		return nil
	}
	category := models.Category(v.String)
	return &category
}

func normalizeSearchTerms(terms []string) []string {
	out := make([]string, 0, len(terms))
	for _, term := range terms {
		if term = strings.TrimSpace(term); term != "" {
			out = append(out, term)
		}
	}
	return out
}

func escapeLike(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `%`, `\%`)
	s = strings.ReplaceAll(s, `_`, `\_`)
	return s
}

