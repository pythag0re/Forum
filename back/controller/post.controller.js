const postModel = require('../model/post.model')

module.exports.setPosts = async(req, res) => {
    if (!req.body.message) {
        res.status(400).json({message: "ajoute un message merci"})
    } else {
        res.json({message: "message bien envoyé..."})
    }
}