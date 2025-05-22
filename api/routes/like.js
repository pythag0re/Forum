const express = require("express");
const router = express.Router();

const likeController = require("../controllers/like");

// Tous les likees
router.get("/", likeController.list);

// Détails d'un likee
router.get("/:id", likeController.read);

// Créer un likee
router.post("/", likeController.create);

// Modifier un likee
router.put("/:id", likeController.updatelike);

// Supprimer un likee
router.delete("/:id", likeController.remove);

module.exports = router;
