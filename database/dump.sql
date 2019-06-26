
create extension if not exists "uuid-ossp";

create extension if not exists "citext";

create type users_query_result as
  (
  users "User"[],
  status integer
  );

alter type users_query_result owner to postgres;

-- Unknown how to generate base type type

alter type citext owner to postgres;

create type type_in_post as
  (
  thread_id integer,
  thread_slug citext,
  author citext,
  messag text,
  parent integer
  );

alter type type_in_post owner to postgres;

create type id_parent_tuple as
  (
  post_id integer,
  parent_id integer
  );

alter type id_parent_tuple owner to postgres;

create type treenode as
  (
  id bigint,
  parent bigint,
  path bigint[]
  );

alter type treenode owner to postgres;

create type servicestatus as
  (
  forum integer,
  post bigint,
  thread integer,
  "user" integer
  );

alter type servicestatus owner to postgres;

create table "Error"
(
  message text
);

alter table "Error"
  owner to postgres;

create table "Status"
(
  forum  integer not null,
  post   bigint  not null,
  thread integer not null,
  "user" integer not null
);

alter table "Status"
  owner to postgres;

create table "User"
(
  about    text,
  email    citext not null
    constraint email
      unique,
  fullname text   not null,
  nickname citext
    constraint nickname
      unique
);

alter table "User"
  owner to postgres;

create table "Forum"
(
  posts   bigint,
  slug    citext not null,
  threads integer,
  title   text   not null,
  "user"  citext not null
    constraint "user"
      references "User" (nickname)
      on update cascade on delete cascade
);

comment on constraint "user" on "Forum" is 'Foreign Key, никнейм  создателя форума';

alter table "Forum"
  owner to postgres;

create unique index forum_slug_uindex
  on "Forum" (slug);

create table "Post"
(
  author     citext not null
    constraint author
      references "User" (nickname)
      on update cascade on delete cascade,
  created    timestamp with time zone,
  forum      citext
    constraint forum
      references "Forum" (slug)
      on update cascade on delete cascade,
  id         bigserial,
  "isEdited" boolean,
  message    text   not null,
  parent     bigint,
  thread     integer
);

alter table "Post"
  owner to postgres;

create unique index post_id_uindex
  on "Post" (id);

create index parent_index
  on "Post" (parent);

create index full_index
  on "Post" (id, parent, thread);

create index thread_id_index
  on "Post" (thread, id);

create table "Thread"
(
  author  citext not null
    constraint author
      references "User" (nickname)
      on update cascade on delete cascade,
  created timestamp with time zone,
  forum   citext
    constraint forum
      references "Forum" (slug)
      on update cascade on delete cascade,
  id      serial,
  message text   not null,
  slug    citext not null,
  title   text   not null,
  votes   integer
);

alter table "Thread"
  owner to postgres;

create unique index thread_slug_uindex
  on "Thread" (slug);

create unique index thread_id_uindex
  on "Thread" (id);

create unique index user_email_uindex
  on "User" (email);

create unique index user_nickname_uindex
  on "User" (nickname);

create table "Vote"
(
  nickname citext  not null
    constraint nickname
      references "User" (nickname)
      on update cascade on delete cascade,
  voice    integer not null,
  thread   citext
    constraint thread
      references "Thread" (slug)
      on update cascade on delete cascade
);

alter table "Vote"
  owner to postgres;

create table forumuser
(
  forum  citext not null
    constraint forum
      references "Forum" (slug)
      on update cascade on delete cascade,
  "user" citext not null
    constraint "user"
      references "User" (nickname)
      on update cascade on delete cascade
);

alter table forumuser
  owner to postgres;

create or replace function "ForumPostInc"() returns trigger
  language plpgsql
as
$$
DECLARE
  postCounter "Forum"%ROWTYPE;
BEGIN
  SELECT * INTO postCounter FROM "Forum" WHERE slug = NEW.forum;
  postCounter.posts := postCounter.posts + 1;
  RETURN NEW;
END;
$$;

alter function "ForumPostInc"() owner to postgres;

create or replace function "CreateOrGetUsers"("aboutArg" text, "emailArg" text, "fullnameArg" text,
                                              "nicknameArg" text) returns SETOF "User"
  language plpgsql
as
$$
DECLARE

BEGIN
  RETURN QUERY SELECT * FROM public."User" Where nickname = "nicknameArg"::citext OR email = "emailArg"::citext;
  IF NOT FOUND THEN
    INSERT INTO public."User" (about, email, fullname, nickname)
    VALUES ("aboutArg"::citext, "emailArg"::citext, "fullnameArg"::citext, "nicknameArg"::citext);
    RETURN QUERY SELECT * FROM public."User" Where nickname = "nicknameArg"::citext AND email = "emailArg"::citext;
  ELSE
    RAISE unique_violation;
  END IF;
END
$$;

alter function "CreateOrGetUsers"(text, text, text, text) owner to postgres;

create or replace function "InsertUser"("aboutArg" text, "emailArg" text, "fullnameArg" text,
                                        "nicknameArg" text) returns SETOF "User"
  language plpgsql
as
$$
DECLARE

BEGIN
  RETURN QUERY SELECT * FROM public."User" Where nickname = "nicknameArg"::citext OR email = "emailArg"::citext;
  IF NOT FOUND THEN
    INSERT INTO public."User" (about, email, fullname, nickname)
    VALUES ("aboutArg"::citext, "emailArg"::citext, "fullnameArg"::citext, "nicknameArg"::citext);
    RETURN QUERY SELECT * FROM public."User" Where nickname = "nicknameArg"::citext AND email = "emailArg"::citext;
  ELSE
    RAISE unique_violation;
  END IF;
END
$$;

alter function "InsertUser"(text, text, text, text) owner to postgres;

create or replace function "GetSimilarUsers"("emailArg" text, "nicknameArg" text) returns SETOF "User"
  language plpgsql
as
$$
DECLARE
BEGIN
  RETURN QUERY SELECT * FROM public."User" Where nickname = "nicknameArg"::citext OR email = "emailArg"::citext;
