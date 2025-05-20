const likeService = require('../services/like');

function list(req, res) {
  const likes = likeService.findAll();
  res.status(200).json(likes);
}

function read(req, res) {
  const likeId = req.params.id;
  const like = likeService.find(likeId);
  if (like)
    res.status(200).json(like);
  else
    res.status(404).json({ message: "like non trouvé" });
}

function create(req, res) {
  const datas = req.body;
  const createdlike = likeService.create(datas);
  if (createdlike)
    res.status(201).json({ message: "like créé" });
  else
    res.status(400).json({ message: "Erreur lors de l'insertion" });
}

function update(req, res) {
  const likeId = req.params.id;
  const datas = req.body;
  const updatedlike = likeService.update(likeId, datas);
  if (updatedlike) {
    res.status(200).json({ message: "like édité" });
  } else {
    res.status(400).json({ message: "Erreur lors de l'édition" });
  }
}

function remove(req, res) {
  const likeId = req.params.id;
  const removedlike = likeService.remove(likeId);
  if (removedlike) {
    res.status(200).json({ message: "like supprimé" });
  } else {
    res.status(400).json({ message: "Erreur lors de la suppression" });
  }
}

module.exports = {
  list,
  read,
  create,
  update,
  remove
};
