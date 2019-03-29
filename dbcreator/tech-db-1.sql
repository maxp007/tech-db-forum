
drop table if exists "Error";
drop table if exists "Porum";
drop table if exists "Post";
drop table if exists "PostUpdate";
drop table if exists "Status";
drop table if exists "Thread";
drop table if exists "ThreadUpdate";
drop table if exists "User";
drop table if exists "UserUpdate" ;
drop table if exists "Vote";


-- Table: public."Error"

-- DROP TABLE public."Error";

CREATE TABLE public."Error"
(
  message text COLLATE pg_catalog."default"
)
  WITH (
    OIDS = FALSE
  )
  TABLESPACE pg_default;

ALTER TABLE public."Error"
  OWNER to postgres;



-- Table: public."Forum"

-- DROP TABLE public."Forum";

CREATE TABLE public."Forum"
(
  posts bigint,
  slug text COLLATE pg_catalog."default",
  threads integer,
  title text COLLATE pg_catalog."default",
  "user" text COLLATE pg_catalog."default"
)
  WITH (
    OIDS = FALSE
  )
  TABLESPACE pg_default;

ALTER TABLE public."Forum"
  OWNER to postgres;



-- Table: public."Post"

-- DROP TABLE public."Post";

CREATE TABLE public."Post"
(
  author text COLLATE pg_catalog."default",
  created text COLLATE pg_catalog."default",
  forum text COLLATE pg_catalog."default",
  id bigint,
  "isEdited" boolean,
  message text COLLATE pg_catalog."default",
  parent bigint,
  thread integer
)
  WITH (
    OIDS = FALSE
  )
  TABLESPACE pg_default;

ALTER TABLE public."Post"
  OWNER to postgres;



-- Table: public."PostUpdate"

-- DROP TABLE public."PostUpdate";

CREATE TABLE public."PostUpdate"
(
  message text COLLATE pg_catalog."default"
)
  WITH (
    OIDS = FALSE
  )
  TABLESPACE pg_default;

ALTER TABLE public."PostUpdate"
  OWNER to postgres;




-- Table: public."Status"

-- DROP TABLE public."Status";

CREATE TABLE public."Status"
(
  forum integer,
  post bigint,
  thread integer,
  "user" integer
)
  WITH (
    OIDS = FALSE
  )
  TABLESPACE pg_default;

ALTER TABLE public."Status"
  OWNER to postgres;



-- Table: public."Thread"

-- DROP TABLE public."Thread";

CREATE TABLE public."Thread"
(
  author text COLLATE pg_catalog."default",
  created text COLLATE pg_catalog."default",
  forum text COLLATE pg_catalog."default",
  id integer,
  message text COLLATE pg_catalog."default",
  slug text COLLATE pg_catalog."default",
  title text COLLATE pg_catalog."default",
  votes integer
)
  WITH (
    OIDS = FALSE
  )
  TABLESPACE pg_default;

ALTER TABLE public."Thread"
  OWNER to postgres;



-- Table: public."ThreadUpdate"

-- DROP TABLE public."ThreadUpdate";

CREATE TABLE public."ThreadUpdate"
(
  message text COLLATE pg_catalog."default",
  title text COLLATE pg_catalog."default"
)
  WITH (
    OIDS = FALSE
  )
  TABLESPACE pg_default;

ALTER TABLE public."ThreadUpdate"
  OWNER to postgres;




-- Table: public."User"

-- DROP TABLE public."User";

CREATE TABLE public."User"
(
  about text COLLATE pg_catalog."default",
  email text COLLATE pg_catalog."default",
  fullname text COLLATE pg_catalog."default",
  nickname text COLLATE pg_catalog."default"
)
  WITH (
    OIDS = FALSE
  )
  TABLESPACE pg_default;

ALTER TABLE public."User"
  OWNER to postgres;



-- Table: public."UserUpdate"

-- DROP TABLE public."UserUpdate";

CREATE TABLE public."UserUpdate"
(
  about text COLLATE pg_catalog."default",
  email text COLLATE pg_catalog."default",
  fullname text COLLATE pg_catalog."default"
)
  WITH (
    OIDS = FALSE
  )
  TABLESPACE pg_default;

ALTER TABLE public."UserUpdate"
  OWNER to postgres;



-- Table: public."Vote"

-- DROP TABLE public."Vote";

CREATE TABLE public."Vote"
(
  nickname text COLLATE pg_catalog."default",
  voice integer
)
  WITH (
    OIDS = FALSE
  )
  TABLESPACE pg_default;

ALTER TABLE public."Vote"
  OWNER to postgres;