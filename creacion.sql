/* Drop Tables */

DROP TABLE IF EXISTS vehiculo CASCADE
;

/* Create Tables */
CREATE TABLE contribuyente
(
id bigint(20) NOT NULL AUTO_INCREMENT,
placa varchar(6) NOT NULL,
marca varchar(50) NOT NULL,
linea varchar(70) NOT NULL,
modelo varchar(4) NOT NULL,
PRIMARY KEY (id)
)
;