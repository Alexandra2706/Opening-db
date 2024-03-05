# https://shikimori.one/api/animes/{id}

## id

Type: `int`

## name

Type: `str`

Description: название анимэ

## russian

Type: `str`

Description: название анимэ на русском 

## english

Type: `str[]`

Description: название анимэ на русском 

## japanese

Type: `str[]`

Description: название анимэ на японском 

## synonyms

Type: `str[]`

Description: альтернативные названия 

## image

Type: `dict{str : str}`

Description: url адрес изображения на сайте 

## url

Type: `str`

Description: url адрес страницы на сайте shikimori

## kind

Type: `str`

Description: Может быть одним из: tv, movie, ova, ona, special, tv_special, music, pv, cm, tv_13, tv_24, tv_48

## score

Type: `double`

Description: рейтинг берется из myanimelist

## status

Type: `str`

Description: Может быть одним из: anons, ongoing, released

## episodes

Type: `int`

Description: количество серий

## episodes_aired

Type: `int`

Description: Количество вышедших эпизодов, актульно только для `status == ongoing`

## aired_on

Type: `date`

Description: начало выпуска

## released_on

Type: `date`

Description: конец выпуска

## rating

Type: `str`

Description: ограничение по возрасту

## duration

Type: `int`

Description: длительность серии в минутах

## description

Type: `str`

Description: описание анимэ

## description_html

Type: `str`

Description: описание анимэ с html тегами

## franchise

Type: `str`

Description: франшиза

## myanimelist_id

Type: `int`

Description: id с сайта myanimelist

## updated_at

Type: `date`

Description: дата обновления

## fansubbers

Type: `str[]`

Description: субтитры

## licensors

Type: `str[]`

Description: лицензиары

## genres

Type: `dict[]`

Description: жанры

## studios

Type: `dict[]`

Description: студии, выпустившие анимэ

## videos

Type: `dict[]`

Description: видео

## screenshots

Type: `dict[]`

Description: кадры

```json
{
    "id": int,
    "name": str, // название анимэ
    "russian": str, // название на русском 
    "english": str[], // название на английском
    "japanese": str[], // название на японском
    "synonyms": str[], // альтернативные названия
    "image": {
        "original": str,
        "preview": str, 
    }, // url изображения на сайте
    "url": str, // адрес страницы на сайте shikimori
    "kind": str , // тип
    "score": double, // рейтинг берется из myanimelist
    "status": str, // статус: anons, ongoing, released
    "episodes": int, // количество серий
    "episodes_aired": 0, // количество вышедших эпизодов
    "aired_on": date, // начало выпуска
    "released_on": date, // конец выпуска
    "rating": str, // возрастной ценз    
    "duration": int, // длительность серии в минутах
    "description": str, // описание
    "description_html": str, // описание с тегами html
    "description_source": null, // Пока опускаем
    "franchise": str, // франшиза
    "myanimelist_id": int, //id с сайта myanimelist
    "updated_at": "2024-03-02T19:31:43.472+03:00", // дата обновления
    "next_episode_at": null,
    "fansubbers": str[], // субтитры    
    "licensors": str [],//лицензиары
    "genres": dict[], // жанры
    "studios": dict[], // студии
    "videos": dict[], // видео    
    "screenshots": dict[], // кадры
}
```