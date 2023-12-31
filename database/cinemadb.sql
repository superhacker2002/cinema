CREATE TABLE roles (
    role_id SERIAL PRIMARY KEY,
    role_name VARCHAR(50) NOT NULL
);

CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    hashed_password VARCHAR(64) NOT NULL,
    email VARCHAR(50),
    credit_card_info VARCHAR(50),
    role_id INTEGER NOT NULL,
    CONSTRAINT users_role_id_fkey FOREIGN KEY (role_id)
        REFERENCES roles (role_id) ON DELETE RESTRICT
);

CREATE TABLE movies (
    movie_id SERIAL PRIMARY KEY,
    title VARCHAR(50) NOT NULL,
    genre VARCHAR(50) NOT NULL,
    release_date DATE NOT NULL,
    duration INTEGER NOT NULL
);

CREATE TABLE halls (
    hall_id SERIAL PRIMARY KEY,
    hall_name VARCHAR(50) NOT NULL,
    capacity INTEGER NOT NULL
);

CREATE TABLE cinema_sessions (
    session_id SERIAL PRIMARY KEY,
    movie_id INTEGER NOT NULL,
    hall_id INTEGER NOT NULL,
    start_time timestamptz NOT NULL,
    end_time timestamptz NOT NULL,
    price DECIMAL(5,2) NOT NULL,
    CONSTRAINT cinema_sessions_movie_id_fkey FOREIGN KEY (movie_id)
        REFERENCES movies (movie_id) ON DELETE CASCADE,
    CONSTRAINT cinema_sessions_hall_id_fkey FOREIGN KEY (hall_id)
        REFERENCES halls (hall_id) ON DELETE CASCADE
);

CREATE TABLE tickets (
    ticket_id SERIAL PRIMARY KEY,
    session_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    seat_number INTEGER NOT NULL,
    CONSTRAINT tickets_session_id_fkey FOREIGN KEY (session_id)
        REFERENCES cinema_sessions (session_id) ON DELETE CASCADE,
    CONSTRAINT tickets_user_id_fkey FOREIGN KEY (user_id)
        REFERENCES users (user_id) ON DELETE CASCADE
);

-- Data setup scripts
INSERT INTO roles (role_name) VALUES ('admin');
INSERT INTO roles (role_name) VALUES ('user');

INSERT INTO users (username, hashed_password, email, credit_card_info, role_id)
VALUES ('admin', '5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8', 'admin@example.com', '1234567890123456', 1);
