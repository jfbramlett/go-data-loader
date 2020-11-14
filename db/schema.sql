--
-- Primary user table
--
create table user (
  id INT AUTO_INCREMENT PRIMARY KEY,
  uuid varchar(60) NOT NULL UNIQUE,
  username varchar(250) NOT NULL UNIQUE,
  password varchar(250) NOT NULL,
  name varchar(250) NOT NULL,
  active TINYINT(1) NOT NULL DEFAULT 1,
  INDEX user_username (username)
);

--
-- Primary pet table
--
create table pet (
  id INT AUTO_INCREMENT PRIMARY KEY,
  uuid varchar(60) NOT NULL UNIQUE,
  name varchar(250) NOT NULL,
  breed varchar(250) NOT NULL,
  birthday date null,
  INDEX user_username (username)
);

--
-- Mapping of users to pets
--
create table user_pet (
  id INT AUTO_INCREMENT PRIMARY KEY,
  user_uuid varchar(60) NOT NULL,
  pet_uuid varchar(60) NOT NULL,
  owner bool default false
  foreign key (user_uuid) references user(uuid),
  foreign key (pet_uuid) references pet(uuid)
);

--
-- Note table
--
create table note (
  id INT AUTO_INCREMENT PRIMARY KEY,
  pet_uuid varchar(60) NOT NULL,
  note text NOT NULL
)