END
$$;

alter function "GetSimilarUsers"(text, text) owner to postgres;

create or replace function "UpdateUserProfile"("aboutArg" text, "emailArg" text, "fullnameArg" text,
                                               "nicknameArg" text) returns SETOF "User"
  language plpgsql
as
$$
DECLARE
  v_count integer := 0;
BEGIN
  SELECT count(*) into v_count FROM public."User" Where nickname = "nicknameArg"::citext;

  IF v_count = 0 THEN
    RAISE EXCEPTION no_data_found;
  ELSE
    if "emailArg" <> '' then
      UPDATE public."User" SET email = COALESCE("emailArg"::citext, email) WHERE nickname = "nicknameArg"::citext;
    end if;

    if "aboutArg" <> '' then
      UPDATE public."User" SET about = COALESCE("aboutArg"::citext, email) WHERE nickname = "nicknameArg"::citext;
    end if;

    if "fullnameArg" <> '' then
      UPDATE public."User" SET fullname = COALESCE("fullnameArg"::citext, email) WHERE nickname = "nicknameArg"::citext;
    end if;

    RETURN QUERY SELECT * FROM public."User" Where nickname = "nicknameArg"::citext;
  END IF;

END
$$;

alter function "UpdateUserProfile"(text, text, text, text) owner to postgres;

create or replace function "CreateOrGetForum"("titleArg" text, "slugArg" text, "userArg" citext) returns SETOF "Forum"
  language plpgsql
as
$$
DECLARE
  forumRow      "Forum"%ROWTYPE;
  forum_counter integer := 0;
  userRow       public."User"%ROWTYPE;
BEGIN
  SELECT count(*) INTO forum_counter FROM public."Forum" WHERE slug = "slugArg"::citext;
  if forum_counter <> 0 THEN
    RAISE unique_violation;
  ELSE
    SELECT * INTO userRow FROM public."User" WHERE nickname = "userArg"::citext;
    INSERT INTO public."Forum" ("posts", "slug", "threads", "title", "user")
    VALUES (0, "slugArg", 0, "titleArg"::citext, userRow.nickname::citext);
    RETURN QUERY SELECT * FROM public."Forum" WHERE slug = "slugArg"::citext;

  end if;
END
$$;

alter function "CreateOrGetForum"(text, text, citext) owner to postgres;



create or replace function "CreateOrGetThread"("slugArg" citext, "authorArg" citext, "createdArg" text,
                                               "messageArg" text, "forumArg" citext,
                                               "titleArg" text) returns SETOF "Thread"
  language plpgsql
as
$$
DECLARE
  threadRow      "Thread"%ROWTYPE;
  thread_counter integer := 0;
  userRow        public."User"%ROWTYPE;
  user_counter   integer := 0;
  time_stamp     timestamptz;
  thread_slug    text;
  forum_row      citext;
  forum_counter  integer := 0;
  temp_counter   INTEGER := 0;
BEGIN
  SELECT count(*) INTO thread_counter FROM public."Thread" WHERE slug = "slugArg"::citext;
  if thread_counter <> 0 THEN
    RAISE unique_violation;
    return;
  end if;

  SELECT count(*) INTO user_counter FROM public."User" WHERE nickname = "authorArg"::citext;
  if user_counter = 0 THEN
    RAISE no_data_found;
    return;
  end if;

  if "createdArg" = '' then
    time_stamp = now();
  else
    time_stamp = "createdArg";
  end if;

  if "slugArg" = '' then
    thread_slug = uuid_generate_v4()::citext;
    if thread_slug = '' then
      thread_slug = 'uid_generate_v4()';
    end if;
  else
    thread_slug = "slugArg";
  end if;

  SELECT "Forum".slug INTO forum_row FROM public."Forum" WHERE slug = "forumArg"::citext;
  SELECT count(*) INTO forum_counter FROM public."Forum" WHERE slug = "forumArg"::citext;


  if forum_counter = 0 then
    raise no_data_found;
    return;
  end if;

  SELECT * INTO userRow FROM public."User" WHERE nickname = "authorArg"::citext;
  INSERT INTO public."Thread" (author, created, forum, "message", slug, title, votes)
  VALUES (userRow.nickname::citext, time_stamp::timestamptz, forum_row::citext, "messageArg", thread_slug::citext,
          "titleArg", 0);


  temp_counter :=
      (SELECT count(*) FROM public.forumuser WHERE forum = forum_row::citext AND "user" = userRow.nickname::citext);
  if temp_counter = 0 then
    INSERT INTO public.forumuser ("forum", "user") VALUES (forum_row::citext, userRow.nickname::citext);
  end if;
  UPDATE public."Forum" SET threads = threads + 1 where slug = forum_row::citext;
  RETURN QUERY SELECT * FROM public."Thread" WHERE slug = thread_slug::citext;

END
$$;

alter function "CreateOrGetThread"(citext, citext, text, text, citext, text) owner to postgres;

create or replace function "CreateOrGetVote"(thread_slug citext, thread_id integer, author citext, vote integer) returns SETOF "Thread"
  language plpgsql
as
$$
DECLARE
  thread_counter integer := 0;
  author_counter integer := 0;
  author_row     "User"%ROWTYPE;
  treadSlug      citext;
  vote_counter   integer := 0;
  vote_voice     integer := 0;
