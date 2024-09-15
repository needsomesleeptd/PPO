--table users

DROP TABLE IF EXISTS ppo.tbl_user;
CREATE  TABLE ppo.tbl_user ( 
	id                   serial  NOT NULL  ,
	nickname             varchar(100)    ,
	name                 varchar(100)    ,
	surame               varchar    ,
	register_dttm        timestamp    ,
	passwd_rep           varchar(100)    ,
	CONSTRAINT pk_tbl_user PRIMARY KEY ( id )
 );

COMMENT ON TABLE ppo.tbl_user IS 'normal user who can only send data';

--table of documents
DROP TABLE IF EXISTS ppo.tbl_document;
CREATE  TABLE ppo.tbl_document ( 
	id                   serial  NOT NULL  ,
	pages_count          integer    ,
	checks_count         integer    ,
	sender_id            integer    ,
	CONSTRAINT pk_tbl_document PRIMARY KEY ( id )
 );



--table of markups_types
DROP TABLE IF EXISTS ppo.markup_type;
CREATE  TABLE ppo.markup_type ( 
	id                   serial  NOT NULL  ,
	description          text    ,
	created_controller_id integer    ,
	CONSTRAINT pk_markup_type PRIMARY KEY ( id )
 );


--table of markups_instances
DROP TABLE IF EXISTS ppo.markup;

CREATE  TABLE ppo.markup ( 
	id                   serial  NOT NULL  ,
	markup_data          json    ,
	markup_type_id       integer    ,
	marked_document_id   integer    ,
	controller_check_id  integer    ,
	CONSTRAINT pk_markup PRIMARY KEY ( id ),
	CONSTRAINT unq_markup_marked_document_id UNIQUE ( marked_document_id ) 
 );

--table of controllers
DROP TABLE IF EXISTS ppo.tbl_normcontroller;
CREATE  TABLE ppo.tbl_normcontroller ( 
	id                   serial  NOT NULL  ,
	nickname             varchar(100)    ,
	name                 varchar(100)    ,
	surname              varchar(100)    ,
	registration_date    date    ,
	controller_group     varchar(100)    ,
	passwd_rep           varchar(100)    ,
	CONSTRAINT pk_tbl_normcontroller PRIMARY KEY ( id )
 );

