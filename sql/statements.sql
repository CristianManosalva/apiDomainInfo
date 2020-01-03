DROP TABLE IF EXISTS domain_info CASCADE;
CREATE TABLE domain_info (
	domain STRING PRIMARY KEY,
	servers_changed BOOL,
	ssl_grade STRING,
	previous_ssl_grade STRING,
	logo STRING,
	title STRING,
	is_down BOOL
);

DROP TABLE IF EXISTS servers_info CASCADE;
CREATE TABLE servers_info (
	address STRING,
	ssl_grade STRING,
	country STRING,
	owner STRING,
	domain STRING
);


