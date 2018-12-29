CREATE table IF NOT EXISTS user(
    id int(64),
    username varchar(60) unique key not null,
    password varchar(100) not null,
    nickname varchar(50),
    avatar varchar(200),
    email varchar(200) unique key,
    ctime varchar(100),
    constraint primary key(id)
)