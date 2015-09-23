
create table car (
	id serial primary key,
	model varchar(30) not null,
	price varchar(50),
	year varchar(40),
	brand varchar(40) unique not null
);

insert into car (model, price, year, brand) 
values ('silverado', '39,000', '2015', 'Chevrolet');

