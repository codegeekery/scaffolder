// Middleware que registra cada petición en consola
function logger(req, res, next) {
  console.log(`[${new Date().toISOString()}] ${req.method} ${req.originalUrl}`);
  next();
}

module.exports = logger;