BEGIN


  -- CHECK IF THREAD BY SLUG EXISTS
  if thread_id = 0 then
    SELECT count(*) INTO thread_counter from public."Thread" where slug = thread_slug::citext;
    SELECT slug INTO treadSlug from public."Thread" where slug = thread_slug::citext;
  end if;

  -- CHECK IF THREAD BY ID EXISTS
  if thread_slug = '' then
    SELECT count(*) INTO thread_counter from public."Thread" where id = thread_id;
    SELECT slug INTO treadSlug from public."Thread" where id = thread_id;
  end if;

  if thread_counter = 0 then
    RAISE EXCEPTION 'THREAD NOT FOUND  slug  %, id %', thread_slug::text,thread_id::text USING ERRCODE = 'no_data_found';
  end if;

  if vote = 0 AND author = '' then
    RETURN QUERY SELECT * from public."Thread" WHERE slug = treadSlug::citext;
  end if;

  if (vote <> -1 AND vote <> 1) then
    RAISE EXCEPTION 'WRONG VOTE VALUE  %', vote::text USING ERRCODE = 'no_data_found';
    return;
  end if;

  SELECT count(*) INTO author_counter from public."User" where nickname = author::citext;
  SELECT * INTO author_row from public."User" where nickname = author::citext;
  if author_counter = 0 then
    RAISE EXCEPTION 'AUTHOR NOT FOUND  nickname %,', author::text USING ERRCODE = 'no_data_found';
    return;
  end if;
  SELECT count(*) INTO vote_counter
  from public."Vote"
  where ("thread" = treadSlug::citext AND "nickname" = author::citext);

  if vote_counter = 0 then
    INSERT INTO public."Vote" ("nickname", "voice", "thread") VALUES (author, vote, treadSlug);
    UPDATE public."Thread" SET "votes"="votes" + vote::integer WHERE "slug" = treadSlug::citext;
  else
    SELECT voice INTO vote_voice
    from public."Vote"
    where ("thread" = treadSlug::citext AND "nickname" = author::citext);

    if vote_voice::integer <> vote::integer then
      UPDATE public."Vote" SET "voice" = vote WHERE "thread" = treadSlug::citext AND "nickname" = author::citext;
      UPDATE public."Thread" SET "votes"="votes" + vote + vote::integer WHERE "slug" = treadSlug::citext;
    end if;
  end if;

  RETURN QUERY SELECT * from public."Thread" WHERE slug = treadSlug::citext;

END
$$;

alter function "CreateOrGetVote"(citext, integer, citext, integer) owner to postgres;

create or replace function "GetThreadDetails"(thread_slug citext, thread_id integer) returns SETOF "Thread"
  language plpgsql
as
$$
DECLARE
  treadSlug      citext;
  thread_counter integer := 0;
BEGIN
  if thread_id = 0 then
    SELECT count(*) INTO thread_counter from public."Thread" where slug = thread_slug::citext;
    SELECT slug INTO treadSlug from public."Thread" where slug = thread_slug::citext;
  end if;

  -- CHECK IF THREAD BY ID EXISTS
  if thread_slug = '' then
    SELECT count(*) INTO thread_counter from public."Thread" where id = thread_id;
    SELECT slug INTO treadSlug from public."Thread" where id = thread_id;
  end if;

  if thread_counter = 0 then
    RAISE EXCEPTION 'THREAD NOT FOUND  slug  %, id %', thread_slug::text,thread_id::text USING ERRCODE = 'no_data_found';
  end if;

  RETURN QUERY SELECT * FROM public."Thread" WHERE slug = treadSlug::citext;
END
$$;

alter function "GetThreadDetails"(citext, integer) owner to postgres;

create or replace function "ThreadPostsTree"(thread_id integer, limitarg integer, sincearg integer, sortarg text,
                                             descarg boolean) returns SETOF "Post"
  language plpgsql
as
$$
DECLARE
  parents_counter  integer := 0;
  parents_id_array integer[];
  result_array     "Post"[];
  result_id_array  bigint[];
  temp_post_row    "Post";
BEGIN
  parents_id_array := ARRAY(SELECT id
                            FROM public."Post"
                            where (id > 448 AND id < 508 AND parent = 0 AND "Post".thread = thread_id)
                            ORDER BY id ASC);
  SELECT * INTO temp_post_row FROM "Post" where id = parents_id_array [ 1];
  if limitarg > 0 then
    result_id_array = array_append(result_id_array, temp_post_row.id);
  end if;
  FOR i in 1..array_length(parents_id_array, 1)
    LOOP
      if array_length(parents_id_array, 1) = limitarg then
        Return QUERY SELECT * FROM public."Post" WHERE id = ANY (result_id_array);
        EXIT;
      else
        SELECT * INTO temp_post_row FROM public."Post" where parent = parents_id_array [ i];
      end if;
    end loop;
  SELECT * FROM public."Post" WHERE id = ANY (ARRAY [448::bigint,449::bigint,507::bigint]);


END
$$;

alter function "ThreadPostsTree"(integer, integer, integer, text, boolean) owner to postgres;

create or replace function "ThreadPostsTree"(thread_id integer, limitarg integer, sincearg integer) returns SETOF "Post"
  language plpgsql
as
$$
DECLARE
  parents_counter  integer := 0;
  parents_id_array integer[];
  result_array     "Post"[];
  result_id_array  bigint[];
  temp_post_row    "Post";
BEGIN
  parents_id_array := ARRAY(SELECT id
                            FROM public."Post"
                            where (id > 448 AND id < 508 AND parent = 0 AND "Post".thread = thread_id)
                            ORDER BY id ASC);
  SELECT * INTO temp_post_row FROM "Post" where id = parents_id_array [ 1];
  if limitarg > 0 then
    result_id_array = array_append(result_id_array, temp_post_row.id);
  end if;
  FOR i in 1..array_length(parents_id_array, 1)
    LOOP
      if array_length(parents_id_array, 1) = limitarg then
        Return QUERY SELECT * FROM public."Post" WHERE id = ANY (result_id_array);
        EXIT;
      else
        SELECT * INTO temp_post_row FROM public."Post" where parent = parents_id_array [ i];
      end if;
    end loop;
  SELECT * FROM public."Post" WHERE id = ANY (ARRAY [448::bigint,449::bigint,507::bigint]);


END
$$;

alter function "ThreadPostsTree"(integer, integer, integer) owner to postgres;

create or replace function "CreatePost"(postsarray type_in_post[]) returns SETOF "Post"
  language plpgsql
