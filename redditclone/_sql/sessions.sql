SET NAMES utf8;
SET time_zone = '+00:00';
SET foreign_key_checks = 0;
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';

DROP TABLE IF EXISTS `sessions`;
CREATE TABLE `sessions` (
                         `id` int(11) NOT NULL AUTO_INCREMENT,
                         `data` text NOT NULL,
                         `userID` text NOT NULL,
                         PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `sessions` (`id`, `data`, `userID`) VALUES
    (0, 'asdfghjkl;', '0');