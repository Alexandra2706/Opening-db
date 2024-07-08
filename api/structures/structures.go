package structures

import (
	"github.com/google/uuid"
)

type AnimeShortInfo struct {
	Id            uuid.UUID   `json:"id"`
	AnimeName     string      `json:"anime_name"`
	NameRussian   string      `json:"name_russian"`
	NameEnglish   *[]string   `json:"name_english"`
	NameJapanese  *[]string   `json:"name_japanese"`
	NameSynonyms  *[]string   `json:"name_synonyms"`
	AnimeStatus   AnimeStatus `json:"anime_status"`
	Episodes      *int        `json:"episodes"`
	EpisodesAired *int        `json:"episodes_aired"`
	AiredOn       Time        `json:"aired_on"`
	ReleasedOn    Time        `json:"released_on"`
	Duration      *int        `json:"duration"`

	//LicensorRu *[]string `json:"licensor_ru"`
	//Franchise  *string   `json:"franchise"`

	UpdatedAt     Time         `json:"updated_at"`
	NextEpisodeAt *string      `json:"next_episode_at"`
	Image         *string      `json:"image"`
	Genres        *[]string    `json:"genres"`
	Studios       *[]uuid.UUID `json:"studios"`
	Videos        *[]string    `json:"videos"`
	Screenshots   *[]string    `json:"screenshots"`

	ShikimoriId   string    `json:"shikimori_id"`
	ShikimoriKind AnimeKind `json:"shikimori_kind"`
}

/*
licensors_ru jsonb, --лицензировано
franchise jsonb, --франшиза

videos varchar(255)[], --REFERENCES Video (id), --эпизоды

-- shikimori data:
shikimori_kind kind NOT NULL, --тип анимэ на сайте shikimori Посмотреть джненрики и темплейты
shikimori_rating rating, --возрастной ценз
shikimori_description varchar, --описание на сайте shikimori
shikimori_description_html varchar, --описание с тегами html на сайте shikimori
shikimori_last_revision timestamp with time zone, --дата обновления на сайте shikimori

-- myanimelist data:
myanimelist_id integer UNIQUE NOT NULL, --id с сайта myanimelist
myanimelist_score real --рейтинг берется из myanimelist

--description_source null, --Пока опускаем не понятно, что это
);
*/