as
$$
DECLARE
  in_post               type_in_post;
  current_date_variable timestamptz := now();
  forum_slug            citext;
  thread_identifier     integer;
  posts_setof           "Post"[];
  out_post              "Post";
BEGIN
  FOR i IN 1 .. array_upper(postsArray, 1)
    LOOP
      in_post := postsArray [ i];

      if in_post.thread_slug <> ''::citext then
        select forum, id into forum_slug, thread_identifier from public."Thread" where slug = in_post.thread_slug;
      elsif in_post.thread_id <> 0 then
        thread_identifier := in_post.thread_id;
        select forum into forum_slug from public."Thread" where id = in_post.thread_id;
      else
        raise no_data_found;
        EXIT;
      end if;
      INSERT INTO out_post (author, created, forum, "isEdited", message, parent, thread)
      VALUES (in_post.author::citext, current_date_variable::timestamptz, forum_slug::citext, false, in_post.messag, 0,
              thread_identifier);

      INSERT INTO public."Post" (author, created, forum, "isEdited", message, parent, thread)
      VALUES (in_post.author::citext, current_date_variable::timestamptz, forum_slug::citext, false, in_post.messag, 0,
              thread_identifier);

    END LOOP;
  return QUERY SELECT * FROM public."Post" WHERE created = current_date_variable AND thread = thread_identifier;
END
$$;

alter function "CreatePost"(type_in_post[]) owner to postgres;

create or replace function buildtreefromnode(node_id bigint, thread_id bigint) returns SETOF treenode
  language plpgsql
as
$$
DECLARE
BEGIN
  RETURN QUERY (WITH RECURSIVE tree AS (
    SELECT message,
           id,
           parent,
           NULL::varchar AS parent_name,
           array [id]    AS path
    FROM "Post"
    WHERE id = node_id
      AND "Post".thread = thread_id
    UNION
    SELECT parent_name,
           f1.id,
           f1.parent,
           tree.message       AS parent_name,
           tree.path || f1.id AS path
    FROM tree
           JOIN "Post" f1 ON f1.parent = tree.id
    ) SELECT tree.id, tree.parent,tree.path FROM tree order by path);
end
$$;

alter function buildtreefromnode(bigint, bigint) owner to postgres;

create or replace function "GetThreadPosts"(thread_slug citext, thread_id integer, limitarg integer, sincearg integer,
                                            sortarg text, descarg text) returns SETOF "Post"
  language plpgsql
