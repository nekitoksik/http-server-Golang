CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    points INTEGER NOT NULL CHECK (points > 0)
);

CREATE INDEX idx_tasks_points ON tasks(points);
