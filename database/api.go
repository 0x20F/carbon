package database

import (
	"co2/types"
)

func Containers() []types.Container {
	db, _ := Get()

	rows, err := db.Query("SELECT * FROM containers;")
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

func Stores() []types.Store {
	db, _ := Get()

	rows, err := db.Query("SELECT * FROM stores;")
	handle(err)

	var stores []types.Store
	for rows.Next() {
		var out types.Store

		err = rows.Scan(&out.Id, &out.Uid, &out.Path, &out.CreatedAt)
		handle(err)

		stores = append(stores, out)
	}

	return stores
}

func AddContainer(container types.Container) types.Container {
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

func AddStore(store types.Store) types.Store {
	db, _ := Get()

	stmt, err := db.Prepare("INSERT INTO stores(uid, path) VALUES(?,?);")
	handle(err)

	res, err := stmt.Exec(store.Uid, store.Path)
	handle(err)

	id, err := res.LastInsertId()
	handle(err)

	store.Id = id
	return store
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
