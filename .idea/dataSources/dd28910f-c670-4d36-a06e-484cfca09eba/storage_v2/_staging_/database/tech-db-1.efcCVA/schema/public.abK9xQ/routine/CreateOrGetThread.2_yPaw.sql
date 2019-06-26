create function "UpdatePostDetails"("postId" integer, "newMessage" text) returns SETOF "Post"
  language plpgsql
as
$$
DECLARE
  temp_counter integer :=0;
  temp_message text;
BEGIN
 temp_counter:=(SELECT count(*) FROM public."Post" where id="postId");
 if temp_counter=0 then
   RAISE EXCEPTION no_data_found;
 end if;
  temp_message:=(SELECT "message" FROM public."Post" WHERE id="postId");
  if temp_message  <> "newMessage" AND "newMessage" IS NOT NULL AND "newMessage" <> '' then
    UPDATE public."Post" SET "message" = COALESCE("newMessage","message"), "isEdited" = true WHERE id = "postId";
  end if;

  RETURN QUERY SELECT * FROM public."Post"  WHERE id = "postId";
END
$$;

alter function "UpdatePostDetails"(integer, text) owner to postgres;

