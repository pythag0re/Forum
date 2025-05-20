const express = require('express');
const app = express();

const hostname = '127.0.0.1';
const port = process.env.PORT || 3000;

const routes = require("./routes");

app.use("/api", routes);

app.listen(port, hostname, () => {
	console.log(`Serveur démarré sur http://${hostname}:${port}`);
});
