-- Функция генерации таблицы жанров
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

--SELECT image_generate();
--SELECT * FROM images_table;
--SELECT * FROM ipfs_object;

-- Функция генерирует таблицу studio_table
CREATE OR REPLACE FUNCTION stidio_generate() RETURNS VOID AS $stidio_generate$
DECLARE
    --Задаем количество записей в таблице image_table
    number integer := 10;
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
$stidio_generate$ LANGUAGE plpgsql;

SELECT stidio_generate();
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


--
-- CREATE TABLE IF NOT EXISTS screenshots_table(
--     hash varchar(64) PRIMARY KEY, --id скриншота IPFS CID
--     original varchar(255), --url исходного изображения
--     preview varchar(255) --url превью
-- );

-- Функция генерирует запись в таблице скриншотов screenshots_table
-- и добавляет запись в таблицу ipfs_object
-- возвращает hash
CREATE OR REPLACE FUNCTION screenshot_generate() RETURNS VARCHAR(64) AS $screenshot_generate$
DECLARE
    current_hash varchar(64) := NULL;
BEGIN
    current_hash := md5(random()::text);
    INSERT INTO screenshots_table (hash, original, preview)
    SELECT
        current_hash AS hash,
        md5(random()::text)::char(10) AS original,
        md5(random()::text)::char(10) AS preview
    FROM
        generate_series(1, 1);
    INSERT INTO ipfs_object (hash, mime_type) VALUES (current_hash, 'image/jpeg');
    RETURN current_hash;
END;
$screenshot_generate$ LANGUAGE plpgsql;

SELECT screenshot_generate();

SELECT * FROM screenshots_table;
SELECT * FROM ipfs_object;

-- CREATE TABLE IF NOT EXISTS animes(
--     -- Основные поля:
--     id uuid PRIMARY KEY, --наш уникальный id
--     anime_name varchar(255) UNIQUE NOT NULL, --название анимэ
--     name_russian varchar(255), --название анимэ на русском
--     name_english varchar[], --название анимэ на английском
--     name_japanese varchar[], --название анимэ на японском
--     name_synonyms varchar[], --синонимы названия анимэ
--     anime_status status NOT NULL, --статус: anons, ongoing, released
--     episodes integer DEFAULT 0, --количество серий
--     episodes_aired integer DEFAULT 0, --количество вышедших эпизодов
--     aired_on timestamp with time zone, --начало выпуска, формат ISO 8601 with TimeZone
--     released_on timestamp with time zone, --конец выпуска, формат ISO 8601 with TimeZone
--     duration integer, --длительность серии в минутах
--     licensors_ru jsonb, --лицензировано
--     franchise jsonb, --франшиза
--     updated_at timestamp with time zone DEFAULT NOW(), --дата обновления, формат ISO 8601 with TimeZone
--     next_episode_at varchar(255), --следующая серия ссылка
--     image varchar(64) REFERENCES images_table (hash), --постер аниме (изображения на сайте shikimori)
--     genres varchar(32)[], --жанры, может быть несколько
--     studios uuid[], --REFERENCES Studio (id), --студии, может быть несколько
--     videos varchar(64)[], --REFERENCES Video (id), --эпизоды
--     screenshots varchar(64)[] -- REFERENCES Sreenshot (id), --кадры
-- );

-- Функция генерирует таблицу anime
CREATE OR REPLACE FUNCTION animes_generate() RETURNS VOID AS $animes_generate$
DECLARE
    -- Задаем количество записей в таблице anime
    number integer := 3;
    number_episodes integer := 0;
    data_issue timestamp with time zone;
    --current_hash varchar(64) := NULL;
BEGIN
    number_episodes :=  CEIL(RANDOM() * 100 + 1);
    data_issue := '2000-01-01 00:00:00'::timestamp +
                  ('2023-12-31 23:59:59'::timestamp - '2000-01-01 00:00:00'::timestamp) * RANDOM();
    INSERT INTO animes (id,anime_name, name_russian, name_english, name_japanese, name_synonyms,
                        anime_status, episodes, episodes_aired, aired_on, released_on, duration,
                        updated_at, next_episode_at, image)
    SELECT
        gen_random_uuid() AS id,
        md5(random()::text)::char(10) AS anime_name,
        md5(random()::text)::char(10) AS name_russian,
        ARRAY[md5(random()::text)::char(10)] AS name_english,
        ARRAY[md5(random()::text)::char(10)]AS name_japanese,
        ARRAY[md5(random()::text)::char(10)]AS name_synonyms,
        (ARRAY['anons', 'ongoing', 'released'])[floor(random() * 3 + 1)]::text::status AS anime_status,
        number_episodes AS episodes,
        CEIL(RANDOM() * number_episodes) AS episodes_aired,
        data_issue AS aired_on,
        data_issue + ('2023-12-31 23:59:59'::timestamp - data_issue) * RANDOM() AS released_on,
        CEIL(RANDOM() * 196 + 5) AS duration,
        data_issue + ('2023-12-31 23:59:59'::timestamp - data_issue) * RANDOM() AS updated_at,
        md5(random()::text)::char(10) AS next_episode_at,
        image_generate() AS image
    FROM
        generate_series(1, number);
END;
$animes_generate$ LANGUAGE plpgsql;

SELECT animes_generate();

SELECT * FROM animes;