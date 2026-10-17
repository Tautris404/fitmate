CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(16) NOT NULL UNIQUE
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role_id INTEGER NOT NULL REFERENCES roles(id),
    first_name VARCHAR(50) NOT NULL,
    gender VARCHAR(7) NOT NULL CHECK (gender IN ('vyras', 'moteris', 'kita')),
    bio TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE user_images (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    image_url VARCHAR(255) NOT NULL
);

CREATE TABLE gyms (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    address VARCHAR(255) NOT NULL,
    city VARCHAR(100) NOT NULL,
    UNIQUE (name, address, city)
);

CREATE TABLE user_gyms (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    gym_id INTEGER NOT NULL REFERENCES gyms(id) ON DELETE CASCADE,
    UNIQUE (user_id, gym_id)
);

CREATE EXTENSION IF NOT EXISTS btree_gist;

CREATE TABLE user_preferred_workout_times (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    day_of_week SMALLINT NOT NULL CHECK (day_of_week BETWEEN 1 AND 7),
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,

    CHECK (start_time BETWEEN TIME '00:00' AND TIME '23:59'),
    CHECK (end_time BETWEEN TIME '00:00' AND TIME '23:59'),
    CHECK (EXTRACT(SECOND FROM start_time) = 0),
    CHECK (EXTRACT(SECOND FROM end_time) = 0),
    CHECK (end_time > start_time),

    EXCLUDE USING gist (
        user_id WITH =,
        day_of_week WITH =,
        tsrange(
            DATE '2000-01-01' + start_time,
            DATE '2000-01-01' + end_time,
            '[)'
        ) WITH &&
    )
);

CREATE TABLE athlete_profiles (
    user_id INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    date_of_birth DATE NOT NULL,
    height_cm DECIMAL(5, 2) NOT NULL CHECK (height_cm > 0),
    weight_kg DECIMAL(5, 2) NOT NULL CHECK (weight_kg > 0),
    experience_years INTEGER NOT NULL CHECK (experience_years >= 0)
);

CREATE TABLE trainer_profiles (
    user_id INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    work_experience_years INTEGER NOT NULL CHECK (work_experience_years >= 0),
    education VARCHAR(255) NOT NULL
);

CREATE TABLE specializations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE trainer_specializations (
    id SERIAL PRIMARY KEY,
    trainer_user_id INTEGER NOT NULL REFERENCES trainer_profiles(user_id) ON DELETE CASCADE,
    specialization_id INTEGER NOT NULL REFERENCES specializations(id) ON DELETE CASCADE,
    UNIQUE (trainer_user_id, specialization_id)
);

CREATE TABLE swipes (
    id SERIAL PRIMARY KEY,
    swiper_user_id INTEGER NOT NULL REFERENCES athlete_profiles(user_id) ON DELETE CASCADE,
    target_user_id INTEGER NOT NULL REFERENCES athlete_profiles(user_id) ON DELETE CASCADE,
    action VARCHAR(9) NOT NULL CHECK (action IN ('patinka', 'nepatinka')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CHECK (swiper_user_id <> target_user_id),
    UNIQUE (swiper_user_id, target_user_id)
);

CREATE TABLE chats (
    id SERIAL PRIMARY KEY,
    user1_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user2_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(10) NOT NULL CHECK (type IN ('match', 'trainer')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CHECK (user1_id <> user2_id)
);

CREATE UNIQUE INDEX unique_chat_users
ON chats (
    LEAST(user1_id, user2_id),
    GREATEST(user1_id, user2_id)
);

CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    chat_id INTEGER NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    sender_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL CHECK (LENGTH(TRIM(content)) > 0),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE auth_sessions (
    id BIGSERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ
);

CREATE INDEX auth_sessions_user_id_idx
    ON auth_sessions(user_id);

CREATE TABLE refresh_tokens (
    id BIGSERIAL PRIMARY KEY,
    session_id BIGINT NOT NULL REFERENCES auth_sessions(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ
);

CREATE INDEX refresh_tokens_session_id_idx
    ON refresh_tokens(session_id);

CREATE UNIQUE INDEX refresh_tokens_one_unused_per_session_idx
    ON refresh_tokens(session_id)
    WHERE used_at IS NULL;

CREATE TABLE reviews (
    id SERIAL PRIMARY KEY,
    trainer_user_id INTEGER NOT NULL REFERENCES trainer_profiles(user_id) ON DELETE CASCADE,
    author_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    rating SMALLINT NOT NULL CHECK (rating BETWEEN 1 AND 5),
    comment TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (trainer_user_id, author_id)
);
