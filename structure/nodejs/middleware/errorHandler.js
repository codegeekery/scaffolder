// Middleware para rutas no encontradas
function notFound(req, res, next) {
  res.status(404).json({ error: `Ruta no encontrada: ${req.originalUrl}` });
}

// Middleware para manejo general de errores
function errorHandler(err, req, res, next) {
  console.error(err.stack);
  res.status(err.status || 500).json({
    error: err.message || 'Error interno del servidor',
  });
}

module.exports = { notFound, errorHandler };
