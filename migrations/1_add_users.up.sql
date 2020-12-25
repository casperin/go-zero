create table users (
  id serial primary key,
  email varchar(80) unique not null,
  password varchar(80) not null
);
