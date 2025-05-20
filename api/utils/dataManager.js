const sqlite3 = require('sqlite3').verbose();
const path = require('path');


const dbPath = path.resolve(__dirname, '../../app.db');

let db = new sqlite3.Database(dbPath, sqlite3.OPEN_READWRITE, (err) => {
    if (err) {
        console.error("Erreur :", err.message);
    } else {
        console.log('Connecté à la base de données SQLite.');
    }
});

// Exécuter une requête SELECT
db.serialize(() => {
    db.each(`SELECT email, password, pseudo FROM Users`, (err, row) => {
        if (err) {
            console.error(err.message);
        }
        console.log(`${row.email}:${row.pseudo}`);
    });
});

function updateLike() {
    const createdAt = new Date().toISOString();
    db.run(
        `INSERT INTO likes (id, user_id, post_id, comment_id, created_at) VALUES (?, ?, ?, ?, ?)`,
        ['1', '1', '5', 'malik', createdAt],
        function (err) {
            if (err) {
                console.error("Erreur d'insertion :", err.message);
            } else {
                console.log(`Like inséré avec l'id ${this.lastID}`);
            }
        }
    );
}
updateLike();

// Fermer la base de données
db.close((err) => {
    if (err) {
        console.error(err.message);
    } else {
        console.log('Fermeture de la base de données.');
    }
});




module.exports = { db }