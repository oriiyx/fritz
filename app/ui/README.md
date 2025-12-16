# Fritz UI - React Frontend

Modern React frontend for the Fritz PIM system built with TypeScript, TanStack Router, TanStack Query, and Tailwind CSS
with daisyUI.

## Features

- 🔐 **OAuth Authentication** - Google and GitHub login support
- 🎨 **Modern UI** - Built with Tailwind CSS and daisyUI components
- 🚀 **Type-Safe Routing** - TanStack Router with full TypeScript support
- 📊 **Data Management** - TanStack Query for efficient server state management
- 📱 **Responsive Design** - Mobile-first approach with drawer navigation
- 🎯 **State Management** - Zustand for client-side state

## Project Structure

```
src/
├── components/
│   └── layout/          # Layout components (Header, Footer, Sidebar, UserMenu)
├── layouts/             # Page layouts (RootLayout, DashboardLayout)
├── pages/               # Page components
├── services/            # API service layer
├── stores/              # Zustand stores
├── lib/                 # Utilities and configurations
├── router.tsx           # Router configuration
└── App.tsx             # Main app component
```

## Getting Started

### Prerequisites

- Node.js 22.20.0 or higher (use `.tool-versions` or install manually)
- yarn

### Installation

1. Install dependencies:

```bash
yarn install
```

2. Copy the environment file:

```bash
cp .env.example .env
```

3. Update `.env` with your API URL:

```
VITE_API_URL=http://localhost:8080
```

### Development

Start the development server:

```bash
yarn run dev
```

The app will be available at `http://localhost:3333`

### Building for Production

```bash
yarn run build
```

The build output will be in the `dist` directory.

### Preview Production Build

```bash
yarn run preview
```
