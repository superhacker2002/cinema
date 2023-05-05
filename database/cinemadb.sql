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
        REFERENCES roles (role_id) ON DELETE CASCADE
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
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
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
VALUES ('admin', '6D4525C2A21F9BE1CCA9E41F3AA402E0765EE5FCC3E7FEA34A169B1730AE386E', 'admin@example.com', '1234567890123456', 1);

INSERT INTO halls (hall_name, capacity) VALUES ('Hall 1', 100);
INSERT INTO halls (hall_name, capacity) VALUES ('Hall 2', 80);
