const express = require('express');
const app = express();
const cors = require('cors');
const sqlite3 = require('sqlite3');


app.use(cors());
app.use(express.json());

const db = new sqlite3.Database('../app.db');

// pour add un like
app.post('/like', (request, resultat) => {
  const { user_id, post_id, comment_id } = request.body;

  const stmt = db.prepare(
    `INSERT INTO likes (user_id, post_id, comment_id) VALUES (?, ?, ?)`
  );

  stmt.run([user_id, post_id || null, comment_id || null], function (err) {
    if (err) {
      return resultat.status(500).json({ error: err.message });
    }
    resultat.status(201).json({ id: this.lastID, message: 'like ajouté' });
  });
});

// pour le delete
app.delete('/like', (request, resultat) => {
  const { user_id, post_id, comment_id } = request.body;

  const stmt = db.prepare(
    `DELETE FROM likes WHERE user_id = ? AND post_id IS ? AND comment_id IS ?`
  );

  stmt.run([user_id, post_id || null, comment_id || null], function (err) {
    if (err) {
      return resultat.status(500).json({ error: err.message });
    }
    resultat.status(200).json({ message: 'Like retiré' });
  });
});

// pour le noombre de like
app.get('/like/count', (request, resultat) => {
  const { post_id, comment_id } = request.query;

  const query = `
    SELECT COUNT(*) as count
    FROM likes
    WHERE post_id IS ? AND comment_id IS ?
  `;

  db.get(query, [post_id || null, comment_id || null], (err, row) => {
    if (err) {
      return resultat.status(500).json({ error: err.message });
    }
    resultat.json({ count: row.count });
  });
});

app.get('/', (req, res) => {
  res.send('API LIKE PRETE');
});
app.listen(3001, () => {
  console.log('API Likes SQLite démarrée sur http://localhost:3001');
});

app.get('/like/has-liked', (req, res) => {
  const { user_id, post_id } = req.query;

  db.get(
    `SELECT COUNT(*) as count FROM likes WHERE user_id = ? AND post_id = ?`,
    [user_id, post_id],
    (err, row) => {
      if (err) return res.status(500).json({ error: err.message });
      res.json({ liked: row.count > 0 });
    }
  );
});