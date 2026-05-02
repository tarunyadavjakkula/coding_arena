const express = require('express');
const router = express.Router();
const { getProblem } = require('./problemController');

// GET /api/problems/:id
router.get('/problems/:id', getProblem);

module.exports = router;
