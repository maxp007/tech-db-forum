create function clearalldata() returns SETOF "Status"
  language plpgsql
 as
$$
DECLARE
  BEGIN
  TRUNCATE TABLE public."Error" CASCADE;
  TRUNCATE TABLE public."Vote" CASCADE;
  TRUNCATE TABLE public."forumuser" CASCADE;
  TRUNCATE TABLE public."Post" CASCADE;
  TRUNCATE TABLE public."Thread" CASCADE;

  TRUNCATE TABLE public."Forum" CASCADE;
  TRUNCATE TABLE public."Status" CASCADE;
  TRUNCATE TABLE public."User" CASCADE;
  RETURN QUERY SELECT * FROM public."GetServiceStatus"();
end
$$;

alter function clearalldata() owner to postgres;

