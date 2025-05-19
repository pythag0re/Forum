const express = require("express");
const { setPosts } = require("../controller/post.controller");
const router = express.Router();

router.get("/", (req, res) => {
    res.json({message: "données"})
})

router.post("/", setPosts);

router.put('/:id', (req, res) => {
    res.json({messageId: req.params.id})
})

router.delete("/:id", (req, res) => {
    res.json({message: "poste supprimé id: " + req.params.id})
})

router.patch("/like-post/:id" , (req, res) => {
    res.json({message: "post liké: id: " + req.params.id })
})

router.patch("/dislike-post/:id" , (req, res) => {
    res.json({message: "post disliké: id: " + req.params.id })
})

module.exports = router