INSERT INTO roles (name) VALUES ('atletas'), ('treneris'), ('administratorius');

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

INSERT INTO gyms (name, address, city) VALUES
    ('Testinis klubas Vilnius Centras', 'Testinė g. 1', 'Vilnius'),
    ('Testinis klubas Vilnius Šiaurė', 'Testinė g. 2', 'Vilnius'),
    ('Testinis klubas Kaunas Centras', 'Testinė g. 3', 'Kaunas'),
    ('Testinis klubas Kaunas Pietūs', 'Testinė g. 4', 'Kaunas'),
    ('Testinis klubas Klaipėda', 'Testinė g. 5', 'Klaipėda'),
    ('Testinis klubas Šiauliai', 'Testinė g. 6', 'Šiauliai'),
    ('Testinis klubas Panevėžys', 'Testinė g. 7', 'Panevėžys'),
    ('Testinis klubas Alytus', 'Testinė g. 8', 'Alytus');

INSERT INTO users (email, password_hash, role_id, first_name, gender, bio, created_at)
VALUES
    ('athlete1@fitmate.test', '$2a$10$sIIkV4x9IWAuglTU6HobSea.I/urxzc3r7LcamRycXOPcnBy.h0pq', 1, 'Jonas', 'vyras', 'Testinis profilis. Vilniuje ieškau partnerio jėgos treniruotėms.', CURRENT_TIMESTAMP - INTERVAL '14 days'),
    ('athlete2@fitmate.test', '$2a$10$sIIkV4x9IWAuglTU6HobSea.I/urxzc3r7LcamRycXOPcnBy.h0pq', 1, 'Ieva', 'moteris', 'Testinis profilis. Sportuoti pradėjau neseniai, patinka rytinės treniruotės.', CURRENT_TIMESTAMP - INTERVAL '14 days'),
    ('athlete3@fitmate.test', '$2a$10$sIIkV4x9IWAuglTU6HobSea.I/urxzc3r7LcamRycXOPcnBy.h0pq', 1, 'Mantas', 'vyras', 'Testinis profilis. Kaune ruošiuosi jėgos trikovės varžyboms.', CURRENT_TIMESTAMP - INTERVAL '14 days'),
    ('athlete4@fitmate.test', '$2a$10$sIIkV4x9IWAuglTU6HobSea.I/urxzc3r7LcamRycXOPcnBy.h0pq', 1, 'Greta', 'moteris', 'Testinis profilis. Domiuosi funkcinėmis treniruotėmis ir mobilumu.', CURRENT_TIMESTAMP - INTERVAL '14 days'),
    ('athlete5@fitmate.test', '$2a$10$sIIkV4x9IWAuglTU6HobSea.I/urxzc3r7LcamRycXOPcnBy.h0pq', 1, 'Tomas', 'vyras', 'Testinis profilis. Klaipėdoje derinu bėgimą ir sporto salę.', CURRENT_TIMESTAMP - INTERVAL '14 days'),
    ('athlete6@fitmate.test', '$2a$10$sIIkV4x9IWAuglTU6HobSea.I/urxzc3r7LcamRycXOPcnBy.h0pq', 1, 'Alex', 'kita', 'Testinis profilis. Šiauliuose ieškau kompanijos savaitgalio treniruotėms.', CURRENT_TIMESTAMP - INTERVAL '14 days'),
    ('trainer1@fitmate.test', '$2a$10$sIIkV4x9IWAuglTU6HobSea.I/urxzc3r7LcamRycXOPcnBy.h0pq', 2, 'Lukas', 'vyras', 'Testinis trenerio profilis. Vilniuje padedu stiprėti ir auginti raumenų masę.', CURRENT_TIMESTAMP - INTERVAL '14 days'),
    ('trainer2@fitmate.test', '$2a$10$sIIkV4x9IWAuglTU6HobSea.I/urxzc3r7LcamRycXOPcnBy.h0pq', 2, 'Eglė', 'moteris', 'Testinis trenerės profilis. Kaune vedu funkcines ir mobilumo treniruotes.', CURRENT_TIMESTAMP - INTERVAL '14 days'),
    ('trainer3@fitmate.test', '$2a$10$sIIkV4x9IWAuglTU6HobSea.I/urxzc3r7LcamRycXOPcnBy.h0pq', 2, 'Paulius', 'vyras', 'Testinis trenerio profilis. Klaipėdoje dirbu su pradedančiaisiais ir ištverme.', CURRENT_TIMESTAMP - INTERVAL '14 days'),
    ('admin@fitmate.test', '$2a$10$sIIkV4x9IWAuglTU6HobSea.I/urxzc3r7LcamRycXOPcnBy.h0pq', 3, 'Administratorius', 'kita', 'Testinė administratoriaus paskyra sporto klubų ir vartotojų valdymui.', CURRENT_TIMESTAMP - INTERVAL '14 days');

INSERT INTO athlete_profiles (user_id, date_of_birth, height_cm, weight_kg, experience_years)
VALUES
    (1, DATE '1998-03-12', 182.00, 82.50, 3),
    (2, DATE '2002-07-24', 165.00, 58.00, 0),
    (3, DATE '1992-11-08', 190.00, 98.00, 8),
    (4, DATE '1996-05-19', 171.00, 65.50, 4),
    (5, DATE '1987-09-03', 177.00, 76.00, 6),
    (6, DATE '2000-01-16', 174.00, 70.00, 1);

