package database

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Init(dbPath string) error {
	if DB != nil {
		DB.Close()
	}

	var err error
	DB, err = sql.Open("sqlite", dbPath+"?_txlock=immediate")
	if err != nil {
		return err
	}

	if err = createTables(); err != nil {
		return err
	}

	if err = seedData(); err != nil {
		log.Printf("Warning: seed data error (may be already seeded): %v", err)
	}

	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

func createTables() error {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS kids (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			avatar TEXT NOT NULL,
			points INTEGER DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS behaviors (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			type TEXT NOT NULL CHECK(type IN ('add', 'sub')),
			points INTEGER NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS prizes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			icon TEXT NOT NULL,
			probability REAL NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS records (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			kid_id INTEGER NOT NULL,
			behavior_id INTEGER NOT NULL,
			points INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			reversed BOOLEAN DEFAULT FALSE,
			reversed_reason TEXT,
			reversed_at DATETIME,
			FOREIGN KEY (kid_id) REFERENCES kids(id),
			FOREIGN KEY (behavior_id) REFERENCES behaviors(id)
		)`,
	}

	for _, t := range tables {
		if _, err := DB.Exec(t); err != nil {
			return err
		}
	}
	return nil
}

func seedData() error {
	var count int
	DB.QueryRow("SELECT COUNT(*) FROM kids").Scan(&count)
	if count > 0 {
		return nil
	}

	kids := []struct {
		name   string
		avatar string
		points int
	}{
		{"大宝", "🧒", 10},
		{"二宝", "👧", 10},
	}
	for _, k := range kids {
		_, err := DB.Exec("INSERT INTO kids (name, avatar, points) VALUES (?, ?, ?)", k.name, k.avatar, k.points)
		if err != nil {
			return err
		}
	}

	behaviors := []struct {
		name   string
		type_  string
		points int
	}{
		{"作业得A", "add", 2},
		{"打扫卫生", "add", 3},
		{"帮忙洗碗", "add", 2},
		{"主动读书", "add", 3},
		{"不系安全带", "sub", 2},
		{"打架", "sub", 3},
		{"超时玩手机", "sub", 2},
		{"迟到", "sub", 1},
	}
	for _, b := range behaviors {
		_, err := DB.Exec("INSERT INTO behaviors (name, type, points) VALUES (?, ?, ?)", b.name, b.type_, b.points)
		if err != nil {
			return err
		}
	}

	prizes := []struct {
		name        string
		icon        string
		probability float64
	}{
		{"Switch半小时", "🎮", 3},
		{"宝可梦3局", "🦄", 3},
		{"XBOX", "🧗", 3},
		{"吃大餐", "🍕", 5},
		{"静坐5分钟", "🧘", 20},
		{"动画片30分钟", "📺", 20},
		{"什么都没有", "😀", 15},
		{"诵经|背古诗", "📝", 15},
		{"零食", "🍪", 10},
		{"卡丁车", "🚘", 3},
		{"抓娃娃", "🧸", 3},
	}
	for _, p := range prizes {
		_, err := DB.Exec("INSERT INTO prizes (name, icon, probability) VALUES (?, ?, ?)", p.name, p.icon, p.probability)
		if err != nil {
			return err
		}
	}

	return nil
}
