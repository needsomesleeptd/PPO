-- note that it must be run with automigrate commands

ALTER TABLE markups ADD CONSTRAINT fk_markup_markup_type FOREIGN KEY ( class_label ) REFERENCES markup_types( id );

ALTER TABLE documents ADD CONSTRAINT fk_document_user FOREIGN KEY ( creator_id ) REFERENCES users( id );

ALTER TABLE documents ADD CONSTRAINT fk_document_user FOREIGN KEY ( creator_id ) REFERENCES users( id );

ALTER TABLE markups ADD CONSTRAINT fk_markup_user FOREIGN KEY ( creator_id ) REFERENCES users( id );

ALTER TABLE markup_types ADD CONSTRAINT fk_markup_types_user FOREIGN KEY ( creator_id ) REFERENCES users( id );


