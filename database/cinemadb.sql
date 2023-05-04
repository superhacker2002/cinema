CREATE TABLE roles (
    role_id SERIAL PRIMARY KEY,
    role_name VARCHAR(50) NOT NULL
);

CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    hashed_password VARCHAR(100) NOT NULL,
    email VARCHAR(50) NOT NULL,
    credit_card_info VARCHAR(50) NOT NULL,
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
    seat_number INTEGER UNIQUE NOT NULL,
    CONSTRAINT tickets_session_id_fkey FOREIGN KEY (session_id)
        REFERENCES cinema_sessions (session_id) ON DELETE CASCADE,
    CONSTRAINT tickets_user_id_fkey FOREIGN KEY (user_id)
        REFERENCES users (user_id) ON DELETE CASCADE
);

-- Data setup scripts
INSERT INTO roles (role_name) VALUES ('admin');
INSERT INTO roles (role_name) VALUES ('user');

INSERT INTO users (username, hashed_password, email, credit_card_info, role_id)
VALUES ('admin', 'hashed_admin_password', 'admin@example.com', '1234567890123456', 1);

INSERT INTO users (username, hashed_password, email, credit_card_info, role_id)
VALUES ('user1', 'hashed_user1_password', 'user1@example.com', '1111222233334444', 2);

INSERT INTO users (username, hashed_password, email, credit_card_info, role_id)
VALUES ('user2', 'hashed_user2_password', 'user2@example.com', '2222333344445555', 2);

INSERT INTO movies (title, genre, release_date, duration)
VALUES ('Avengers: Endgame', 'Action, Adventure, Drama', '2019-04-26', 181);

INSERT INTO movies (title, genre, release_date, duration)
VALUES ('The Lion King', 'Animation, Adventure, Drama', '2019-07-19', 118);

INSERT INTO movies (title, genre, release_date, duration)
VALUES ('The Dark Knight', 'Action, Crime, Drama', '2008-07-18', 152);

INSERT INTO movies (title, genre, release_date, duration)
VALUES ('Forrest Gump', 'Drama, Romance', '1994-07-06', 142);

INSERT INTO movies (title, genre, release_date, duration)
VALUES ('The Shawshank Redemption', 'Drama', '1994-09-23', 142);

INSERT INTO halls (hall_name, capacity) VALUES ('Hall 1', 100);
INSERT INTO halls (hall_name, capacity) VALUES ('Hall 2', 80);

INSERT INTO cinema_sessions (movie_id, hall_id, start_time, end_time, price)
VALUES (1, 1, '2023-05-10 19:00:00', '2023-05-10 22:00:00', 10.00);

INSERT INTO cinema_sessions (movie_id, hall_id, start_time, end_time, price)
VALUES (1, 2, '2023-05-10 20:00:00', '2023-05-10 22:00:00', 10.00);

INSERT INTO tickets (session_id, user_id, seat_number) VALUES (1, 2, 1);
INSERT INTO tickets (session_id, user_id, seat_number) VALUES (1, 3, 2);
INSERT INTO tickets (session_id, user_id, seat_number) VALUES (1, 3, 3);

