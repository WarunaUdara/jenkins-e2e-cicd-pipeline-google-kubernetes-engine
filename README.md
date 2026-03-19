# Jenkins End-to-End CI/CD Pipeline to Kubernetes

End-to-end Jenkins CI/CD pipeline deploying a Go REST API Todo application to Kubernetes using ArgoCD and GitOps principles.

![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go)
![Docker](https://img.shields.io/badge/Docker-Multi--stage-2496ED?logo=docker)
![Jenkins](https://img.shields.io/badge/Jenkins-CI%2FCD-D24939?logo=jenkins)
![Kubernetes](https://img.shields.io/badge/Kubernetes-Deployment-326CE5?logo=kubernetes)
![ArgoCD](https://img.shields.io/badge/ArgoCD-GitOps-EF7B4D?logo=argo)

## Architecture Overview

```
Developer → GitHub (App Repo) → [Webhook] → Jenkins Pipeline
                                              ↓
                                    Build Docker Image
                                              ↓
                                    Push to DockerHub
                                              ↓
                                GitHub (Manifest Repo)
                                              ↓
                                    ArgoCD (Monitoring)
                                              ↓
                                         Kubernetes
                                              ↓
                                         End Users
```

## Features

- **Go REST API** with full CRUD operations
- **Multi-stage Docker build** for optimized image size (12MB)
- **Jenkins Pipeline** with 5 automated stages
- **GitOps** workflow with ArgoCD
- **Kubernetes** deployment with health checks
- **Zero-downtime deployments** with rolling updates

## Tech Stack

| Component | Technology |
|-----------|-----------|
| Application | Go + Gin Framework |
| Containerization | Docker (Multi-stage build) |
| CI/CD | Jenkins |
| Container Registry | DockerHub |
| Orchestration | Kubernetes (GKE) |
| GitOps | ArgoCD |
| Source Control | GitHub |

## Project Structure

```
.
├── main.go                    # Application entry point
├── go.mod                     # Go module dependencies
├── handlers/
│   ├── todo.go               # CRUD handlers
│   └── health.go             # Health check handlers
├── models/
│   └── todo.go               # Todo data model
├── repository/
│   └── todo_repo.go          # Thread-safe in-memory storage
├── static/
│   └── css/
│       └── style.css         # Frontend styling
├── templates/
│   └── index.html            # API documentation page
├── Dockerfile                # Multi-stage Docker build
├── .dockerignore             # Docker build optimization
├── Jenkinsfile               # CI/CD pipeline definition
└── README.md                 # This file
```

## API Endpoints

### Todo Operations

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/todos` | List all todos |
| GET | `/api/todos/:id` | Get a specific todo |
| POST | `/api/todos` | Create a new todo |
| PUT | `/api/todos/:id` | Update a todo |
| DELETE | `/api/todos/:id` | Delete a todo |

### Health Checks

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Application health status |
| GET | `/ready` | Readiness probe for Kubernetes |

## Quick Start

### Prerequisites

- Go 1.26+
- Docker
- Git

### Local Development

```bash
# Clone repository
git clone https://github.com/WarunaUdara/jenkins-e2e-cicd-pipeline-google-kubernetes-engine.git
cd jenkins-e2e-cicd-pipeline-google-kubernetes-engine

# Install dependencies
go mod download

# Run application
go run main.go

# Access the application
open http://localhost:8080
```

### Docker

```bash
# Build image
docker build -t warunaudara/todo-app-go:latest .

# Run container
docker run -d -p 8080:8080 warunaudara/todo-app-go:latest

# Test endpoints
curl http://localhost:8080/health
curl http://localhost:8080/api/todos
```

## Jenkins Pipeline

### Pipeline Stages

1. **Checkout Application Code** - Clone app repository
2. **Build Docker Image** - Create optimized Docker image
3. **Push to DockerHub** - Upload image to registry
4. **Checkout K8s Manifests** - Clone manifest repository
5. **Update & Push Manifests** - Update image tag and trigger ArgoCD sync

### Required Jenkins Credentials

| Credential ID | Type | Description |
|---------------|------|-------------|
| `github-credentials` | Username/Password | GitHub PAT with `repo` scope |
| `dockerhub-credentials` | Username/Password | DockerHub username and access token |

### Environment Variables

```groovy
DOCKERHUB_USERNAME = 'warunaudara'
DOCKERHUB_REPO = 'warunaudara/todo-app-go'
IMAGE_TAG = BUILD_NUMBER
APP_REPO = 'https://github.com/WarunaUdara/jenkins-e2e-cicd-pipeline-google-kubernetes-engine.git'
MANIFEST_REPO = 'https://github.com/WarunaUdara/jenkins-e2e-cicd-pipeline-manifests.git'
```

## Docker Image

**Image:** `warunaudara/todo-app-go:latest`  
**Size:** 12MB (compressed) / 42.8MB (uncompressed)  
**Base:** Alpine Linux  

### Optimizations

- Multi-stage build (builder + runtime)
- Static binary compilation (CGO_ENABLED=0)
- Debug symbols stripped (-ldflags="-w -s")
- Non-root user for security
- Minimal Alpine base image

## Kubernetes Deployment

### Resources

- **Deployment**: 3 replicas with rolling update strategy
- **Service**: LoadBalancer type for external access
- **Health Checks**: Liveness and readiness probes
- **Resource Limits**: Memory and CPU constraints

### ArgoCD GitOps

- **Auto-sync**: Enabled
- **Self-heal**: Automatic recovery
- **Prune**: Remove orphaned resources
- **Manifest Repo**: [jenkins-e2e-cicd-pipeline-manifests](https://github.com/WarunaUdara/jenkins-e2e-cicd-pipeline-manifests)

## API Usage Examples

### Create a Todo

```bash
curl -X POST http://localhost:8080/api/todos \
  -H "Content-Type: application/json" \
  -d '{"title": "Learn Jenkins CI/CD"}'
```

### List All Todos

```bash
curl http://localhost:8080/api/todos
```

### Update a Todo

```bash
curl -X PUT http://localhost:8080/api/todos/1 \
  -H "Content-Type: application/json" \
  -d '{"is_completed": true}'
```

### Delete a Todo

```bash
curl -X DELETE http://localhost:8080/api/todos/1
```

## CI/CD Workflow

1. **Developer** pushes code to GitHub
2. **GitHub Webhook** triggers Jenkins build
3. **Jenkins** executes 5-stage pipeline:
   - Builds Docker image
   - Pushes to DockerHub
   - Updates Kubernetes manifests
4. **ArgoCD** detects manifest changes
5. **Kubernetes** performs rolling update
6. **Application** deployed with zero downtime

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Author

**Waruna Udara**

- GitHub: [@WarunaUdara](https://github.com/WarunaUdara)
- DockerHub: [warunaudara](https://hub.docker.com/u/warunaudara)

## Acknowledgments

- Jenkins for CI/CD automation
- ArgoCD for GitOps workflow
- Kubernetes for orchestration
- Docker for containerization
- Go community for excellent frameworks

---

**Built with ❤️ using Go, Docker, Jenkins, Kubernetes, and ArgoCD**
