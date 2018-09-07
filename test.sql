use entrytask;

-- create table if not exists user(
-- 	id int unique auto_increment,
--     username varchar(20),
--     password varchar(30) not null,
--     nickname varchar(20) not null,
--     avatar varchar(200),
--     primary key (username)
--     
-- )

DELIMITER $$

CREATE PROCEDURE BatchInsert(IN init INT, IN user_count INT)
DETERMINISTIC
BEGIN
      DECLARE id int;
      DECLARE var int;

      SET var = 0;
      SET id = init;
     
      
      WHILE var < user_count DO
          INSERT INTO USER(id,username, password, nickname, avatar) VALUES ( null,CONCAT('mhh', id), "a123456",concat('dou',id) ,"/go1.png" );
          SET id = id + 1;
          SET var = var + 1;
      END WHILE;
      
END$$

DELIMITER ;
all BatchInsert(1,1000)
--  call BatchInsert(1, 10000000)
