package database

import (
	"co2/types"
)

func Containers() []types.Container {
	db, _ := Get()

	rows, err := db.Query("SELECT * FROM containers")
	handle(err)

	var containers []types.Container
	for rows.Next() {
		var out types.Container

		err = rows.Scan(&out.Id, &out.Uid, &out.Name, &out.ComposeFile, &out.CreatedAt)
		handle(err)

		containers = append(containers, out)
	}

	return containers
}

func InsertContainer(container types.Container) types.Container {
	db, _ := Get()

	stmt, err := db.Prepare("INSERT INTO containers(uid, name, compose_file) VALUES(?,?,?);")
	handle(err)

	res, err := stmt.Exec(container.Uid, container.Name, container.ComposeFile)
	handle(err)

	id, err := res.LastInsertId()
	handle(err)

	container.Id = id
	return container
}

func DeleteContainer(container types.Container) int64 {
	db, _ := Get()

	stmt, err := db.Prepare("DELETE FROM containers WHERE uid=? AND name=?;")
	handle(err)

	res, err := stmt.Exec(container.Uid, container.Name)
	handle(err)

	affect, err := res.RowsAffected()
	handle(err)

	return affect
}
