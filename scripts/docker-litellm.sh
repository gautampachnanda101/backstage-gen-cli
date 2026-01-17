#!/bin/bash
# Docker helper script for LiteLLM
# Usage: ./docker-litellm.sh [start|stop|status|logs|restart]

set -e

CONTAINER_NAME="litellm"
IMAGE="ghcr.io/berriai/litellm:main-latest"
PORT="4000"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_docker() {
    if ! command -v docker &> /dev/null; then
        print_error "Docker is not installed. Please install Docker first."
        echo "Visit: https://docs.docker.com/get-docker/"
        exit 1
    fi

    if ! docker info &> /dev/null; then
        print_error "Docker daemon is not running. Please start Docker."
        exit 1
    fi
}

start_litellm() {
    check_docker

    if docker ps -q -f name="$CONTAINER_NAME" | grep -q .; then
        print_warning "LiteLLM container is already running"
        return 0
    fi

    # Check if container exists but is stopped
    if docker ps -aq -f name="$CONTAINER_NAME" | grep -q .; then
        print_status "Starting existing LiteLLM container..."
        docker start "$CONTAINER_NAME"
    else
        print_status "Creating and starting LiteLLM container..."
        docker run -d \
            --name "$CONTAINER_NAME" \
            -p "$PORT:$PORT" \
            --restart unless-stopped \
            "$IMAGE"
    fi

    print_status "Waiting for LiteLLM to be ready..."
    for i in {1..30}; do
        if curl -s "http://localhost:$PORT/health" &> /dev/null; then
            print_status "LiteLLM is ready on http://localhost:$PORT"
            return 0
        fi
        sleep 1
        echo -n "."
    done
    echo ""
    print_warning "LiteLLM may still be starting. Check logs with: $0 logs"
}

stop_litellm() {
    check_docker

    if docker ps -q -f name="$CONTAINER_NAME" | grep -q .; then
        print_status "Stopping LiteLLM container..."
        docker stop "$CONTAINER_NAME"
        print_status "LiteLLM stopped"
    else
        print_warning "LiteLLM container is not running"
    fi
}

status_litellm() {
    check_docker

    if docker ps -q -f name="$CONTAINER_NAME" | grep -q .; then
        print_status "LiteLLM is running"
        echo ""
        echo "Container details:"
        docker ps -f name="$CONTAINER_NAME" --format "  ID: {{.ID}}\n  Image: {{.Image}}\n  Status: {{.Status}}\n  Ports: {{.Ports}}"
        echo ""
        if curl -s "http://localhost:$PORT/health" &> /dev/null; then
            print_status "Health check: OK"
        else
            print_warning "Health check: Not responding"
        fi
    elif docker ps -aq -f name="$CONTAINER_NAME" | grep -q .; then
        print_warning "LiteLLM container exists but is stopped"
        echo "Run '$0 start' to start it"
    else
        print_warning "LiteLLM container does not exist"
        echo "Run '$0 start' to create and start it"
    fi
}

logs_litellm() {
    check_docker

    if docker ps -aq -f name="$CONTAINER_NAME" | grep -q .; then
        docker logs -f "$CONTAINER_NAME"
    else
        print_error "LiteLLM container does not exist"
        exit 1
    fi
}

restart_litellm() {
    stop_litellm
    start_litellm
}

remove_litellm() {
    check_docker

    if docker ps -q -f name="$CONTAINER_NAME" | grep -q .; then
        print_status "Stopping LiteLLM container..."
        docker stop "$CONTAINER_NAME"
    fi

    if docker ps -aq -f name="$CONTAINER_NAME" | grep -q .; then
        print_status "Removing LiteLLM container..."
        docker rm "$CONTAINER_NAME"
        print_status "LiteLLM container removed"
    else
        print_warning "LiteLLM container does not exist"
    fi
}

usage() {
    echo "Docker helper script for LiteLLM"
    echo ""
    echo "Usage: $0 [command]"
    echo ""
    echo "Commands:"
    echo "  start     Start LiteLLM container"
    echo "  stop      Stop LiteLLM container"
    echo "  status    Show LiteLLM status"
    echo "  logs      Show LiteLLM logs (follow)"
    echo "  restart   Restart LiteLLM container"
    echo "  remove    Remove LiteLLM container"
    echo ""
    echo "Environment variables:"
    echo "  PORT      LiteLLM port (default: 4000)"
    echo ""
    echo "Examples:"
    echo "  $0 start"
    echo "  $0 logs"
    echo "  PORT=8000 $0 start"
}

# Override port from environment
if [ -n "$LITELLM_PORT" ]; then
    PORT="$LITELLM_PORT"
fi

case "${1:-}" in
    start)
        start_litellm
        ;;
    stop)
        stop_litellm
        ;;
    status)
        status_litellm
        ;;
    logs)
        logs_litellm
        ;;
    restart)
        restart_litellm
        ;;
    remove)
        remove_litellm
        ;;
    *)
        usage
        exit 1
        ;;
esac
