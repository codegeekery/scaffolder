# Node.js Pro Project

Basic, professional structure using Express.

## Folders

```
proyecto-node-pro/
├── config/         # App configuration (port, etc.)
├── controllers/    # Business logic separated from routes
├── middleware/      # Middlewares (logger, error handling)
├── routes/          # Route definitions
├── utils/           # Reusable helper functions
├── package.json
└── server.js
```

## Usage

```bash
npm install
npm run dev
```

`npm run dev` uses Node.js's native `--watch` flag (v18.11+) to automatically restart the server on file changes, without relying on nodemon.

Then visit `http://localhost:3000/`.