as
$$
DECLARE
  treadId                            integer := 0;
  thread_counter                     integer := 0;
  limit_posts                        integer := 0;
  since_posts                        integer := 0;
  post_row_temp                      "Post"%ROWTYPE;
  post_row_temp_array                "Post"[];
  post_row_temp_array_for_child_rows "Post"[];
  post_result_array                  "Post"[];
  posts_in_array                     integer := 0;
  tree_result_ids_array              bigint[];
  parent_id_array                    bigint[];
  child_id_array                     bigint[];
  filtered_with_since_array          bigint[];
  tree_since_row                     "Post";
  tree_since_row_id                  bigint;
  node_id_found                      bool;
  MIN_ID                             bigint;
  PARENT_MIN_ID                      bigint;
  parents_counter INTEGER:=0;
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

  if treadId is NULL then
    raise exception raise_exception ;
  end if;
  if limitArg = 0 then
    limit_posts := -1;
  else
    limit_posts := limitArg;
  end if;

  if sinceArg::INTEGER = 0 then
    since_posts := 0;
  else
    since_posts := sinceArg;
  end if;

  if sortArg = 'flat' OR sortArg = '' then
    ------------------------------ FLAT-------------------------
    if descArg = 'true' then
      if since_posts <> 0 then
        RETURN QUERY SELECT *
                     FROM public."Post"
                     WHERE (thread = treadId AND id < since_posts)
                     ORDER BY id DESC
                     LIMIT limit_posts;
      else
        RETURN QUERY SELECT *
                     FROM public."Post"
                     WHERE (thread = treadId)
                     ORDER BY id DESC
                     LIMIT limit_posts;
      end if;


    else
      RETURN QUERY SELECT *
                   FROM public."Post"
                   WHERE (thread = treadId AND id > since_posts)
                   ORDER BY id ASC
                   LIMIT limit_posts;
    end if;
  elsif sortArg = 'tree' then
    if descArg = 'true' then

      if since_posts <> 0 then
        --RAISE 'TREE SORTARG ASC';
        tree_result_ids_array := ARRAY(WITH RECURSIVE tree
          AS
          (
          SELECT message,
                 id,
                 parent,
                 NULL::varchar AS parent_name,
                 array [id]    AS path
          FROM "Post"
          WHERE parent = 0
            AND "Post".thread = treadId
          UNION
          SELECT parent_name,
                 f1.id,
                 f1.parent,
                 tree.message       AS parent_name,
                 tree.path || f1.id AS path
          FROM tree
                 JOIN "Post" f1 ON f1.parent = tree.id
          )
          SELECT tree.id
          FROM tree
          ORDER BY path DESC);

        --RAISE EXCEPTION '% % % % % ',tree_result_ids_array[1], (tree_result_ids_array)[2],(tree_result_ids_array)[3],(tree_result_ids_array)[4],(tree_result_ids_array)[5];
        --https://coderwall.com/p/jmtskw/use-in-instead-of-any-in-postgresql

        -------------------FILTERING PART--------------------------

        FOR i IN 1..array_length(tree_result_ids_array, 1)
          LOOP
            if tree_result_ids_array [ i] = since_posts then
              node_id_found := true;
              CONTINUE;
            end if;

            if array_length(filtered_with_since_array, 1) IS NULL AND node_id_found then
              filtered_with_since_array = array_append(filtered_with_since_array, tree_result_ids_array [ i]);
              CONTINUE;
            end if;

            if array_length(filtered_with_since_array, 1) > 0 AND
               array_length(filtered_with_since_array, 1) < limit_posts and node_id_found then
              filtered_with_since_array = array_append(filtered_with_since_array, tree_result_ids_array [ i]);
            end if;

            if array_length(filtered_with_since_array, 1) = limit_posts and node_id_found then
              EXIT;
            end if;

          end loop;
        ------------------------------------------------------------

        FOR i IN 1..array_length(filtered_with_since_array, 1)
          LOOP
            SELECT * INTO post_row_temp FROM public."Post" WHERE id = filtered_with_since_array [ i];

            post_row_temp_array = array_append(post_row_temp_array, post_row_temp);
          end loop;
        RETURN QUERY SELECT * FROM unnest(post_row_temp_array);
      else
        --RAISE 'TREE SORTARG ASC';
        tree_result_ids_array := ARRAY(WITH RECURSIVE tree
          AS
          (
          SELECT message,
                 id,
                 parent,
                 NULL::varchar AS parent_name,
                 array [id]    AS path
          FROM "Post"
          WHERE parent = 0
            AND "Post".thread = treadId
          UNION
          SELECT parent_name,
                 f1.id,
                 f1.parent,
                 tree.message       AS parent_name,
                 tree.path || f1.id AS path
          FROM tree
                 JOIN "Post" f1 ON f1.parent = tree.id
          )
          SELECT tree.id
          FROM tree
          ORDER BY path DESC);

        --RAISE EXCEPTION '% % % % % ',tree_result_ids_array[1], (tree_result_ids_array)[2],(tree_result_ids_array)[3],(tree_result_ids_array)[4],(tree_result_ids_array)[5];
        --https://coderwall.com/p/jmtskw/use-in-instead-of-any-in-postgresql

        FOR I IN 1..array_length(tree_result_ids_array, 1)
          LOOP
            SELECT * INTO post_row_temp FROM public."Post" WHERE id = tree_result_ids_array [ i];
            post_row_temp_array = array_append(post_row_temp_array, post_row_temp);
          end loop;
        RETURN QUERY SELECT * FROM unnest(post_row_temp_array) LIMIT limit_posts;
      end if;

      --------------------------------------------------------------------------------------
      -------------------------------------------------------------------------------------
    else
      if since_posts <> 0 then
        tree_result_ids_array := ARRAY(WITH RECURSIVE tree
          AS
          (
          SELECT message,
                 id,
                 parent,
                 NULL::varchar AS parent_name,
                 array [id]    AS path
          FROM "Post"
          WHERE parent = 0
            AND "Post".thread = treadId
          UNION
          SELECT parent_name,
                 f1.id,
                 f1.parent,
                 tree.message       AS parent_name,
                 tree.path || f1.id AS path
          FROM tree
                 JOIN "Post" f1 ON f1.parent = tree.id
          )
          SELECT tree.id
          FROM tree
          ORDER BY path);

        FOR i IN 1..array_length(tree_result_ids_array, 1)
          LOOP
            if tree_result_ids_array [ i] = since_posts then
              node_id_found := true;
              CONTINUE;
            end if;

            if array_length(filtered_with_since_array, 1) IS NULL AND node_id_found then
              filtered_with_since_array = array_append(filtered_with_since_array, tree_result_ids_array [ i]);
              CONTINUE;
            end if;

            if array_length(filtered_with_since_array, 1) > 0 AND
               array_length(filtered_with_since_array, 1) < limit_posts and node_id_found then
              filtered_with_since_array = array_append(filtered_with_since_array, tree_result_ids_array [ i]);
            end if;

            if array_length(filtered_with_since_array, 1) = limit_posts and node_id_found then
              EXIT;
            end if;

          end loop;


        FOR i IN 1..array_length(filtered_with_since_array, 1)
          LOOP
            SELECT * INTO post_row_temp FROM public."Post" WHERE id = filtered_with_since_array [ i];

            post_row_temp_array = array_append(post_row_temp_array, post_row_temp);
          end loop;


        RETURN QUERY SELECT * FROM unnest(post_row_temp_array);

      else
        tree_result_ids_array := ARRAY(WITH RECURSIVE tree
          AS
          (
          SELECT message,
                 id,
                 parent,
                 NULL::varchar AS parent_name,
                 array [id]    AS path
          FROM "Post"
          WHERE parent = 0
            AND "Post".thread = treadId
          UNION
          SELECT parent_name,
                 f1.id,
                 f1.parent,
                 tree.message       AS parent_name,
                 tree.path || f1.id AS path
          FROM tree
                 JOIN "Post" f1 ON f1.parent = tree.id
          )
          SELECT tree.id
          FROM tree
          ORDER BY path);

        FOR I IN 1..array_length(tree_result_ids_array, 1)
          LOOP
            SELECT * INTO post_row_temp FROM public."Post" WHERE id = tree_result_ids_array [ i];
            post_row_temp_array = array_append(post_row_temp_array, post_row_temp);
          end loop;

        RETURN QUERY SELECT * FROM unnest(post_row_temp_array) LIMIT limit_posts;

      end if;

    end if;

    ---------------------------         PARENT TREE      -------------------------------------------------------
  elsif sortArg = 'parent_tree' then
    if descArg = 'true' then
      --------------------------------DESC=TRUE --------------------------------------------------------------
      if since_posts <> 0 then
        --------------------------------SINCE !=0 DESC --------------------------------------------------------------
        parent_id_array :=
            ARRAY(SELECT id
                  FROM public."Post"
                  WHERE parent = 0
                    AND "Post".thread = treadId
                    AND id < since_posts
                  ORDER BY id ASC
                  LIMIT 1);
        if array_length(parent_id_array, 1) = 0 then
          RETURN QUERY SELECT * FROM unnest(parent_id_array);
          return;
        end if;

        if array_length(parent_id_array, 1) = 1 then
          PARENT_MIN_ID := (SELECT parent FROM public."Post" WHERE id = since_posts);
          if parent_id_array [ 1] = PARENT_MIN_ID then
            return;
          end if;
        end if;

        --https://coderwall.com/p/jmtskw/use-in-instead-of-any-in-postgresql
        FOR i IN 1..array_length(parent_id_array, 1)
          LOOP

            child_id_array := ARRAY(WITH RECURSIVE tree
              AS
              (
              (SELECT message,
                      id,
                      parent,
                      NULL::varchar AS parent_name,
                      array [id]    AS path
               FROM "Post"
               WHERE parent = 0
                 AND "Post".thread = treadId
                 AND id = parent_id_array [ i]
              )
              UNION
              SELECT parent_name,
                     f1.id,
                     f1.parent,
                     tree.message       AS parent_name,
                     tree.path || f1.id AS path
              FROM tree
                     JOIN "Post" f1 ON f1.parent = tree.id
              )
              SELECT tree.id
              FROM tree
              ORDER BY path);
            tree_result_ids_array = array_cat(tree_result_ids_array, child_id_array);
          end loop;

        FOR I IN 1..array_length(tree_result_ids_array, 1)
          LOOP
            SELECT * INTO post_row_temp FROM public."Post" WHERE id = tree_result_ids_array [ i];
            post_row_temp_array = array_append(post_row_temp_array, post_row_temp);
          end loop;

        RETURN QUERY SELECT * FROM unnest(post_row_temp_array);
        return;
      else
        --------------------------------SINCE = 0 DESC --------------------------------------------------------------
        parent_id_array :=
            ARRAY(SELECT id
                  FROM public."Post"
                  WHERE parent = 0
                    AND "Post".thread = treadId
                  ORDER BY id DESC
                  LIMIT limit_posts);
        if array_length(parent_id_array, 1) = 0 then
          RETURN QUERY SELECT * FROM unnest(parent_id_array);
          return;
        end if;
        --https://coderwall.com/p/jmtskw/use-in-instead-of-any-in-postgresql
        FOR i IN 1..array_length(parent_id_array, 1)
          LOOP
            child_id_array := ARRAY(WITH RECURSIVE tree
              AS
              (
              (SELECT message,
                      id,
                      parent,
                      NULL::varchar AS parent_name,
                      array [id]    AS path
               FROM "Post"
               WHERE parent = 0
                 AND "Post".thread = treadId
                 AND id = parent_id_array [ i]
              )
              UNION
              SELECT parent_name,
                     f1.id,
                     f1.parent,
                     tree.message       AS parent_name,
                     tree.path || f1.id AS path
              FROM tree
                     JOIN "Post" f1 ON f1.parent = tree.id
              )
              SELECT tree.id
              FROM tree
              ORDER BY path);
            tree_result_ids_array = array_cat(tree_result_ids_array, child_id_array);
          end loop;

        FOR I IN 1..array_length(tree_result_ids_array, 1)
          LOOP
            SELECT * INTO post_row_temp FROM public."Post" WHERE id = tree_result_ids_array [ i];
            post_row_temp_array = array_append(post_row_temp_array, post_row_temp);
          end loop;

        RETURN QUERY SELECT * FROM unnest(post_row_temp_array);
        return;
      end if;
    else
      --------------------------------DESC=FALSE --------------------------------------------------------------
      if since_posts <> 0 then
        tree_result_ids_array := ARRAY(WITH RECURSIVE tree
          AS
          (
          (SELECT message,
                  id,
                  parent,
                  NULL::varchar AS parent_name,
                  array [id]    AS path
           FROM "Post"
           WHERE parent = 0
             AND "Post".thread = treadId
          )
          UNION
          SELECT parent_name,
                 f1.id,
                 f1.parent,
                 tree.message       AS parent_name,
                 tree.path || f1.id AS path
          FROM tree
                 JOIN "Post" f1 ON f1.parent = tree.id
          )
          SELECT tree.id
          FROM tree
          ORDER BY path);

        if array_length(tree_result_ids_array, 1) IS NULL or array_length(tree_result_ids_array, 1) = 0 then
          RAISE NOTICE 'empty tree_result_ids_array ARRAY';
          RETURN;
        end if;

        node_id_found := false;
        FOR i IN 1..array_length(tree_result_ids_array, 1)
          LOOP
            if tree_result_ids_array [ i] = since_posts then
              node_id_found := true;
              CONTINUE;
            end if;

            IF node_id_found then
              SELECT * INTO post_row_temp FROM public."Post" where id = tree_result_ids_array [ i];
              if post_row_temp.parent = 0 then
                if parents_counter < limit_posts then
                  --append parent
                  filtered_with_since_array := array_append(filtered_with_since_array, tree_result_ids_array [ i]);
                  parents_counter := parents_counter + 1;
                else
                  EXIT;
                end if;
              else
                filtered_with_since_array := array_append(filtered_with_since_array, tree_result_ids_array [ i]);
              end if;
            end if;
          end loop;

        if array_length(filtered_with_since_array, 1) IS NULL or array_length(filtered_with_since_array, 1) = 0 then
          RAISE NOTICE 'empty filtered_with_since_array ARRAY';
          RETURN;
        end if;

        FOR i in 1..array_length(filtered_with_since_array, 1)
          LOOP
            SELECT * INTO post_row_temp FROM public."Post" WHERE id = filtered_with_since_array [ i];
            RAISE NOTICE '% %', post_row_temp.id, post_row_temp.parent;
            post_row_temp_array:=array_append(post_row_temp_array, post_row_temp);
          end loop;

        RETURN QUERY SELECT * FROM unnest(post_row_temp_array);
        return;



      else
        tree_result_ids_array := ARRAY(WITH RECURSIVE tree   AS   (
          (SELECT message,
                  id,
                  parent,
                  NULL::varchar AS parent_name,
                  array [id]    AS path
           FROM "Post"
           WHERE parent = 0
             AND "Post".thread = treadId
           LIMIT limit_posts)
          UNION
          SELECT parent_name,
                 f1.id,
                 f1.parent,
                 tree.message       AS parent_name,
                 tree.path || f1.id AS path
          FROM tree
                 JOIN "Post" f1 ON f1.parent = tree.id
          )
          SELECT tree.id
          FROM tree
          ORDER BY path);

        FOR i IN 1..array_length(tree_result_ids_array, 1)
          LOOP
            SELECT * INTO post_row_temp FROM public."Post" WHERE id = tree_result_ids_array [ i];
            post_row_temp_array = array_append(post_row_temp_array, post_row_temp);
          end loop;
      end if;

      RETURN QUERY SELECT * FROM unnest(post_row_temp_array);

    end if;
  end if;
