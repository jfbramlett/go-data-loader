use test;

create table sample (
  id INT AUTO_INCREMENT PRIMARY KEY,
  uuid varchar(60) NOT NULL UNIQUE,
  asset_uuid varchar(60) NOT NULL,
  chord varchar(10),
  skey varchar(10),
  bpm INT,
  name varchar(60) NOT NULL,
  INDEX test_asset_uuid (asset_uuid)
);
