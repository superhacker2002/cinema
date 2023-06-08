INSERT INTO users (username, hashed_password, email, credit_card_info, role_id)
VALUES ('user1', 'A3163B169544206384021627139043454DD8C7D926746F6D01A11FA904D90C03', 'user1@example.com', '1111222233334444', 2),
       ('user2', '5ABFAC4EC9F3459E7FA7C22476615FBA9F2E98125C3D38FDA867993D30735CD8', 'user2@example.com', '2222333344445555', 2);

INSERT INTO movies (title, genre, release_date, duration)
VALUES ('Avengers: Endgame', 'Action, Adventure, Drama', '2019-04-26', 181),
       ('The Shawshank Redemption', 'Drama', '1994-09-23', 142),
       ('The Lion King', 'Animation, Adventure, Drama', '2019-07-19', 118),
       ('The Dark Knight', 'Action, Crime, Drama', '2008-07-18', 152),
       ('Forrest Gump', 'Drama, Romance', '1994-07-06', 142),
       ('The Shawshank Redemption', 'Drama', '1994-09-23', 142);

INSERT INTO halls(hall_name, capacity)
VALUES ('Small Hall', 10),
       ('Big Hall', 100),
       ('Very Big Hall', 200),
       ('IMAX', 100);

INSERT INTO cinema_sessions (movie_id, hall_id, start_time, end_time, price)
VALUES (3, 1, '2023-05-29 14:00:00 +04', '2023-05-29 16:00:00 +04', 10.00),
       (3, 2, '2023-05-29 14:00:00 +04', '2023-05-29 16:00:00 +04', 10.00),
       (3, 3, '2023-05-29 14:00:00 +04', '2023-05-29 16:00:00 +04', 10.00),
       (3, 4, '2023-05-29 14:00:00 +04', '2023-05-29 16:00:00 +04', 10.00),
       (3, 1, '2023-05-29 17:00:00 +04', '2023-05-29 19:00:00 +04', 10.00),
       (3, 2, '2023-05-29 17:00:00 +04', '2023-05-29 19:00:00 +04', 10.00),
       (3, 3, '2023-05-29 17:00:00 +04', '2023-05-29 19:00:00 +04', 10.00),
       (3, 4, '2023-05-29 17:00:00 +04', '2023-05-29 19:00:00 +04', 10.00),
       (3, 1, '2023-05-22 08:00:00 +04', '2023-05-22 10:00:00 +04', 10.00);


INSERT INTO tickets (session_id, user_id, seat_number)
VALUES (1, 2, 1),
       (1, 3, 2),
       (1, 3, 3);