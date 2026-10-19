# FitMate

FitMate is a social networking application designed to help people who train at the gym find workout partners or personal trainers.

## System Purpose

The goal of the project is to create a social network that helps people who work out at the gym find a workout partner or a personal trainer.

Users will be able to create an account and a personal profile. In their profile, they will be able to specify their personal fitness goals, experience, gym, preferred workout time, and other useful information.

Users will be shown profiles of other users based on selected filters. While viewing a profile, a user will be able to express interest in another person or skip the profile. If two users express interest in each other, they will be able to communicate with each other through messages.

The platform will also allow users to search for personal trainers. Trainers will be able to create a special trainer profile where they can provide information about their qualifications, specialization, experience, and other relevant details. Users will be able to choose a suitable trainer and send them a message.

The system will consist of a client application and an API. It will be used by athletes, personal trainers, and administrators.

## Functional Requirements

### Unregistered User

An unregistered user will be able to:

1. View the platform's landing page.
2. Register in the system.

### Registered User (Athlete)

A registered user will be able to:

1. Log in to the system.
2. Log out of the system.
3. View their profile.
4. Edit their profile information.
5. View other users' profiles based on selected filters.
6. Express interest in another user.
7. Skip an unwanted user profile.
8. View mutual matches and their profiles.
9. Exchange messages with users they have matched with.
10. View a list of trainers based on selected filters.
11. View a selected trainer's profile.

### Trainer

A trainer will be able to:

1. Log in to the system.
2. Log out of the system.
3. View their profile.
4. Edit their profile.
5. View users who are interested in the trainer and their profiles.
6. Exchange messages with interested users.

### Administrator

An administrator will be able to:

1. Log in to the system.
2. Log out of the system.
3. View the list of athletes and trainers.
4. Manage athletes and trainers by editing their information or removing them.

## Technologies

- **Go 1.26** with `net/http` for the backend API.
- **PostgreSQL 17** with `pgx` for database access.
- **JWT** for authentication and role-based access control.
- **bcrypt** for password hashing.
- **Docker and Docker Compose** for running the API and database.
- **OpenAPI 3.0** for API documentation.
- **Postman** for API testing.

## Running the Project

Install **Git** and **Docker Desktop**. Start Docker Desktop before running the commands below. These examples use PowerShell on Windows.

1. **Download the project:**

   ```powershell
   git clone https://github.com/Tautris404/fitmate.git
   cd fitmate
   ```

   If you already have the project, open a terminal in the repository root, which contains `docker-compose.yml`.

2. **Create the environment file:**

   ```powershell
   Copy-Item .env.example .env
   ```

   Open `.env` and configure the database credentials, ports and token settings. `JWT_SECRET` must contain at least 32 bytes.

3. **Build and start the API and database:**

   ```powershell
   docker compose up --build -d
   ```

   Run this command from the repository root. Docker builds the Go backend and starts both containers in the background. The first build may take a few minutes.

   On the first database startup, the SQL files in `backend/database/migrations` create the schema and insert the initial data. Initialization runs only when the database volume is empty.

4. **Access the API:**

   Open `http://localhost:<API_PORT>/gyms` in a browser or Postman, replacing `<API_PORT>` with the value configured in `.env`.

After changing backend code, run `docker compose up --build -d` again to rebuild it.

To stop the project:

```powershell
docker compose down
```

The database and uploaded images remain stored in Docker volumes.

## API Documentation

The [OpenAPI specification](docs/fitmate-api.yaml) describes the implemented endpoints, request bodies, responses and access restrictions.

## Running Postman Tests

Start the project using the instructions above, then open the Postman desktop application.

1. Import [fitmate-postman-tests.json](backend/tests/fitmate-postman-tests.json).
2. In Postman settings, set the [working directory](https://learning.postman.com/docs/getting-started/installation/settings/working-directory) to the repository root, which contains `docker-compose.yml`.
3. Open the `fitmate tests` collection and select **Run**. Use **1 iteration**, **0 ms delay** and keep all requests selected in their existing order.
4. Start the run and view the results.
