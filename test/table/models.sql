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
		email text DEFAULT '' NOT NULL,
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
		email text DEFAULT '' NOT NULL,
		passwords text DEFAULT '' NOT NULL,
		key_secret text DEFAULT '' NOT NULL,
		atCreate timestamp DEFAULT now ()
	)

	INSERT INTO models (id,"document",nombre,address,birthdate,age,amount,credits,passwords,key_secret,atcreate,email) VALUES
	 ('550e8400-e29b-41d4-a716-446655440001'::uuid,'12345678901','Juan Perez','Lima','1995-05-10',29,150.5,2,'hash1','secret1','2025-01-01 10:00:00.000',''),
	 ('550e8400-e29b-41d4-a716-446655440002'::uuid,'12345678902','Maria Lopez','Arequipa','1990-08-20',34,0.0,0,'hash2','secret2','2025-01-02 11:30:00.000',''),
	 ('550e8400-e29b-41d4-a716-446655440003'::uuid,'12345678903','Carlos Ruiz','Cusco',NULL,40,1200.0,5,'hash3','secret3','2025-01-03 09:15:00.000',''),
	 ('550e8400-e29b-41d4-a716-446655440004'::uuid,'12345678904','Ana Torres','Lima','2000-01-01',25,300.75,1,'hash4','secret4','2025-01-04 14:45:00.000',''),
	 ('550e8400-e29b-41d4-a716-446655440005'::uuid,'12345678905','Pedro Sanchez','Piura','1985-12-12',39,50.0,0,'hash5','secret5','2025-01-05 08:00:00.000',''),
	 ('550e8400-e29b-41d4-a716-446655440006'::uuid,'12345678906','Lucia Diaz','Trujillo',NULL,22,0.0,0,'hash6','secret6','2025-01-06 16:20:00.000',''),
	 ('550e8400-e29b-41d4-a716-446655440007'::uuid,'12345678907','Miguel Castro','Chiclayo','1998-03-18',26,800.0,3,'hash7','secret7','2025-01-07 12:10:00.000',''),
	 ('550e8400-e29b-41d4-a716-446655440008'::uuid,'12345678908','Sofia Ramirez','Lima','1992-11-05',32,100.0,1,'hash8','secret8','2025-01-08 18:50:00.000',''),
	 ('550e8400-e29b-41d4-a716-446655440009'::uuid,'12345678909','Diego Flores','Tacna',NULL,45,2000.0,10,'hash9','secret9','2025-01-09 07:40:00.000',''),
	 ('550e8400-e29b-41d4-a716-446655440010'::uuid,'12345678910','Valeria Mendoza','Iquitos','2002-06-25',22,25.0,0,'hash10','secret10','2025-01-10 20:00:00.000','');
