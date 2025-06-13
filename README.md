# social-network

## Overview

This project is a Facebook-like social network with various features such as followers, profiles, posts, groups, notifications, and private chats. The project is built using a combination of frontend and backend technologies, and it is containerized using Docker.

## Technologies Used

### Frontend

The frontend is developed using [React.js](https://reactjs.org/), a popular JavaScript library for building user interfaces. React.js helps create a dynamic and interactive user experience.

### Backend

The backend is implemented in [Go](https://golang.org/), utilizing the [Gorilla](https://pkg.go.dev/github.com/gorilla/websocket) websocket package for real-time communication. The database management is handled by [SQLite](https://www.sqlite.org/) with migrations managed using [golang-migrate](https://github.com/golang-migrate/migrate/).

### Containerization

[Docker](https://www.docker.com/) is used to containerize the application. There are two Docker images: one for the backend and another for the frontend. Docker makes it easy to deploy and manage the application across different environments.

## Project Structure

The project is organized into frontend and backend folders, each with its own set of functionalities.

- `frontend/`: Contains the React.js application for the user interface.
- `backend/`: Contains the Go application for the server-side logic, handling requests, and interacting with the database.

## Running the Application

### Docker Compose

To run the entire application using Docker Compose, follow these steps:

1. Make sure you have Docker and Docker Compose installed on your machine.

2. Clone the repository:

   ```bash
   git clone https://01.kood.tech/git/Aleksander/social-network.git
   cd social-network
   ```

3. Run Docker Compose:

   ```bash
   docker-compose up -d
   ```

4. Access the application in your browser at `http://localhost:3000`.

### Running Separately

If you prefer to run the frontend and backend separately:

#### Backend

1. Navigate to the `backend/` directory:

   ```bash
   cd backend
   ```

2. Run the backend server:

   ```bash
   go run .
   ```

#### Frontend

1. Navigate to the `frontend/` directory:

   ```bash
   cd frontend
   ```

2. Install dependencies:

   ```bash
   npm install
   ```

3. Run the frontend:

   ```bash
   npm start
   ```

4. Access the application in your browser at `http://localhost:3000`.

## Authors

- [Aleksander](https://01.kood.tech/git/Aleksander)
- [Jklimenk](https://01.kood.tech/git/jklimenk)
- [Jegor_petsorin](https://01.kood.tech/git/jegor_petsorin)
