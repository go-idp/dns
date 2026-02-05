# DNS CLI Documentation

This directory contains the documentation website for DNS CLI, built with VitePress.

## Development

### Prerequisites

- Node.js 20+
- pnpm 8+

### Install Dependencies

```bash
pnpm install
```

### Development Server

```bash
pnpm docs:dev
```

Visit `http://localhost:5173` to view the documentation.

### Build

```bash
pnpm docs:build
```

### Preview Build

```bash
pnpm docs:preview
```

## Deployment

The documentation is automatically deployed to GitHub Pages via GitHub Actions when changes are pushed to the `master` branch.

The workflow is defined in `.github/workflows/docs.yml`.

## Project Structure

```
docs/
├── .vitepress/          # VitePress configuration
│   └── config.mts       # VitePress config file
├── guide/               # User guide documentation
├── api/                 # API reference documentation
├── examples/            # Code examples
├── index.md             # Homepage
├── package.json         # npm/pnpm configuration
├── tsconfig.json        # TypeScript configuration
└── .npmrc               # npm configuration
```

## Adding New Documentation

1. Create a new markdown file in the appropriate directory (`guide/`, `api/`, or `examples/`)
2. Add the page to the sidebar in `.vitepress/config.mts`
3. Commit and push to trigger automatic deployment
