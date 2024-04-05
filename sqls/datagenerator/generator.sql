
CREATE OR REPLACE FUNCTION genre_generate() RETURNS VOID AS $genre_generate$
DECLARE
    --Задаем количество записей в таблице жанров
    number integer := 10;
BEGIN
    INSERT INTO genres_table (id, shikimori_id, genre_name, russian)
    SELECT
        md5(random()::text)::char(32) AS id,
        s AS shikimori_id,
        md5(random()::text)::char(10) AS genre_name,
        md5(random()::text)::char(10) AS russian
    FROM
        generate_series(1, number) s;

END;
$genre_generate$ LANGUAGE plpgsql;

SELECT genre_generate();

-- DO $$ BEGIN
--     INSERT INTO genres_table (id, shikimori_id, genre_name, russian)
--     SELECT
--         md5(random()::text)::char(32) AS id,
--         s AS shikimori_id,
--         md5(random()::text)::char(10) AS genre_name,
--         md5(random()::text)::char(10) AS russian
--     FROM
--         generate_series(1, 10) s;
-- END $$;

SELECT * FROM genres_table;


CREATE OR REPLACE FUNCTION image_generate() RETURNS VOID AS $image_generate$
DECLARE
    --Задаем количество записей в таблице image_table
    number integer := 10;
    hash_image varchar(64) := NULL;
BEGIN
    INSERT INTO images_table (hash, source_img)
    SELECT
        md5(random()::text) AS hash,
        md5(random()::text)::char(10) AS source_img
    FROM
        generate_series(1, number);
    FOR hash_image IN SELECT hash FROM images_table LOOP
      INSERT INTO ipfs_object (hash) VALUES (hash_image);
    END LOOP;
END;
$image_generate$ LANGUAGE plpgsql;

SELECT image_generate();

SELECT * FROM images_table;
SELECT * FROM ipfs_object;

-- CREATE TABLE IF NOT EXISTS studio_table(
--     id uuid PRIMARY KEY, --id студии
--     shikimori_id integer, --id с сайта shikimori
--     studio_name varchar(100), --название студии
--     image varchar(64) REFERENCES images_table (hash) --url логотипа студии
-- );


CREATE OR REPLACE FUNCTION stidio_generate() RETURNS VOID AS $stidio_generate$
DECLARE
    --Задаем количество записей в таблице image_table
    number integer := 10;
    hash_image varchar(64) := NULL;
BEGIN
    INSERT INTO studio_table (id, shikimori_id, studio_name)
    SELECT
        gen_random_uuid() AS hash,
        s AS shikimori_id,
        md5(random()::text)::char(10) AS studio_name
    FROM
        generate_series(1, number) s;
--     FOR hash_image IN SELECT hash FROM images_table LOOP
--         INSERT INTO studio_table (image) VALUES (hash_image);
--     END LOOP;
END;
$stidio_generate$ LANGUAGE plpgsql;

SELECT stidio_generate();

SELECT * FROM studio_table;


-- CREATE TABLE IF NOT EXISTS video_table(
--     hash varchar(64) PRIMARY KEY, --id видео IPFS CID
--     shikimori_id integer, --id с сайта shikimori
--     url varchar(255), --ссылка на youtube видео
--     player_url varchar(255),  --ссылка на youtube видео на весь экран
--     video_name varchar(100), --название эпизода
--     video_kind kind_video --тип видео
-- );

-- CREATE TYPE kind_video AS ENUM ('pv', 'op', 'cm', 'ed', 'mv');

CREATE OR REPLACE FUNCTION video_generate() RETURNS VOID AS $video_generate$
DECLARE
    --Задаем количество записей в таблице video_table
    number integer := 10;
    hash_video varchar(64) := NULL;
BEGIN
    INSERT INTO video_table (hash, shikimori_id, url, player_url, video_name, video_kind)
    SELECT
        md5(random()::text) AS hash,
        s AS shikimori_id,
        md5(random()::text)::char(20) AS url,
        md5(random()::text)::char(20) AS player_url,
        md5(random()::text)::char(20) AS video_name,
        (array['pv', 'op', 'cm', 'ed', 'mv'])[floor(random() * 5 + 1)]::text::kind_video AS video_kind
    FROM
        generate_series(1, number) s;
    FOR hash_video IN SELECT hash FROM video_table LOOP
        INSERT INTO ipfs_object (hash) VALUES (hash_video);
    END LOOP;
END;
$video_generate$ LANGUAGE plpgsql;

SELECT video_generate();

SELECT * FROM video_table;
SELECT * FROM ipfs_object;

