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

create table metadata (
  id INT AUTO_INCREMENT PRIMARY KEY,
  asset_type INT NOT NULL,
  name varchar(60) NOT NULL,
  datatype varchar(60) NOT NULL
);

insert into metadata (asset_type, name, datatype) values(1, 'chord', 'string');
insert into metadata (asset_type, name, datatype) values(1, 'key', 'string');
insert into metadata (asset_type, name, datatype) values(1, 'bpm', 'int');
insert into metadata (asset_type, name, datatype) values(1, 'name', 'string');

create table asset_data (
  id INT AUTO_INCREMENT PRIMARY KEY,
  asset_uuid varchar(60) not null,
  asset_metadata_id INT NOT NULL,
  value varchar(256) NOT NULL,
  foreign key (asset_metadata_id) references metadata(id),
  index asset_data_asset_uuid (asset_uuid)
);
