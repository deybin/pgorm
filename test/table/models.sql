create table
	models (
		id uuid PRIMARY KEY NOT NULL,
		document varchar(11) DEFAULT '' NOT NULL,
		nombre text DEFAULT '' NOT NULL,
		address text DEFAULT '' NOT NULL,
		birthdate date NULL,
		age int DEFAULT 0 NOT NULL,
		amount float8 DEFAULT 0.0 NOT NULL,
		credits int DEFAULT 0 NOT NULL,
		passwords text DEFAULT '' NOT NULL,
		key_secret text DEFAULT '' NOT NULL,
		atCreate timestamp DEFAULT now ()
	)

	create table
	models2 (
		id uuid PRIMARY KEY NOT NULL,
		document varchar(11) DEFAULT '' NOT NULL,
		nombre text DEFAULT '' NOT NULL,
		address text DEFAULT '' NOT NULL,
		birthdate date NULL,
		age int DEFAULT 0 NOT NULL,
		amount float8 DEFAULT 0.0 NOT NULL,
		credits int DEFAULT 0 NOT NULL,
		passwords text DEFAULT '' NOT NULL,
		key_secret text DEFAULT '' NOT NULL,
		atCreate timestamp DEFAULT now ()
	)