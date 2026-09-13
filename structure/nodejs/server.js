const express = require('express');

const config = require('./config/config');
const logger = require('./middleware/logger');
const { notFound, errorHandler } = require('./middleware/errorHandler');
const routes = require('./routes');

const app = express();

app.use(logger);
app.use('/', routes);

app.use(notFound);
app.use(errorHandler);

app.listen(config.port, () => {
  console.log(`Servidor corriendo en http://localhost:${config.port}`);
});
