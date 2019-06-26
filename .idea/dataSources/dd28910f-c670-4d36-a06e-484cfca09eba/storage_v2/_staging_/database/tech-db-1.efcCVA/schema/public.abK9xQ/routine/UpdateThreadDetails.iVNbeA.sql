create function "UpdateThreadDetails"(thread_id integer, thread_slug text, new_title text, new_message text) returns SETOF "Thread"
  language plpgsql
as
$$
DECLARE
thread_counter integer;
treadId BIGINT;
BEGIN
if thread_id = 0 then
    SELECT count(*) INTO thread_counter from public."Thread" where slug = thread_slug::citext;
    SELECT id INTO treadId from public."Thread" where slug = thread_slug::citext;
  end if;

  -- CHECK IF THREAD BY ID EXISTS
  if thread_slug = '' then
    SELECT count(*) INTO thread_counter from public."Thread" where id = thread_id;
    SELECT id INTO treadId from public."Thread" where id = thread_id;
  end if;

  if thread_counter = 0 OR thread_counter IS NULL then
    RAISE EXCEPTION 'THREAD NOT FOUND  slug  %, id %', thread_slug::text,thread_id::text USING ERRCODE = 'no_data_found';
  end if;
if  treadId IS NULL then
  RAISE EXCEPTION 'THREAD NOT FOUND  slug  %, id %', thread_slug::text,thread_id::text USING ERRCODE = 'no_data_found';

end if;
        if "new_message" <> '' then
             UPDATE "Thread" SET message = COALESCE("new_message"::citext, message) WHERE id = treadId;
        end if;

        if "new_title" <> '' then
              UPDATE "Thread"  SET title = COALESCE("new_title"::citext, title) WHERE id = treadId;
        end if;

        RETURN QUERY SELECT * FROM public."Thread" Where  id = treadId;

END
$$;

alter function "UpdateThreadDetails"(integer, text, text, text) owner to postgres;