END
$$;

alter function "GetThreadPosts"(citext, integer, integer, integer, text, text) owner to postgres;




create or replace function "UpdateThreadDetails"(thread_id integer, thread_slug text, new_title text, new_message text) returns SETOF "Thread"
  language plpgsql
as
$$
DECLARE
  thread_counter integer;
  treadId        BIGINT;
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

  if thread_counter = 0 then
    RAISE EXCEPTION 'THREAD NOT FOUND  slug  %, id %', thread_slug::text,thread_id::text USING ERRCODE = 'no_data_found';
  end if;

  if "new_message" <> '' then
    UPDATE "Thread" SET message = COALESCE("new_message"::citext, message) WHERE id = treadId;
  end if;

  if "new_title" <> '' then
    UPDATE "Thread" SET title = COALESCE("new_title"::citext, title) WHERE id = treadId;
  end if;

  RETURN QUERY SELECT * FROM public."Thread" Where id = treadId;

END
$$;

alter function "UpdateThreadDetails"(integer, text, text, text) owner to postgres;

create or replace function "GetForumUsers"(forum_slug citext, limitarg integer, sincearg citext, descarg text) returns SETOF "User"
  language plpgsql
as
$$
DECLARE
  temp_counter INTEGER := 0;
  limit_Arg    INTEGER := 0;
