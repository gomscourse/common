package prettier

import (
	"fmt"
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"
)

func TestQueryPrettier(t *testing.T) {
	t.Parallel()

	var (
		id       = gofakeit.Int64()
		author   = gofakeit.FirstName()
		safeName = fmt.Sprintf("\"%s\"", author)
		chatID   = gofakeit.Int64()
	)

	t.Run("without squirrel", func(t *testing.T) {
		t.Parallel()
		t.Run("query with no params", func(t *testing.T) {
			t.Parallel()

			q := "SELECT id, author, chat_id FROM message"
			pretty := Pretty(q, PlaceholderDollar)
			require.Equal(t, q, pretty)
		})

		t.Run("query with spaces", func(t *testing.T) {
			t.Parallel()

			q := "   SELECT id, author, chat_id FROM message   "
			pretty := Pretty(q, PlaceholderDollar)
			require.Equal(t, "SELECT id, author, chat_id FROM message", pretty)
		})

		t.Run("select by one param", func(t *testing.T) {
			t.Parallel()

			q := "SELECT id, author, chat_id FROM message WHERE id = $1"
			pretty := Pretty(q, PlaceholderDollar, id)
			require.Equal(t, fmt.Sprintf("SELECT id, author, chat_id FROM message WHERE id = %v", id), pretty)
		})

		t.Run("select with IN clause", func(t *testing.T) {
			t.Parallel()

			q := "SELECT id, author, chat_id FROM message WHERE id IN ($1,$2,$3)"
			pretty := Pretty(q, PlaceholderDollar, id, 123, 321)
			require.Equal(t, fmt.Sprintf("SELECT id, author, chat_id FROM message WHERE id IN (%v,%v,%v)", id, 123, 321), pretty)
		})

		t.Run("insert", func(t *testing.T) {
			t.Parallel()

			q := "INSERT INTO message (author, chat_id) VALUES ($1,$2)"
			pretty := Pretty(q, PlaceholderDollar, author, chatID)
			require.Equal(t, fmt.Sprintf("INSERT INTO message (author, chat_id) VALUES (%v,%v)", safeName, chatID), pretty)
		})

		t.Run("update", func(t *testing.T) {
			t.Parallel()

			q := "UPDATE message SET author = $1, chat_id = $2 WHERE author = $1"
			pretty := Pretty(q, PlaceholderDollar, []byte(author), chatID)
			require.Equal(t, fmt.Sprintf("UPDATE message SET author = %v, chatID = %v WHERE author = %v", safeName, chatID, safeName), pretty)
		})
	})

	t.Run("with squirrel", func(t *testing.T) {
		t.Parallel()
		t.Run("query with no params", func(t *testing.T) {
			t.Parallel()

			builder := sq.Select("id, author, chatID").
				PlaceholderFormat(sq.Dollar).
				From("message")

			query, v, err := builder.ToSql()
			require.NoError(t, err)

			pretty := Pretty(query, PlaceholderDollar, v...)
			require.Equal(t, "SELECT id, author, chat_id FROM message", pretty)
		})

		t.Run("query with spaces", func(t *testing.T) {
			t.Parallel()

			builder := sq.Select("id, author, chatID").
				PlaceholderFormat(sq.Dollar).
				From("message    ")

			query, v, err := builder.ToSql()
			require.NoError(t, err)

			pretty := Pretty(query, PlaceholderDollar, v...)

			require.Equal(t, "SELECT id, author, chat_id FROM message", pretty)
		})

		t.Run("SELECT by one param", func(t *testing.T) {
			t.Parallel()

			builder := sq.Select("id, author, chatID").
				PlaceholderFormat(sq.Dollar).
				From("message").
				Where(sq.Eq{"id": id})

			query, v, err := builder.ToSql()
			require.NoError(t, err)

			pretty := Pretty(query, PlaceholderDollar, v...)
			require.Equal(t, fmt.Sprintf("SELECT id, author, chat_id FROM message WHERE id = %v", id), pretty)
		})

		t.Run("select with IN clause", func(t *testing.T) {
			t.Parallel()

			builder := sq.Select("id, author, chatID").
				PlaceholderFormat(sq.Dollar).
				From("message").
				Where(sq.Eq{"id": []int64{id, 123, 321}})

			query, v, err := builder.ToSql()
			require.NoError(t, err)

			pretty := Pretty(query, PlaceholderDollar, v...)
			require.Equal(t, fmt.Sprintf("SELECT id, author, chat_id FROM message WHERE id IN (%v,%v,%v)", id, 123, 321), pretty)
		})

		t.Run("insert", func(t *testing.T) {
			t.Parallel()

			builder := sq.Insert("message").
				PlaceholderFormat(sq.Dollar).
				Columns("author, chatID").
				Values(author, chatID)

			query, v, err := builder.ToSql()
			require.NoError(t, err)

			pretty := Pretty(query, PlaceholderDollar, v...)
			require.Equal(t, fmt.Sprintf("INSERT INTO message (author, chat_id) VALUES (%v,%v)", safeName, chatID), pretty)
		})

		t.Run("update", func(t *testing.T) {
			t.Parallel()

			builder := sq.Update("message").
				PlaceholderFormat(sq.Dollar).
				Set("author", author).
				Set("chatID", chatID).
				Where(sq.Eq{"author": author})

			query, v, err := builder.ToSql()
			require.NoError(t, err)

			pretty := Pretty(query, PlaceholderDollar, v...)
			require.Equal(t, fmt.Sprintf("UPDATE message SET author = %v, chat_id = %v WHERE author = %v", safeName, chatID, safeName), pretty)
		})
	})
}
