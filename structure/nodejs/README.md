# Proyecto Node.js Pro

Estructura básica y profesional con Express.

## Carpetas

```
proyecto-node-pro/
├── config/         # Configuración de la app (puerto, etc.)
├── controllers/    # Lógica de negocio separada de las rutas
├── middleware/      # Middlewares (logger, manejo de errores)
├── routes/          # Definición de rutas
├── utils/           # Funciones auxiliares reutilizables
├── package.json
└── server.js
```

## Uso

```bash
npm install
npm run dev
```

`npm run dev` usa el flag `--watch` nativo de Node.js (v18.11+) para reiniciar el servidor automáticamente al detectar cambios, sin depender de nodemon.

Luego visita `http://localhost:3000/`.