BEGIN
  temp_counter := (SELECT count(*) FROM public."Forum" WHERE slug = forum_slug::citext);
  if temp_counter = 0 then
    RAISE EXCEPTION no_data_found;
  end if;

  if limitarg = 0 then
    limit_Arg := 100000000;
  else
    limit_Arg := limitarg;
  end if;

  if descArg = 'true' then
    if sincearg = ''
    then
      RETURN QUERY SELECT about,email,fullname,nickname
                   FROM forumuser
                          INNER JOIN
                        public."User" ON "User".nickname = forumuser.user
                   WHERE forum = forum_slug::citext
                   ORDER BY lower(nickname) DESC
                   LIMIT limit_Arg;

    else
      RETURN QUERY SELECT about,email,fullname,nickname
                   FROM forumuser
                          INNER JOIN
                        public."User" ON "User".nickname = forumuser.user
                   WHERE forum = forum_slug::citext
                     AND nickname < sinceArg
                   ORDER BY lower(nickname) DESC
                   LIMIT limit_Arg;

    end if;
  else
    if sincearg = ''
    then
      RETURN QUERY SELECT about,email,fullname,nickname
                   FROM forumuser
                          INNER JOIN
                        public."User" ON "User".nickname = forumuser.user
                   WHERE forum = forum_slug::citext
                   ORDER BY lower(nickname) ASC
                   LIMIT limit_Arg;

    else
      RETURN QUERY SELECT about,email,fullname,nickname
                   FROM forumuser
                          INNER JOIN
                        public."User" ON "User".nickname = forumuser.user
                   WHERE forum = forum_slug::citext
                     AND nickname > sinceArg
                   ORDER BY lower(nickname) ASC
                   LIMIT limit_Arg;

    end if;
  end if;


END
$$;

alter function "GetForumUsers"(citext, integer, citext, text) owner to postgres;

create or replace function "CreateOrGetThread"("postId" integer, "newMessage" text) returns SETOF "Post"
  language plpgsql
as
$$
DECLARE
  temp_counter integer := 0;
  temp_message text;
BEGIN
  temp_counter := (SELECT count(*) FROM public."Post" where id = "postId");
  if temp_counter = 0 then
    RAISE EXCEPTION no_data_found;
  end if;
  temp_message := (SELECT "message" FROM public."Post" WHERE id = "postId");
  if temp_message <> "newMessage" then
    UPDATE public."Post" SET message = "newMessage", "isEdited" = true WHERE id = "postId";
  end if;

  RETURN QUERY SELECT * FROM public."Post" WHERE id = "postId";
END
$$;

alter function "CreateOrGetThread"(integer, text) owner to postgres;

create or replace function "CreatePostUsingFieldArrays"(author_array citext[], message_array text[],
                                                        parent_array integer[], array_len integer, thread_slug citext,
                                                        thread_id integer, parent_is_the_same boolean,
                                                        author_is_the_same boolean) returns SETOF "Post"
  language plpgsql
as
$$
DECLARE
  thread_counter  integer := 0;
  thread_row      public."Thread"%ROWTYPE;
  forum_counter   integer := 0;
  parent_counter  integer := 0;
  author_counter  integer := 0;
  common_parent   integer := 0;
  common_author   text    := '';
  current_date_my timestamptz;
  temp_counter    INTEGER;
