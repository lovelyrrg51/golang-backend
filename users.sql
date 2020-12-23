CREATE TABLE `users` (
  `id` int(6) unsigned NOT NULL AUTO_INCREMENT,
  `user_name` varchar(30) NOT NULL,
  `name` varchar(30) NOT NULL,
  `authToken` varchar(30) NOT NULL,
  PRIMARY KEY (`id`)
) 