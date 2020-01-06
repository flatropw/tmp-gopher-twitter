CREATE table users (
    id SERIAL PRIMARY KEY,
    login varchar UNIQUE ,
    email varchar UNIQUE,
    password varchar,
    token varchar NULL
)

create table tweets
(
	id serial not null
		constraint tweets_pk
			primary key,
	message text not null,
	user_id int not null
		constraint tweets_users_id_fk
			references users,
	created_at bigint default 0
);

create table subscriptions
(
	id int not null
		constraint subscriptions_pk
			primary key,
	subscriber_id int not null
		constraint subscriptions_users_id_fk
			references users
				on delete cascade,
	subscribed_id int not null
		constraint subscriptions_users_id_fk_2
			references users
				on delete cascade
);

