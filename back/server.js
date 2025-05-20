const express = require("express")
const connectDB = require("./config/db")
const dotenv = require("dotenv").config();
const cors = require("cors");
const port = process.env.PORT || 5000

//connection a la db
connectDB()

var app = express();

// middleware: permet de traiter les données de la request
app.use(cors());
app.use(express.json())
app.use(express.urlencoded({ extended: false}))

app.use("/post", require("./routes/post.routes"))

// Gestion des erreurs globale
app.use((err, req, res, next) => {
    console.error(err.stack);
    res.status(500).json({ message: "Une erreur est survenue", error: err.message });
});

app.listen(port, () => console.log("serveur démarré au port " + port))

