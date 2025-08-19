# React + TypeScript + Vite

This template provides a minimal setup to get React working in Vite with HMR and some ESLint rules.

Currently, two official plugins are available:

- [@vitejs/plugin-react](https://github.com/vitejs/vite-plugin-react/blob/main/packages/plugin-react) uses [Babel](https://babeljs.io/) for Fast Refresh
- [@vitejs/plugin-react-swc](https://github.com/vitejs/vite-plugin-react/blob/main/packages/plugin-react-swc) uses [SWC](https://swc.rs/) for Fast Refresh

# Team Management Frontend

A React TypeScript frontend application that interacts with two Go microservices:
- **User Service**: GraphQL API for user management and team operations
- **Asset Service**: REST API for folder and note management

## Features

### User Management (GraphQL)
- User authentication (login/logout)
- User creation
- Fetch all users

### Team Management (REST)
- Create teams
- View all teams
- Add/remove team members
- Add/remove team managers

### Asset Management (REST)
- Create and manage folders
- Create, view, update, and delete notes
- Share folders and notes with other users
- View sharing permissions

## Technology Stack

- **React 19** with TypeScript
- **Vite** for build tooling
- **Apollo Client** for GraphQL communication
- **Axios** for REST API communication
- **Context API** for state management

## Prerequisites

Make sure your backend services are running:
- User Service: `http://localhost:8090`
- Asset Service: `http://localhost:8080`

## Getting Started

1. Install dependencies:
```bash
npm install
```

2. Start the development server:
```bash
npm run dev
```

3. Open your browser and navigate to `http://localhost:5173`

## API Endpoints

### User Service (GraphQL - http://localhost:8090/user/query)
- `fetchUsers`: Get all users
- `createUser`: Create a new user
- `login`: Authenticate user
- `logout`: Sign out user

### Team Service (REST - http://localhost:8090/teams)
- `GET /teams`: Get all teams
- `POST /teams`: Create a new team
- `GET /teams/:id`: Get team by ID
- `POST /teams/:id/members`: Add team member
- `DELETE /teams/:id/members/:memberId`: Remove team member
- `POST /teams/:id/managers`: Add team manager
- `DELETE /teams/:id/managers/:managerId`: Remove team manager

### Asset Service (REST - http://localhost:8080/api/v1)
- `GET /folders`: Get all folders
- `POST /folders`: Create a new folder
- `GET /folders/:id`: Get folder by ID
- `DELETE /folders/:id`: Delete folder
- `GET /notes`: Get all notes
- `POST /notes`: Create a new note
- `GET /notes/:id`: Get note by ID
- `PUT /notes/:id`: Update note
- `DELETE /notes/:id`: Delete note
- `POST /folders/:id/share`: Share folder
- `POST /notes/:id/share`: Share note

## Authentication

The application uses JWT tokens for authentication. Upon successful login:
- Token is stored in localStorage
- Token is automatically included in API requests
- User session persists across browser refreshes

## Project Structure

```
src/
├── components/          # React components
│   ├── Login.tsx       # User authentication
│   ├── Dashboard.tsx   # Main dashboard
│   ├── Teams.tsx       # Team management
│   ├── Assets.tsx      # Asset management
│   └── Header.tsx      # Navigation header
├── contexts/           # React contexts
│   └── AuthContext.tsx # Authentication state
├── services/           # API service layers
│   ├── api.ts          # Base API configuration
│   ├── userService.ts  # GraphQL user operations
│   ├── teamService.ts  # REST team operations
│   └── assetService.ts # REST asset operations
├── types/              # TypeScript type definitions
│   └── index.ts        # Shared interfaces
├── App.tsx             # Main application component
└── main.tsx            # Application entry point
```

## Usage

1. **Login**: Use your email and password to authenticate
2. **Teams**: Create and manage teams, add/remove members and managers
3. **Assets**: Create folders and notes, organize your content
4. **Sharing**: Share folders and notes with team members

## Development

- Run tests: `npm run test` (when tests are added)
- Build for production: `npm run build`
- Preview production build: `npm run preview`
- Lint code: `npm run lint`

## Notes

- The frontend assumes the backend services are running on localhost
- CORS must be configured on the backend services to allow frontend access
- Make sure JWT middleware is properly configured on protected routes

You can also install [eslint-plugin-react-x](https://github.com/Rel1cx/eslint-react/tree/main/packages/plugins/eslint-plugin-react-x) and [eslint-plugin-react-dom](https://github.com/Rel1cx/eslint-react/tree/main/packages/plugins/eslint-plugin-react-dom) for React-specific lint rules:

```js
// eslint.config.js
import reactX from 'eslint-plugin-react-x'
import reactDom from 'eslint-plugin-react-dom'

export default tseslint.config([
  globalIgnores(['dist']),
  {
    files: ['**/*.{ts,tsx}'],
    extends: [
      // Other configs...
      // Enable lint rules for React
      reactX.configs['recommended-typescript'],
      // Enable lint rules for React DOM
      reactDom.configs.recommended,
    ],
    languageOptions: {
      parserOptions: {
        project: ['./tsconfig.node.json', './tsconfig.app.json'],
        tsconfigRootDir: import.meta.dirname,
      },
      // other options...
    },
  },
])
```
