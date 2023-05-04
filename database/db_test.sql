-- Get all movies in the database:
SELECT * FROM movies;

-- Get all cinema sessions for a specific movie:
SELECT * FROM cinema_sessions WHERE movie_id = 1;

--Get all tickets for a specific cinema session:
SELECT * FROM tickets WHERE session_id = 1;

--Get all tickets for a specific user:
SELECT * FROM tickets WHERE user_id = 3;

--Get the number of available seats for a specific cinema session:
SELECT halls.capacity - COUNT(tickets.ticket_id) as available_seats
FROM cinema_sessions
INNER JOIN halls ON cinema_sessions.hall_id = halls.hall_id
LEFT JOIN tickets ON cinema_sessions.session_id = tickets.session_id
WHERE cinema_sessions.session_id = 1
GROUP BY halls.capacity;




