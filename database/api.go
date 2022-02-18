package database

import (
	"co2/types"
)

// Gets all the containers currently registered in the database
// and maps them to our own custom Container structure.
func Containers() []types.Container {
	db, _ := Get()

	rows, err := db.Query("SELECT * FROM containers;")
	handle(err)

	var containers []types.Container
	for rows.Next() {
		var out types.Container

		err = rows.Scan(
			&out.Id,
			&out.DockerUid,
			&out.Uid,
			&out.Name,
			&out.Image,
			&out.ServiceName,
			&out.ComposeFile,
			&out.Ports,
			&out.Status,
			&out.CreatedAt,
		)
		handle(err)

		containers = append(containers, out)
	}

	return containers
}

// Gets all the stores currently registered in the database
// and maps them to our own custom Store structure.
func Stores() []types.Store {
	db, _ := Get()

	rows, err := db.Query("SELECT * FROM stores;")
	handle(err)

	var stores []types.Store
	for rows.Next() {
		var out types.Store

		err = rows.Scan(&out.Id, &out.Uid, &out.Path, &out.Env, &out.CreatedAt)
		handle(err)

		stores = append(stores, out)
	}

	return stores
}

// Adds a new container to the database.
// And updates the ID of the provided container to match the
// inserted one.
func AddContainer(container types.Container) types.Container {
	db, _ := Get()

	stmt, err := db.Prepare("INSERT INTO containers(docker_uid, uid, name, image, service_name, compose_file, ports, status) VALUES(?,?,?,?,?,?,?,?);")
	handle(err)

	res, err := stmt.Exec(
		container.DockerUid,
		container.Uid,
		container.Name,
		container.Image,
		container.ServiceName,
		container.ComposeFile,
		container.Ports,
		container.Status,
	)
	handle(err)

	id, err := res.LastInsertId()
	handle(err)

	container.Id = id
	return container
}

// Adds a new store to the database.
// And updates the ID of the provided store to match the
// inserted one.
func AddStore(store types.Store) types.Store {
	db, _ := Get()

	stmt, err := db.Prepare("INSERT INTO stores(uid, path, env) VALUES(?,?,?);")
	handle(err)

	res, err := stmt.Exec(store.Uid, store.Path, store.Env)
	handle(err)

	id, err := res.LastInsertId()
	handle(err)

	store.Id = id
	return store
}

// Deletes a container from the database and returns the
// deleted ID.
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

// Deletes a store from the database and returns the
// deleted ID.
func DeleteStore(store types.Store) int64 {
	db, _ := Get()

	stmt, err := db.Prepare("DELETE FROM stores WHERE uid=?;")
	handle(err)

	res, err := stmt.Exec(store.Uid)
	handle(err)

	affect, err := res.RowsAffected()
	handle(err)

	return affect
}
