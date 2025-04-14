CREATE TABLE IF NOT EXISTS tasks
(
    id SERIAL not null
        constraint pk_tasks
            primary key,
    title TEXT NOT NULL,
    description TEXT,
    status TEXT CHECK (status IN ('new', 'in_progress', 'done')) DEFAULT 'new',
    created_at timestamptz,
    updated_at timestamptz
);