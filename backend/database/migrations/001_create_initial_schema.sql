CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(16) NOT NULL UNIQUE
);

INSERT INTO roles (name) VALUES ('atletas'), ('treneris'), ('administratorius');

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
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

CREATE table gyms (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    address VARCHAR(255) NOT NULL,
    city VARCHAR(100) NOT NULL
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
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    date_of_birth DATE NOT NULL,
    height_cm DECIMAL(5, 2) NOT NULL CHECK (height_cm > 0),
    weight_kg DECIMAL(5, 2) NOT NULL CHECK (weight_kg > 0),
    experience_years INTEGER NOT NULL CHECK (experience_years >= 0)
);

CREATE TABLE trainer_profiles (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    certification VARCHAR(255) NOT NULL,
    work_experience_years INTEGER NOT NULL CHECK (work_experience_years >= 0),
    education VARCHAR(255) NOT NULL
);

CREATE TABLE specializations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);

INSERT INTO specializations (name) VALUES
('Jėgos treniruotės'),
('Raumenų masės auginimas'),
('Kultūrizmas'),
('Jėgos trikovė'),
('Olimpinė sunkioji atletika'),
('Funkcinės treniruotės'),
('CrossFit treniruotės'),
('Kalistenika'),
('Treniruotės su kūno svoriu'),
('Bendras fizinis pasirengimas'),

('Svorio metimas'),
('Riebalų mažinimas'),
('Svorio kontrolė'),
('Kūno kompozicijos gerinimas'),

('Kardio treniruotės'),
('Ištvermės treniruotės'),
('Bėgimas'),
('Maratono pasiruošimas'),
('Dviračių sportas'),
('Plaukimas'),
('Triatlono treniruotės'),

('HIIT treniruotės'),
('Žiedinės treniruotės'),
('Intensyvios grupinės treniruotės'),

('Mobilumo treniruotės'),
('Lankstumo treniruotės'),
('Tempimo pratimai'),
('Laikysenos gerinimas'),
('Pusiausvyros treniruotės'),
('Koordinacijos lavinimas'),

('Kūno centro treniruotės'),
('Nugaros stiprinimas'),
('Apatinės kūno dalies treniruotės'),
('Viršutinės kūno dalies treniruotės'),
('Viso kūno treniruotės'),

('Sportinis pasirengimas'),
('Atletinis rengimas'),
('Greičio lavinimas'),
('Vikrumo lavinimas'),
('Sprogstamosios jėgos lavinimas'),

('Futbolo fizinis rengimas'),
('Krepšinio fizinis rengimas'),
('Teniso fizinis rengimas'),
('Bokso treniruotės'),
('Kovos menų fizinis rengimas'),

('Reabilitacinės treniruotės'),
('Traumų prevencija'),
('Treniruotės po traumų'),
('Korekcinės treniruotės'),

('Senjorų fizinis aktyvumas'),
('Jaunimo fizinis rengimas'),
('Pradedančiųjų treniruotės'),
('Moterų fitnesas'),
('Vyrų fitnesas'),

('Treniruotės nėščiosioms'),
('Treniruotės po gimdymo'),

('Joga'),
('Pilatesas'),

('Mitybos konsultacijos'),
('Sveikos gyvensenos konsultacijos'),

('Asmeninės treniruotės'),
('Grupinės treniruotės'),
('Nuotolinės treniruotės');

CREATE TABLE trainer_specializations (
    id SERIAL PRIMARY KEY,
    trainer_id INTEGER NOT NULL REFERENCES trainer_profiles(id) ON DELETE CASCADE,
    specialization_id INTEGER NOT NULL REFERENCES specializations(id) ON DELETE CASCADE,
    UNIQUE (trainer_id, specialization_id)
);

CREATE TABLE swipes (
    id SERIAL PRIMARY KEY,
    swiper_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    target_user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    action VARCHAR(9) NOT NULL CHECK (action IN ('patinka', 'nepatinka')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CHECK (swiper_id <> target_user_id),
    UNIQUE (swiper_id, target_user_id)
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