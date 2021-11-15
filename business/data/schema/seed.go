package schema

import "github.com/jmoiron/sqlx"

func Seed(db *sqlx.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	if _, err := tx.Exec(seeds); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	return tx.Commit()
}

const seeds = `
INSERT INTO users (user_id, name, email, roles, password_hash, date_created, date_updated) VALUES
	('5cf37266-3473-4006-984f-9325122678b7', 'Admin Gopher', 'admin@example.com', '{ADMIN,USER}', '$2a$10$1ggfMVZV6Js0ybvJufLRUOWHS5f6KneuP0XwwHpJ8L8ipdry9f2/a', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('45b5fbd3-755f-4379-8f07-a58d4a30fa2f', 'User Gopher', 'user@example.com', '{USER}', '$2a$10$9/XASPKBbJKVfCAZKDH.UuhsuALDr5vVm6VrYA9VFR8rccK86C1hW', '2019-03-24 00:00:00', '2019-03-24 00:00:00')
	ON CONFLICT DO NOTHING;

INSERT INTO supermarkets (supermarket_id, name, address, user_id, date_created, date_updated) VALUES
	('a2b0639f-2cc6-44b8-b97b-15d69dbb511e', 'Magnum', 'Nur-Sultan avenue', '45b5fbd3-755f-4379-8f07-a58d4a30fa2f', '2019-01-01 00:00:01.000001+00', '2019-01-01 00:00:01.000001+00'),
	('72f8b983-3eb4-48db-9ed0-e45cc6bd716b', 'Small', 'Kabanbay batura','45b5fbd3-755f-4379-8f07-a58d4a30fa2f', '2019-01-01 00:00:02.000001+00', '2019-01-01 00:00:02.000001+00')
	ON CONFLICT DO NOTHING;

INSERT INTO categories (category_id, name, supermarket_id, date_created, date_updated) VALUES
	('daaae037-3b1c-44fe-8922-106430fba304', 'fruit', 'a2b0639f-2cc6-44b8-b97b-15d69dbb511e', '2019-01-01 00:00:01.000001+00', '2019-01-01 00:00:01.000001+00'),
	('c62961a6-ed9c-4277-9bd7-8b985bf15321', 'cans', '72f8b983-3eb4-48db-9ed0-e45cc6bd716b', '2019-01-01 00:00:02.000001+00', '2019-01-01 00:00:02.000001+00')
	ON CONFLICT DO NOTHING;

INSERT INTO goods (good_id, name, price, category_id, date_created, date_updated) VALUES
	('fcfd133c-8384-4c06-9f2e-33d65d764294', 'apple', 1234, 'daaae037-3b1c-44fe-8922-106430fba304', '2019-01-01 00:00:01.000001+00', '2019-01-01 00:00:01.000001+00'),
	('08a43bcd-527c-4af8-8042-d81ab299c7e7', 'peach', 12345, 'daaae037-3b1c-44fe-8922-106430fba304', '2019-01-01 00:00:02.000001+00', '2019-01-01 00:00:02.000001+00')
	ON CONFLICT DO NOTHING;`

const deleteAll = `
DELETE FROM goods;
DELETE FROM categories;
DELETE FROM supermarkets;
DELETE FROM users;`

func DeleteAll(db *sqlx.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(deleteAll); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}
