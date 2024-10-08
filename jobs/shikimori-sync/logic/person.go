package logic

import (
	"context"
	"encoding/json"
	"log"
	"shikimori-sync/postgres"
	"shikimori-sync/s3"
	shikimori_api "shikimori-sync/shikimori-api"
)

func CreateOrUpdatePerson(pid int) {
	person, err := shikimori_api.GetPersonInfo(pid)
	if err != nil {
		log.Fatalf("Error in person %d: %q", pid, err.Error())
	}
	log.Printf("Person %d, name %q exists", pid, person.Name)

	// Добавляем image в s3 и images_table, если все ок, то получаем url для таблица person, иначе записываем nil в таблицу person
	err = s3.CreateOrUpdateImage(imageEndpoint + person.Image.Original)
	if err != nil {
		log.Fatalf("Error in anime image: %q", err.Error())
	}

	birthday, err := json.Marshal(person.BirthOn)
	if err != nil {
		log.Fatalf("Error in json birthday: %q", err.Error())
	}

	log.Println("birthday:", string(birthday))

	deceased, err := json.Marshal(person.DeceasedOn)
	if err != nil {
		log.Fatalf("Error in json deceased: %q", err.Error())
	}

	log.Println("deceased:", string(deceased))

	//преобразуем [][] в map
	var gRoles = map[string]int{}
	var grouppedRoles [][]interface{}

	err = json.Unmarshal(*person.GrouppedRoles, &grouppedRoles)
	if err != nil {
		log.Println("error:", err)
	}

	for _, i := range grouppedRoles {
		gRoles[i[0].(string)] = int(i[1].(float64))
	}

	log.Println("gRoles:", gRoles)

	row := postgres.Conn.QueryRow(context.Background(), `
			INSERT INTO public.person(
				people_name,  russian,        japanese,  image,			                          
				shikimori_id, job_title,      birthday,  deceased,
			    website,      groupped_roles, producer,  mangaka,   
			    seyu                     
				) 
			VALUES ($1, $2,  $3,  $4,
			        $5, $6,  $7,  $8,
			        $9, $10, $11, $12,
			        $13
			        )
			ON CONFLICT (shikimori_id)
			DO UPDATE SET
				people_name = $1, russian = $2, japanese = $3, image = $4,
				shikimori_id = $5, job_title = $6, birthday = to_jsonb($7), deceased = to_jsonb($8),
				website = $9, groupped_roles = to_jsonb($10), producer = $11, mangaka = $12, 
				seyu = $13
			RETURNING  id;`,
		person.Name, person.Russian, person.Japanese, imageEndpoint+person.Image.Original,
		person.Id, person.JobTitle, birthday, deceased,
		person.Website, gRoles, person.Producer, person.Mangaka,
		person.Seyu,
	)

	var id string
	err = row.Scan(&id)
	if err != nil {
		log.Fatalf("Error when saving person %d: %q", pid, err.Error())
	}
	log.Printf("Person %q saved!", person.Name)
}
