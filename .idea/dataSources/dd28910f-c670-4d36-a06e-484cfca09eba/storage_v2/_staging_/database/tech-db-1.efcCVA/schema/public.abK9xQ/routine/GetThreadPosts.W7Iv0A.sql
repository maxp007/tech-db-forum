create function "GetThreadPosts"(thread_slug citext, thread_id integer, limitarg integer, sincearg integer,
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
          RAISE EXCEPTION 'empty tree_result_ids_array ARRAY'  USING ERRCODE = 'P0001';
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
          RAISE NOTICE 'empty filtered_with_since_array ARRAY' USING ERRCODE = 'P0001';
          ;
          RETURN;
        end if;

        FOR i in 1..array_length(filtered_with_since_array, 1)
          LOOP
            SELECT * INTO post_row_temp FROM public."Post" WHERE id = filtered_with_since_array [ i];
            post_row_temp_array:=array_append(post_row_temp_array, post_row_temp);
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
             AND "Post".thread = 33693
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
            if tree_result_ids_array [ i] = 331856 then
              node_id_found := true;
              CONTINUE;
            end if;

            IF node_id_found then
              SELECT * INTO post_row_temp FROM public."Post" where id = tree_result_ids_array [ i];
              if post_row_temp.parent = 0 then
                if parents_counter < 3 then
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

