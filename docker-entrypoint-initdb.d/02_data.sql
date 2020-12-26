INSERT INTO users (login, password, roles)
VALUES ('user1', '$2a$04$zWiNgAN9OXAX9iEFsqYXyuXxtXgEOn5qOIpw4x9Wnp6DJsTAOwnpO', '{ADMIN USER}'),
       ('user2', '$2a$04$.BoH0.ylRB.B9DK7u76m9O/8M7mWuU4NvlF5t6nLUQlK0gJ3C5/aa', '{USER}'),
       ('user3', '$2a$04$v0IPoKBvZ8xUd9dB0YYrme5j6Mw9q.s9xBEXTzYTQw66vbvN0XWMW', '{USER}');

INSERT INTO cards (number, owner, balance)
VALUES ('1111 1111 1111 1111', 1, 1000),
       ('2222 2222 2222 2222', 2, 10000),
       ('3333 3333 3333 3333', 3, 1000000);
