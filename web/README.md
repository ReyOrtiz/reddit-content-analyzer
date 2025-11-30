# Reddit Content Analyzer - Web App

A React.js web application for searching relevant subreddit posts by topic.

## Features

- **Topic Search**: Enter a topic to search for relevant Reddit posts
- **Subreddit Management**: Add, remove, and edit subreddits to search
- **Filters**:
  - Limit: Control the number of results (1-100)
  - Relevance Threshold: Set the relevance threshold (0-1)
  - Created After: Filter posts by creation date using a datetime picker

## Getting Started

### Prerequisites

- Node.js (v16 or higher)
- npm or yarn

### Installation

1. Install dependencies:
```bash
npm install
```

2. Start the development server:
```bash
npm run dev
```

The app will be available at `http://localhost:3000`

### Building for Production

```bash
npm run build
```

The built files will be in the `dist` directory.

## Backend Integration

This app connects to the Go backend API running on `http://localhost:8080`. Make sure the backend server is running before using the web app.

The app uses a proxy configuration in `vite.config.js` to forward API requests to the backend.

