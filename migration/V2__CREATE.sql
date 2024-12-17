CREATE TABLE delivery."user" (
	id uuid NOT NULL,
	login varchar NOT NULL,
	"password" varchar NOT NULL,
	username varchar NOT NULL,
	"role" varchar NULL,
	refresh_token varchar NULL,
	CONSTRAINT user_pkey PRIMARY KEY (id),
	CONSTRAINT user_unique UNIQUE (login)
);

CREATE TABLE delivery.courier (
	id uuid NOT NULL,
	"name" varchar NULL,
	surname varchar NULL,
	status varchar NULL,
	performance_indicator int4 NULL,
	userid uuid NOT NULL,
	CONSTRAINT courier_pkey PRIMARY KEY (id),
	CONSTRAINT courier_userid_fkey FOREIGN KEY (userid) REFERENCES delivery."user"(id)
);

CREATE TABLE delivery.delivery (
	id uuid NOT NULL,
	courier_id uuid NULL,
	delivery_date varchar NOT NULL,
	delivery_status varchar NOT NULL,
	delivery_comment varchar NULL,
	CONSTRAINT delivery_pk PRIMARY KEY (id)
);