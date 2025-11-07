# Deployment Preview System

A lightweight deployment preview system that automatically creates isolated Docker container environments for each Git branch via webhooks.

## Overview

This system listens for Git webhook events (e.g., from GitHub, GitLab) and automatically spins up preview environments in Docker containers. Each branch gets its own isolated nginx container, making it easy to review and test changes before merging.

## Features

- **Webhook-driven**: Automatically triggered by Git push events
- **Branch-based previews**: Each branch gets its own isolated container
- **Automatic cleanup**: Removes existing containers before creating new ones to avoid conflicts
- **Dynamic port allocation**: Each preview environment gets a unique port
- **Public URL exposure**: Uses ngrok to make webhook endpoints accessible
- **Container management**: Filters, stops, and removes existing containers by name

## Architecture

```
Git Repository → Webhook → Gin Server → Docker Client → Preview Container
                              ↓
                           ngrok (Public URL)
```

## Prerequisites

- Go 1.25.3 or higher
- Docker installed and running
- ngrok account (for webhook exposure)

## Installation

1. Clone the repository:
```bash
git clone https://github.com/shoiam/deployment-preview-system.git
cd deployment-preview-system
```

2. Install dependencies:
```bash
go mod download
```

3. Ensure Docker is running:
```bash
docker ps
```

## Usage

1. Start the server:
```bash
go run main.go
```

2. The application will output a public ngrok URL:
```
App public URL is: https://xxxx-xx-xx-xxx-xxx.ngrok.io
```

3. Configure your Git repository webhook:
   - **Payload URL**: `https://your-ngrok-url.ngrok.io/webhook`
   - **Content type**: `application/json`
   - **Events**: Push events

4. When you push to a branch, the system will:
   - Receive the webhook payload
   - Extract the branch name
   - Create a preview container named `preview-app-{branch-name}`
   - Return a preview URL (e.g., `http://localhost:32768`)

## API

### POST /webhook

Receives Git webhook payloads and creates preview environments.

**Request Body:**
```json
{
  "ref": "refs/heads/feature-branch"
}
```

**Response:**
```json
{
  "status": "created",
  "branch": "feature-branch",
  "preview_url": "http://localhost:32768"
}
```

## Project Structure

```
.
├── main.go              # Main application entry point
├── dockerClient/
│   └── dockerClient.go  # Docker container management
├── go.mod               # Go module dependencies
└── README.md            # Project documentation
```

## How It Works

1. **Webhook Reception**: The Gin server listens for POST requests on `/webhook`
2. **Branch Extraction**: Extracts the branch name from the `ref` field
3. **Container Cleanup**: Checks for existing containers with the same name and removes them
4. **Container Creation**: Creates a new nginx:alpine container with auto-assigned ports
5. **URL Generation**: Returns the preview URL with the dynamically allocated port

## Dependencies

- [gin-gonic/gin](https://github.com/gin-gonic/gin) - HTTP web framework
- [docker/docker](https://github.com/docker/docker) - Docker client library
- [ngrok](https://ngrok.com/) - Secure tunneling for webhooks

## License

MIT License