INSERT INTO trainer_profiles (user_id, work_experience_years, education)
VALUES
    (7, 7, 'Sporto bakalauro studijos, asmeninio trenerio kvalifikacija.'),
    (8, 12, 'Kineziterapijos studijos ir funkcinio rengimo kursai.'),
    (9, 3, 'Fizinio aktyvumo studijos ir ištvermės treniravimo kursai.');

INSERT INTO user_gyms (user_id, gym_id)
VALUES
    (1, 1),
    (1, 2),
    (2, 1),
    (2, 2),
    (3, 3),
    (3, 4),
    (4, 3),
    (4, 4),
    (5, 5),
    (6, 6),
    (7, 1),
    (7, 2),
    (8, 3),
    (8, 4),
    (9, 5);

INSERT INTO user_preferred_workout_times (user_id, day_of_week, start_time, end_time)
VALUES
    (1, 1, '18:00', '19:30'),
    (1, 3, '18:00', '19:30'),
    (2, 2, '07:00', '08:00'),
    (2, 4, '07:00', '08:00'),
    (3, 1, '19:00', '21:00'),
    (3, 5, '18:00', '20:00'),
    (4, 3, '17:30', '19:00'),
    (4, 6, '10:00', '11:30'),
    (5, 2, '06:30', '08:00'),
    (5, 7, '09:00', '10:30'),
    (6, 6, '12:00', '13:30'),
    (6, 7, '12:00', '13:30'),
    (7, 1, '16:00', '21:00'),
    (7, 3, '16:00', '21:00'),
    (8, 2, '08:00', '14:00'),
    (8, 4, '14:00', '20:00'),
    (9, 5, '15:00', '20:00'),
    (9, 6, '09:00', '14:00');

INSERT INTO trainer_specializations (trainer_user_id, specialization_id)
VALUES
    (7, 1),
    (7, 2),
    (7, 61),
    (8, 6),
    (8, 25),
    (8, 28),
    (9, 15),
    (9, 16),
    (9, 52);

INSERT INTO swipes (swiper_user_id, target_user_id, action, created_at)
VALUES
    (1, 2, 'patinka', CURRENT_TIMESTAMP - INTERVAL '3 days'),
    (2, 1, 'patinka', CURRENT_TIMESTAMP - INTERVAL '3 days'),
    (3, 4, 'patinka', CURRENT_TIMESTAMP - INTERVAL '3 days'),
    (4, 3, 'patinka', CURRENT_TIMESTAMP - INTERVAL '3 days'),
    (1, 3, 'nepatinka', CURRENT_TIMESTAMP - INTERVAL '3 days'),
    (5, 6, 'patinka', CURRENT_TIMESTAMP - INTERVAL '3 days'),
    (6, 5, 'nepatinka', CURRENT_TIMESTAMP - INTERVAL '3 days');

INSERT INTO chats (user1_id, user2_id, type, created_at)
VALUES
    (1, 2, 'match', CURRENT_TIMESTAMP - INTERVAL '2 days'),
    (3, 4, 'match', CURRENT_TIMESTAMP - INTERVAL '2 days'),
    (1, 7, 'trainer', CURRENT_TIMESTAMP - INTERVAL '2 days'),
    (5, 9, 'trainer', CURRENT_TIMESTAMP - INTERVAL '2 days');

INSERT INTO messages (chat_id, sender_id, content, created_at)
VALUES
    (1, 1, 'Labas! Gal pirmadienį pasportuojame kartu?', CURRENT_TIMESTAMP - INTERVAL '90 minutes'),
    (1, 2, 'Labas! Šį kartą galiu 18 valandą.', CURRENT_TIMESTAMP - INTERVAL '85 minutes'),
    (2, 3, 'Sveika! Kokią treniruotę planuoji savaitgalį?', CURRENT_TIMESTAMP - INTERVAL '80 minutes'),
    (2, 4, 'Planuoju funkcinę treniruotę šeštadienį ryte.', CURRENT_TIMESTAMP - INTERVAL '75 minutes'),
    (3, 1, 'Sveiki, norėčiau aptarti jėgos treniruočių planą.', CURRENT_TIMESTAMP - INTERVAL '70 minutes'),
    (3, 7, 'Sveiki! Parašykite, kokie jūsų tikslai ir patirtis.', CURRENT_TIMESTAMP - INTERVAL '65 minutes'),
    (4, 5, 'Ar padedate pasiruošti pusmaratoniui?', CURRENT_TIMESTAMP - INTERVAL '60 minutes'),
    (4, 9, 'Taip, galime derinti bėgimą su jėgos treniruotėmis.', CURRENT_TIMESTAMP - INTERVAL '55 minutes');

INSERT INTO reviews (trainer_user_id, author_id, rating, comment, created_at)
VALUES
    (7, 1, 5, 'Aiškiai paaiškino pratimų techniką ir sudarė tinkamą planą.', CURRENT_TIMESTAMP - INTERVAL '1 day'),
    (7, 2, 4, 'Kantrus treneris, padėjo drąsiau pradėti sportuoti.', CURRENT_TIMESTAMP - INTERVAL '1 day'),
    (8, 4, 5, 'Patiko individualiai pritaikyti mobilumo pratimai.', CURRENT_TIMESTAMP - INTERVAL '1 day'),
    (9, 5, 3, 'Naudingi patarimai, tačiau norėčiau daugiau laisvų treniruočių laikų.', CURRENT_TIMESTAMP - INTERVAL '1 day');
