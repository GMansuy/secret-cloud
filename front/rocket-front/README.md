# Rocket-Front

A modern web application for creating and managing Kubernetes clusters with Microsoft OAuth authentication.

## Overview

Rocket-Front is a Next.js web application built with TypeScript and React that provides a user interface for creating and managing Kubernetes clusters. The application features Microsoft OAuth authentication and communicates with a backend API for cluster operations.

## Features

- **Secure Authentication**: Microsoft Azure AD integration
- **Cluster Management**: Create and configure Kubernetes clusters
- **User-Friendly Interface**: Simple form-based cluster creation
- **Responsive Design**: Works on various device sizes

## Application Structure
```markdown


src/app/
├── auth-config.tsx       # Authentication configuration with Microsoft OAuth
├── page.tsx              # Root page (redirects to login)
├── login/
│   └── page.tsx          # Microsoft login page
├── callback/
│   └── page.tsx          # OAuth callback handler
└── cluster/
└── page.tsx          # Cluster creation interface
```

## Authentication Flow

1. User accesses the application and is redirected to the login page
2. User clicks "Login" to authenticate with Microsoft
3. After successful authentication, the callback page:
   - Processes the authentication response
   - Stores the access token in localStorage
   - Redirects to the cluster page

## Cluster Management

The application allows users to create Kubernetes clusters with customizable parameters:
- Cluster name
- Control plane machine count
- Worker machine count

## API Integration

The application communicates with a backend API at `http://localhost:8080`:

```
POST /cluster
{
"name": "my-cluster",
"controlplaneMachineCount": 1,
"workerMachineCount": 2
}
```

## Getting Started

### Prerequisites

- Node.js (v16+)
- npm or yarn
- Azure AD application configured for authentication

### Installation

```bash
# Clone the repository
git clone <repository-url>
cd rocket-front

# Install dependencies
npm install

# Start development server
npm run dev
```

Open [http://localhost:3000](http://localhost:3000) with your browser to see the application.

### Configuration

Authentication settings are configured in `src/app/auth-config.tsx`:
- `authority`: Microsoft identity platform URL
- `client_id`: Azure AD application ID
- `redirect_uri`: Callback URL after authentication
- `scope`: Required OAuth scopes

## Troubleshooting

### Common Issues

- **Authentication Errors**: Verify Azure AD configuration in `auth-config.tsx`
- **API Connection Errors**: Ensure backend API is running at `http://localhost:8080`
- **404 Errors**: Check that all required Next.js pages are properly defined

## Technologies

- [Next.js](https://nextjs.org/)
- [React](https://reactjs.org/)
- [TypeScript](https://www.typescriptlang.org/)
- [oidc-client-ts](https://github.com/authts/oidc-client-ts) for authentication
```