const express = require("express");
const router = express.Router();

const likeRoutes = require("./like");

router.use("/like", likeRoutes);

module.exports = router;
