CREATE TABLE `message`
(
    `id` int NOT NULL AUTO_INCREMENT,
    `email` varchar(255),
    `event` varchar(255) NOT NULL,
    `title` varchar(255) NOT NULL,
    `message` text NOT NULL,
    `timestamp` int NOT NULL,
    PRIMARY KEY(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;