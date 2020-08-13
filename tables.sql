CREATE TABLE lines
   ( line_name character varying(25) CONSTRAINT line_name_pk PRIMARY KEY,
    line_current_value numeric(5,3),
    line_latest_value numeric(5,3)
   );

INSERT INTO lines (line_name) VALUES ('football');
INSERT INTO lines (line_name) VALUES ('baseball');
INSERT INTO lines (line_name) VALUES ('soccer');

--DROP TABLE lines;