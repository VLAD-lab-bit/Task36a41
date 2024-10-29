package storage

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"Task36a41/pkg/rss"

	_ "github.com/lib/pq" // PostgreSQL драйвер
)

// Storage представляет структуру для работы с БД.
type Storage struct {
	db *sql.DB
}

// MockStorage - имитация хранилища для тестирования
type MockStorage struct {
	posts []rss.Post
}

// NewMockStorage создает новый экземпляр MockStorage с заданными постами
func NewMockStorage(posts []rss.Post) *MockStorage {
	return &MockStorage{posts: posts}
}

// GetLastNPosts возвращает последние N публикаций из имитационного хранилища
func (m *MockStorage) GetLastNPosts(n int) ([]rss.Post, error) {
	if n > len(m.posts) {
		n = len(m.posts)
	}
	return m.posts[:n], nil
}

// New создает новое подключение к базе данных.
func New(connectionString string) (*Storage, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %v", err)
	}

	// Проверим соединение
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("could not ping database: %v", err)
	}

	return &Storage{db: db}, nil
}

// Close закрывает соединение с базой данных.
func (s *Storage) Close() error {
	return s.db.Close()
}

// SavePost сохраняет одну публикацию в БД.
func (s *Storage) SavePost(post rss.Post) error {
	query := `
		INSERT INTO posts (title, content, pub_time, link)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (link) DO NOTHING
	`

	// Попробуем несколько форматов для парсинга даты
	var pubTime time.Time
	var err error
	timeFormats := []string{
		time.RFC1123Z,
		time.RFC1123,
		"Mon, 2 Jan 2006 15:04:05 -0700",
		"Mon, 2 Jan 2006 15:04:05 MST",
	}

	for _, format := range timeFormats {
		pubTime, err = time.Parse(format, post.PubDate)
		if err == nil {
			break
		}
	}

	if err != nil {
		return fmt.Errorf("could not parse publication time: %v", err)
	}

	unixTime := pubTime.Unix()

	// Выполняем вставку в базу данных
	_, err = s.db.Exec(query, post.Title, post.Content, unixTime, post.Link)
	if err != nil {
		return fmt.Errorf("could not insert post: %v", err)
	}

	return nil
}

// SavePosts сохраняет несколько публикаций в БД.
func (s *Storage) SavePosts(posts []rss.Post) error {
	for _, post := range posts {
		if err := s.SavePost(post); err != nil {
			log.Printf("Error saving post: %v", err)
		}
	}
	return nil
}

// GetLastNPosts возвращает последние N публикаций.
func (s *Storage) GetLastNPosts(n int) ([]rss.Post, error) {
	query := `
		SELECT title, content, pub_time, link
		FROM posts
		ORDER BY pub_time DESC
		LIMIT $1
	`

	rows, err := s.db.Query(query, n)
	if err != nil {
		return nil, fmt.Errorf("could not get posts: %v", err)
	}
	defer rows.Close()

	var posts []rss.Post
	for rows.Next() {
		var post rss.Post
		var pubTime int64

		// Извлекаем pub_time как Unix timestamp
		if err := rows.Scan(&post.Title, &post.Content, &pubTime, &post.Link); err != nil {
			return nil, fmt.Errorf("could not scan post: %v", err)
		}

		// Преобразуем Unix timestamp обратно в строку в формате RFC1123Z
		post.PubDate = time.Unix(pubTime, 0).Format(time.RFC1123Z)
		posts = append(posts, post)
	}

	return posts, nil
}
