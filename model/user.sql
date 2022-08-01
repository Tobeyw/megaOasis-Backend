CREATE TABLE `user`
(
    `id` int NOT NULL AUTO_INCREMENT,
    `username` varchar(255),
    `bio` varchar(255),
    `address` varchar(255) NOT NULL UNIQUE,
    `email` varchar(255),
    `twitter` varchar(255),
    `avatar` varchar(255),
    `banner` varchar(255),
    `timestamp` int NOT NULL,
    PRIMARY KEY(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;