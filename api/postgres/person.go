package postgres

import (
	"api/structures"
	"context"
)

func ListPerson(limit int, offset int) ([]structures.PersonInfo, error) {

	rows, err := Conn.Query(context.Background(), `
		SELECT id, people_name, russian, japanese,
		       image, shikimori_id, job_title, birthday,
		       deceased, website, groupped_roles, producer,
		       mangaka, seyu, updated_at
		FROM public.person LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rowSlice []structures.PersonInfo
	for rows.Next() {
		var r structures.PersonInfo
		var uncovertedGroupedRoles [][]any
		err := rows.Scan(
			&r.Id, &r.PeopleName, &r.Russian, &r.Japanese,
			&r.Image, &r.ShikimoriId, &r.JobTitle, &r.Birthday,
			&r.Deceased, &r.Website, &uncovertedGroupedRoles, &r.Producer,
			&r.Mangaka, &r.Seyu, &r.UpdatedAt.T,
		)
		if err != nil {
			return nil, err
		}
		// ToDo: переделать внутри шикимори синк
		r.GrouppedRoles = map[string]int{}
		for _, i := range uncovertedGroupedRoles {
			r.GrouppedRoles[i[0].(string)] = int(i[1].(float64))
		}
		rowSlice = append(rowSlice, r)
	}

	return rowSlice, nil
}
