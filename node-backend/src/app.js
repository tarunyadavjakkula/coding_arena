const express = require('express');
const problemRoutes = require('./problemRoutes');

const app = express();
app.use(express.json());

// Mount the problem routes
app.use('/api', problemRoutes);

module.exports = app;
