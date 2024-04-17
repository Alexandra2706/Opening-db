-- Функция генерации таблицы жанров
CREATE OR REPLACE FUNCTION genre_generate() RETURNS VOID AS $genre_generate$
DECLARE
    --Задаем количество записей в таблице жанров
    number integer := 30;
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
SELECT * FROM genres_table;

-- Функция генерирует запись в таблице изображений images_table
-- и добавляет запись в таблицу ipfs_object
-- возвращает hash
CREATE OR REPLACE FUNCTION image_generate() RETURNS VARCHAR(64) AS $image_generate$
DECLARE
    current_hash varchar(64) := NULL;
BEGIN
    current_hash := md5(random()::text);
    INSERT INTO images_table (hash, source_img)
    SELECT
        current_hash AS hash,
        md5(random()::text)::char(10) AS source_img
    FROM
        generate_series(1, 1);
    INSERT INTO ipfs_object (hash, mime_type) VALUES (current_hash, 'image/jpeg');
    RETURN current_hash;
END;
$image_generate$ LANGUAGE plpgsql;


-- Функция генерирует таблицу studio_table
CREATE OR REPLACE FUNCTION studio_generate() RETURNS VOID AS $studio_generate$
DECLARE
    --Задаем количество записей в таблице image_table
    number integer := 100;
BEGIN
    INSERT INTO studio_table (id, shikimori_id, studio_name, image)
    SELECT
        gen_random_uuid() AS hash,
        s AS shikimori_id,
        md5(random()::text)::char(10) AS studio_name,
        image_generate() AS image
    FROM
        generate_series(1, number) s;
END;
$studio_generate$ LANGUAGE plpgsql;

SELECT studio_generate();
SELECT * FROM studio_table;
SELECT * FROM images_table;


-- Функция генерирует таблицу video_table
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
        INSERT INTO ipfs_object (hash, mime_type) VALUES (hash_video, 'video/mp4');
    END LOOP;
END;
$video_generate$ LANGUAGE plpgsql;

SELECT video_generate();

SELECT * FROM video_table;
SELECT * FROM ipfs_object;


-- Функция генерирует таблицу anime
CREATE OR REPLACE FUNCTION animes_generate() RETURNS VOID AS $animes_generate$
DECLARE
    -- Задаем количество записей в таблице anime
    number integer := 3;
    number_episodes integer := 0; --количество серий
    data_issue timestamp with time zone;
    i integer;
    total_genres integer := 0; --количество записей в таблице genres_table
    genre_number integer := 0; --количество выбранных записей из таблицы genres_table
    array_genres varchar(32)[] := NULL; --массив id выбранных записей из таблицы genres_table

    row_number_studios integer := 0; --количество записей в таблице studio_table
    studio_number integer :=0; --количество выбранных записей из таблицы studio_table
    array_studios uuid[] := NULL; --массив id выбранных записей из таблицы studio_table

    row_number_videos integer := 0; --количество записей в таблице video_table
    video_number integer :=0; --количество выбранных записей из таблицы video_table
    array_videos varchar(64)[] := NULL; --массив hash выбранных записей из таблицы video_table

    screenshots_number integer :=0;
    max_screenshots_number integer :=5;
    screenshots_array varchar(64)[];

BEGIN
    data_issue := '2000-01-01 00:00:00'::timestamp +
                  ('2023-12-31 23:59:59'::timestamp - '2000-01-01 00:00:00'::timestamp) * RANDOM();
    SELECT COUNT(*) INTO total_genres FROM genres_table; --считает кол-во строк в таблице genres_table
    SELECT COUNT(*) INTO row_number_studios FROM studio_table; --считает кол-во строк в таблице studio_table
    SELECT COUNT(*) INTO row_number_videos FROM video_table; --считает кол-во строк в таблице video_table

    FOR i IN 1..5 LOOP
        number_episodes :=  ROUND(RANDOM() * 100);

        genre_number := ROUND(RANDOM() * total_genres); --количество выбранных записаей из таблицы genres_table
        IF genre_number != 0 THEN
            SELECT ARRAY(SELECT id INTO array_genres FROM genres_table ORDER BY RANDOM() LIMIT genre_number);
            --RAISE NOTICE 'array_genres = %', array_genres;
        END IF;

        studio_number := ROUND(RANDOM()*row_number_studios); --количество выбранных записаей из таблицы studio_table
        IF studio_number != 0 THEN
            SELECT ARRAY(SELECT id INTO array_studios FROM studio_table ORDER BY RANDOM() LIMIT studio_number);
        END IF;

        video_number := ROUND(RANDOM()*row_number_studios); --количество выбранных записаей из таблицы video_table
        IF video_number != 0 THEN
            SELECT ARRAY(SELECT hash INTO array_videos FROM video_table ORDER BY RANDOM() LIMIT video_number);
        END IF;

        screenshots_number := ROUND(RANDOM()*max_screenshots_number); --количество скриншотов для таблицы animes
        --RAISE NOTICE 'screenshots_number = %', screenshots_number;
        IF screenshots_number != 0 THEN
            SELECT ARRAY(SELECT image_generate() INTO screenshots_array FROM generate_series(1, screenshots_number));
        END IF;

        INSERT INTO animes (id,anime_name, name_russian, name_english, name_japanese, name_synonyms,
                            anime_status, episodes, episodes_aired, aired_on, released_on, duration,
                            updated_at, next_episode_at, image, genres, studios, videos, screenshots,
                            shikimori_id, shikimori_kind, shikimori_rating, shikimori_description,
                            shikimori_description_html,shikimori_last_revision, myanimelist_id, myanimelist_score)
        VALUES (
        gen_random_uuid(),
        md5(random()::text)::char(10),
        md5(random()::text)::char(10),
        ARRAY[md5(random()::text)::char(10)],
        ARRAY[md5(random()::text)::char(10)],
        ARRAY[md5(random()::text)::char(10)],
        (ARRAY['anons', 'ongoing', 'released'])[floor(random() * 3 + 1)]::status,
        number_episodes,
        ROUND(RANDOM() * number_episodes),
        data_issue,
        data_issue + ('2023-12-31 23:59:59'::timestamp - data_issue) * RANDOM(),
        ROUND(RANDOM() * 199 + 1),
        data_issue + ('2023-12-31 23:59:59'::timestamp - data_issue) * RANDOM(),
        md5(random()::text)::char(10),
        image_generate(),
        array_genres,
        array_studios,
        array_videos,
        screenshots_array,
        i,
        (ARRAY['tv', 'movie', 'ova', 'ona', 'other'])[floor(random() * 4 + 1)]::kind,
        (ARRAY['r_plus', 'pg_13', 'r', 'g', 'rx'])[floor(random() * 5)]::rating,
        md5(random()::text),
        md5(random()::text),
        '2000-01-01 00:00:00'::timestamp + ('2024-04-9 23:59:59'::timestamp - '2000-01-01 00:00:00') * RANDOM(),
        i,
        RANDOM()*5
        );
    END LOOP;
END;
$animes_generate$ LANGUAGE plpgsql;

SELECT animes_generate();

SELECT * FROM animes;
SELECT * FROM ipfs_object;