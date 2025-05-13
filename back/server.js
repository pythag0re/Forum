const express = require("express")
const connectDB = require("./config/db")
const dotenv = require("dotenv").config();
const port = 5000

//connection a la db
connectDB

var app = express();

// middleware: permet de traiter les données de la request
app.use(express.json())
app.use(express.urlencoded({ extended: false}))

app.use("/post", require("./routes/post.routes"))

app.listen(port, () => console.log("serveur démaré au port " + port))

