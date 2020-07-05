User table
CREATE TABLE IF NOT EXISTS Users(
    `id` bigint NOT NULL AUTO_INCREMENT,
    `username` varchar(100) NOT NULL UNIQUE,
    `password` varchar(400) NOT NULL,
    PRIMARY KEY (`id`),
    INDEX(`username`)
)  ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8;