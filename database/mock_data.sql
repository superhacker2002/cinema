INSERT INTO users (username, hashed_password, email, credit_card_info, role_id)
VALUES ('user1', 'A3163B169544206384021627139043454DD8C7D926746F6D01A11FA904D90C03', 'user1@example.com', '1111222233334444', 2);

INSERT INTO users (username, hashed_password, email, credit_card_info, role_id)
VALUES ('user2', '5ABFAC4EC9F3459E7FA7C22476615FBA9F2E98125C3D38FDA867993D30735CD8', 'user2@example.com', '2222333344445555', 2);

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

INSERT INTO cinema_sessions (movie_id, hall_id, start_time, end_time, price)
VALUES (1, 1, '2023-05-10 19:00:00', '2023-05-10 22:00:00', 10.00);

INSERT INTO cinema_sessions (movie_id, hall_id, start_time, end_time, price)
VALUES (1, 2, '2023-05-10 20:00:00', '2023-05-10 22:00:00', 10.00);

INSERT INTO cinema_sessions (movie_id, hall_id, start_time, end_time, price)
VALUES (3, 2, '2023-05-17 19:00:00', '2023-05-10 22:00:00', 10.00);

INSERT INTO cinema_sessions (movie_id, hall_id, start_time, end_time, price)
VALUES (3, 2, '2023-05-18 20:00:00', '2023-05-10 22:00:00', 10.00);

INSERT INTO tickets (session_id, user_id, seat_number) VALUES (1, 2, 1);
INSERT INTO tickets (session_id, user_id, seat_number) VALUES (1, 3, 2);
INSERT INTO tickets (session_id, user_id, seat_number) VALUES (1, 3, 3);