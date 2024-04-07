ALTER TABLE ppo.tbl_document ADD CONSTRAINT fk_tbl_document_markup FOREIGN KEY ( marked_document_id ) REFERENCES ppo.markup( id );
ALTER TABLE ppo.tbl_document ADD CONSTRAINT fk_tbl_document_tbl_user FOREIGN KEY ( sender_id ) REFERENCES ppo.tbl_user( id ) ON DELETE CASCADE;

ALTER TABLE ppo.markup_type ADD CONSTRAINT fk_markup_type_tbl_normcontroller FOREIGN KEY ( created_controller_id ) REFERENCES ppo.tbl_normcontroller( id );

ALTER TABLE ppo.markup ADD CONSTRAINT fk_markup_markup_type FOREIGN KEY ( markup_type_id ) REFERENCES ppo.markup_type( id );
ALTER TABLE ppo.markup ADD CONSTRAINT fk_markup_tbl_normcontroller FOREIGN KEY ( controller_check_id ) REFERENCES ppo.tbl_normcontroller( id );