-- Пароль для всех: "password123"
INSERT INTO users (username, password_hash, balance, referrer_id) VALUES
    ('alice', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 100, NULL),
    ('bob', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 250, 1),
    ('charlie', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 50, 1),
    ('david', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 500, 2),
    ('eve', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 150, 2),
    ('frank', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 300, NULL),
    ('grace', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 75, 6),
    ('heidi', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 400, NULL),
    ('ivan', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 200, 8),
    ('judy', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 125, 8)
ON CONFLICT (username) DO NOTHING;

INSERT INTO tasks (title, description, points) VALUES
    ('Подписаться на Telegram канал', 'Подпишитесь на наш официальный Telegram канал и получите бонус', 50),
    ('Подписаться на Twitter', 'Подпишитесь на наш Twitter аккаунт', 50),
    ('Пригласить друга', 'Используйте реферальный код и пригласите друга', 100),
    ('Пройти верификацию', 'Пройдите верификацию личности', 150),
    ('Оставить отзыв', 'Оставьте отзыв о нашем сервисе', 25),
    ('Поделиться в соцсетях', 'Поделитесь ссылкой на наш проект в социальных сетях', 75),
    ('Пройти опрос', 'Примите участие в нашем опросе', 30),
    ('Написать статью', 'Напишите статью о нашем проекте', 200),
    ('Создать мем', 'Создайте мем о проекте', 40),
    ('Участвовать в конкурсе', 'Примите участие в ежемесячном конкурсе', 150)
ON CONFLICT DO NOTHING;


INSERT INTO user_tasks (user_id, task_id) VALUES
    (1, 1),  -- alice выполнила задание 1
    (1, 2),  -- alice выполнила задание 2
    (2, 1),  -- bob выполнил задание 1
    (2, 3),  -- bob выполнил задание 3
    (2, 5),  -- bob выполнил задание 5
    (3, 1),  -- charlie выполнил задание 1
    (4, 1),  -- david выполнил задание 1
    (4, 2),  -- david выполнил задание 2
    (4, 3),  -- david выполнил задание 3
    (4, 4),  -- david выполнил задание 4
    (5, 2),  -- eve выполнила задание 2
    (5, 5),  -- eve выполнила задание 5
    (6, 1),  -- frank выполнил задание 1
    (6, 3),  -- frank выполнил задание 3
    (6, 6),  -- frank выполнил задание 6
    (7, 4),  -- grace выполнила задание 4
    (8, 1),  -- heidi выполнила задание 1
    (8, 2),  -- heidi выполнила задание 2
    (8, 3),  -- heidi выполнила задание 3
    (9, 5),  -- ivan выполнил задание 5
    (10, 1)  -- judy выполнила задание 1
ON CONFLICT DO NOTHING;
