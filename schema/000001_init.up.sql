
CREATE TABLE mtrx
(
    id      serial       not null unique,
    name    varchar(255) not null unique,
    type 	varchar(255) not null unique,
	value   varchar(255) not null unique,
    delta	varchar(255) not null
);