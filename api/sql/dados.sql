insert into usuarios(nome, nick, email, senha)
values 
("Usuario1", "usuario_1", "usuario1@gmail.com", "$2a$10$F1aBzVZQAsjc2W4sEeYeVu9xaN8lsHrOOcS6eGp9who.3JD/OO3LG"),
("Usuario2", "usuario_2", "usuario2@gmail.com", "$2a$10$F1aBzVZQAsjc2W4sEeYeVu9xaN8lsHrOOcS6eGp9who.3JD/OO3LG"),
("Usuario3", "usuario_3", "usuario3@gmail.com", "$2a$10$F1aBzVZQAsjc2W4sEeYeVu9xaN8lsHrOOcS6eGp9who.3JD/OO3LG");

insert into seguidores(usuario_id, seguidor_id)
values 
(1, 2),
(1, 3),
(3, 1);

insert into publicacoes(titulo, conteudo, autor_id)
values
("Publicacao 1", "Publicacao 1", 1 ),
("Publicacao 2", "Publicacao 2", 2 ),
("Publicacao 3", "Publicacao 3", 3 );