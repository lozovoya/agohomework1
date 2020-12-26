
INSERT INTO users (login, password, full_name, passport, birthdate, status)
VALUES ('iivanov','123','ivan ivanov','aaa111', '2000.12.12', 'ACTIVE'),
       ('ppetrov','456','petr petrov','bbb222', '1990.11.11', 'ACTIVE'),
       ('vvasilev','789','vasya vasilev','ccc333', '1980.10.10', 'ACTIVE');

INSERT INTO cards (number, owner, balance)
VALUES ('1111 1111 1111 1111', 'ivan ivanov', 1000),
       ('2222 2222 2222 2222', 'petr petrov', 10000),
       ('3333 3333 3333 3333', 'vasya vasiliev', 1000000);
