CREATE table users (
    id SERIAL PRIMARY KEY,
    login varchar UNIQUE ,
    email varchar UNIQUE,
    password varchar,
    token varchar NULL
)