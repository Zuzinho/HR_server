create table users (
    id serial primary key,
    name text not null,
    surname text not null,
    patronymic text null,
    age int null,
    gender text null,
    nation text null
)