const express = require("express");
const router = express.Router();

const artistRoutes = require("./like");

router.use("/like", likeRoutes);

module.exports = router;