BEGIN

  -- CHECK IF THREAD BY SLUG EXISTS
  current_date_my := now();
  if thread_id = 0 then
    SELECT count(*) INTO thread_counter from public."Thread" where slug = thread_slug::citext;
    SELECT * INTO thread_row from public."Thread" where slug = thread_slug::citext;
  end if;

  -- CHECK IF THREAD BY ID EXISTS
  if thread_slug = '' then
    SELECT count(*) INTO thread_counter from public."Thread" where id = thread_id;
    SELECT * INTO thread_row from public."Thread" where id = thread_id;
  end if;

  if thread_counter = 0 OR thread_row IS NULL then
    RAISE EXCEPTION no_data_found;
  end if;

  if array_length(message_array, 1) IS NULL
  then
    RETURN;
  end if;
  -- CHECK IF FORUM BY SLUG EXISTS
  SELECT count(*) INTO forum_counter from public."Forum" where slug = thread_row.forum::citext;
  if forum_counter = 0 then
    RAISE EXCEPTION no_data_found;
  end if;

  if author_is_the_same = true AND parent_is_the_same = true then
    common_parent = parent_array [ 1];
    common_author = author_array [ 1];
    if common_parent <> 0 then
      temp_counter := (SELECT thread FROM "Post" WHERE common_parent = id);
      if temp_counter <> thread_row.id OR temp_counter IS NULL then
        raise exception raise_exception;
      end if;

    end if;
    INSERT INTO public."Post" ("author", "created", "forum", "isEdited", "message", "parent", "thread")
    VALUES (common_author, current_date_my, thread_row.forum, false, unnest(message_array), common_parent,
            thread_row.id);
    UPDATE public."Forum"
    SET posts = posts + array_length(message_array, 1)
    where "Forum".slug = thread_row.forum::citext;

    temp_counter := (SELECT count(*)
                     FROM public.forumuser
                     WHERE forum = thread_row.forum::citext
                       AND "user" = common_author::citext);
    if temp_counter = 0 then
      INSERT INTO public.forumuser ("forum", "user") VALUES (thread_row.forum::citext, common_author::citext);
    end if;

    RETURN QUERY SELECT * FROM public."Post" where created = current_date_my;
    return;
  else
    -- IF authors are the same then check it's existance if it not 0, else ok
    if author_is_the_same = true then
      common_author = author_array [ 1];
      if common_author <> '' then
        SELECT count(*) INTO author_counter from public."User" where nickname = common_author::citext;
        if author_counter = 0 then
          RAISE EXCEPTION 'COMMON AUTHOR NOT FOUND  nickname %', common_author USING ERRCODE = 'no_data_found';
        end if;
      end if;
    else
      FOR idx IN 1..array_len
        LOOP
          SELECT count(*) INTO author_counter from public."User" where nickname = author_array [ idx]::citext;
          if author_counter = 0 then
            RAISE EXCEPTION 'VARY AUTHOR NOT FOUND  nickname %', author_array [ idx] USING ERRCODE = 'no_data_found';
          end if;
        END LOOP;
    end if;

    -- IF parents are the same then check it's existance if it not 0, else ok
    if parent_is_the_same = true then
      common_parent = parent_array [ 1];
      if common_parent <> 0 then
        SELECT count(*) INTO parent_counter from public."Post" where id = common_parent::integer;
        if parent_counter = 0 then
          RAISE EXCEPTION 'COMMON PARENT NOT FOUND  common_parent %', common_parent USING ERRCODE = 'no_data_found';
        end if;
      end if;
    else
      FOR idx IN 1..array_len
        LOOP
          if parent_array [ idx] <> 0 then
            temp_counter := (SELECT thread from public."Post" where id = parent_array [ idx]::BIGINT);
            if temp_counter <> thread_row.id OR temp_counter IS NULL then
              RAISE EXCEPTION raise_exception;
            end if;
          end if;
        END LOOP;

      FOR idx IN 1..array_len
        LOOP
          SELECT count(*) INTO parent_counter from public."Post" where id = parent_array [ idx]::BIGINT;
          if parent_counter = 0 AND parent_array [ idx] <> 0 then
            RAISE EXCEPTION 'VARY PARENT NOT FOUND  parent %', parent_array [ idx] USING ERRCODE = 'no_data_found';
          end if;
        END LOOP;
    end if;

    INSERT INTO public."Post" ("author", "created", "forum", "isEdited", "message", "parent", "thread")
    VALUES (unnest(author_array), current_date_my, thread_row.forum, false, unnest(message_array), unnest(parent_array),
            thread_row.id);
    FOR i in 1..array_length(author_array, 1)
      LOOP
        temp_counter := (SELECT count(*)
                         FROM public.forumuser
                         WHERE forum = thread_row.forum::citext
                           AND "user" = author_array [ i]::citext);
        if temp_counter = 0 then
          INSERT INTO public.forumuser ("forum", "user") VALUES (thread_row.forum::citext, author_array [ i]::citext);
        end if;
      end loop;

    UPDATE public."Forum"
    SET posts = posts + array_length(author_array, 1)
    where "Forum".slug = thread_row.forum::citext;

    RETURN QUERY SELECT * FROM public."Post" where created = current_date_my ORDER BY id ASC;
    return;
  end if;
END
$$;

alter function "CreatePostUsingFieldArrays"(citext[], text[], integer[], integer, citext, integer, boolean, boolean) owner to postgres;

create or replace function clearalldata() returns SETOF "Status"
  language plpgsql
as
$$
DECLARE
BEGIN
  DELETE FROM public."Error" CASCADE;
  DELETE FROM public."Vote" CASCADE;
  DELETE FROM public."forumuser" CASCADE;
  DELETE FROM public."Post" CASCADE;
  DELETE FROM public."Thread" CASCADE;

  DELETE FROM public."Forum" CASCADE;
  DELETE FROM public."Status" CASCADE;
  DELETE FROM public."User" CASCADE;

  RETURN QUERY SELECT * FROM public."GetServiceStatus"();
end
$$;

alter function clearalldata() owner to postgres;

create or replace function "GetServiceStatus"() returns "Status"
  language plpgsql
as
$$
DECLARE
  status_output ServiceStatus;
BEGIN

  status_output.forum := COALESCE((SELECT count(*) from public."Forum"), 0::INTEGER);
  status_output.post := COALESCE((SELECT count(*) from public."Post"), 0::INTEGER);
  status_output.thread := COALESCE((SELECT count(*) from public."Thread"), 0::INTEGER);
  status_output."user" := COALESCE((SELECT count(*) from public."User"), 0::INTEGER);

  return "status_output";

END
$$;

create sequence posts_sequence
  as integer
  maxvalue 2147483647;

alter sequence posts_sequence owner to postgres;


create sequence thread_id_serial
  as integer
  maxvalue 2147483647;

alter sequence thread_id_serial owner to postgres;

