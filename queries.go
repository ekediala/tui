package main

const CreateTableStmt = `
	CREATE TABLE IF NOT EXISTS notes (
		id INTEGER NOT NULL PRIMARY KEY,
		title TEXT NOT NULL,
		body TEXT NOT NULL
	);
`

const GetNotesStmt = `
	SELECT id, title, body FROM notes ORDER BY id DESC;
`

const UpsertNoteStmt = `
	INSERT INTO notes (id, title, body)
	VALUES(?, ?, ?)
	ON CONFLICT(id) DO UPDATE
	SET title=excluded.title, body=excluded.body;
`
