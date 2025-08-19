#!/bin/bash

# Development startup script for the golang-training project
# This script starts all services in the correct order

echo "🚀 Starting Golang Training Development Environment"
echo "=================================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to check if port is in use
check_port() {
    local port=$1
    if lsof -Pi :$port -sTCP:LISTEN -t >/dev/null ; then
        return 0  # Port is in use
    else
        return 1  # Port is free
    fi
}

# Function to kill process on port
kill_port() {
    local port=$1
    echo -e "${YELLOW}Killing any existing process on port $port...${NC}"
    sudo fuser -k $port/tcp 2>/dev/null || true
    sleep 2
}

# Clean up existing processes
echo -e "${YELLOW}Cleaning up existing processes...${NC}"
kill_port 8080  # User service
kill_port 7070  # Asset service
kill_port 5173  # Frontend

echo ""
echo -e "${GREEN}Starting services...${NC}"

# Start User Service
echo -e "${BLUE}1. Starting User Service (port 8080)...${NC}"
cd user-service
if [ ! -f "cmd/api/main.go" ]; then
    echo -e "${RED}❌ User service main.go not found${NC}"
    exit 1
fi

gnome-terminal --title="User Service" -- bash -c "
    echo 'Starting User Service on port 8080...';
    go run cmd/api/main.go;
    echo 'User Service stopped. Press any key to close...';
    read -n 1;
" &

sleep 3

# Check if user service started
if check_port 8080; then
    echo -e "${GREEN}✅ User Service started on port 8080${NC}"
else
    echo -e "${RED}❌ User Service failed to start${NC}"
fi

# Start Asset Service
echo -e "${BLUE}2. Starting Asset Service (port 7070)...${NC}"
cd ../asset-service
if [ ! -f "cmd/api/main.go" ]; then
    echo -e "${RED}❌ Asset service main.go not found${NC}"
    exit 1
fi

gnome-terminal --title="Asset Service" -- bash -c "
    echo 'Starting Asset Service on port 7070...';
    go run cmd/api/main.go;
    echo 'Asset Service stopped. Press any key to close...';
    read -n 1;
" &

sleep 3

# Check if asset service started
if check_port 7070; then
    echo -e "${GREEN}✅ Asset Service started on port 7070${NC}"
else
    echo -e "${RED}❌ Asset Service failed to start${NC}"
fi

# Start Frontend
echo -e "${BLUE}3. Starting Frontend (port 5173)...${NC}"
cd ../team-fe
if [ ! -f "package.json" ]; then
    echo -e "${RED}❌ Frontend package.json not found${NC}"
    exit 1
fi

gnome-terminal --title="Frontend" -- bash -c "
    echo 'Starting Frontend on port 5173...';
    npm run dev;
    echo 'Frontend stopped. Press any key to close...';
    read -n 1;
" &

sleep 5

# Check if frontend started
if check_port 5173; then
    echo -e "${GREEN}✅ Frontend started on port 5173${NC}"
else
    echo -e "${RED}❌ Frontend failed to start${NC}"
fi

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}🎉 All services started successfully!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo -e "${BLUE}Services:${NC}"
echo -e "📊 User Service (GraphQL + Teams):  ${BLUE}http://localhost:8080${NC}"
echo -e "📁 Asset Service (Folders/Notes):   ${BLUE}http://localhost:7070${NC}"
echo -e "🌐 Frontend Application:           ${BLUE}http://localhost:5173${NC}"
echo ""
echo -e "${BLUE}API Endpoints:${NC}"
echo -e "🔍 GraphQL Playground:             ${BLUE}http://localhost:8080/user/query${NC}"
echo -e "🏥 User Service Health:            ${BLUE}http://localhost:8080/health${NC}"
echo -e "🏥 Asset Service Health:           ${BLUE}http://localhost:7070/health${NC}"
echo ""
echo -e "${YELLOW}📋 Quick Start Guide:${NC}"
echo -e "1. Open ${BLUE}http://localhost:5173${NC} in your browser"
echo -e "2. Create a user account (manager or member)"
echo -e "3. Login with your credentials"
echo -e "4. Start creating teams and managing assets!"
echo ""
echo -e "${YELLOW}To stop all services:${NC}"
echo -e "Press Ctrl+C in each terminal window or run:"
echo -e "sudo fuser -k 8080/tcp && sudo fuser -k 7070/tcp && sudo fuser -k 5173/tcp"

# Return to original directory
cd ..

echo ""
echo -e "${GREEN}Development environment ready! 🚀${NC}"
