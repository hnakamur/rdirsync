package rdirsync

import (
	"os/user"
	"strconv"
	"sync"
)

type userGroupDB struct {
	muUser       sync.Mutex
	nameToUidMap map[string]int
	uidToNameMap map[int]string

	muGroup      sync.Mutex
	nameToGidMap map[string]int
	gidToNameMap map[int]string
}

func newUserGroupDB() *userGroupDB {
	return &userGroupDB{
		nameToUidMap: make(map[string]int),
		uidToNameMap: make(map[int]string),
		nameToGidMap: make(map[string]int),
		gidToNameMap: make(map[int]string),
	}
}

func (db *userGroupDB) LookupUser(name string) (int, error) {
	db.muUser.Lock()
	defer db.muUser.Unlock()

	uid, exists := db.nameToUidMap[name]
	if exists {
		return uid, nil
	}

	user_, err := user.Lookup(name)
	if err != nil {
		return -1, err
	}
	uid, err = strconv.Atoi(user_.Uid)
	if err != nil {
		return -1, err
	}
	db.cacheUser(uid, name)

	return uid, nil
}

func (db *userGroupDB) LookupUid(uid int) (string, error) {
	db.muUser.Lock()
	defer db.muUser.Unlock()

	name, exists := db.uidToNameMap[uid]
	if exists {
		return name, nil
	}

	user_, err := user.LookupId(strconv.Itoa(uid))
	if err != nil {
		return "", err
	}
	name = user_.Name
	db.cacheUser(uid, name)

	return name, nil
}

func (db *userGroupDB) cacheUser(uid int, name string) {
	db.nameToUidMap[name] = uid
	db.uidToNameMap[uid] = name
}

func (db *userGroupDB) LookupGroup(name string) (int, error) {
	db.muGroup.Lock()
	defer db.muGroup.Unlock()

	gid, exists := db.nameToGidMap[name]
	if exists {
		return gid, nil
	}

	group, err := user.LookupGroup(name)
	if err != nil {
		return -1, err
	}
	gid, err = strconv.Atoi(group.Gid)
	if err != nil {
		return -1, err
	}
	db.cacheGroup(gid, name)

	return gid, nil
}

func (db *userGroupDB) LookupGid(gid int) (string, error) {
	db.muGroup.Lock()
	defer db.muGroup.Unlock()

	name, exists := db.gidToNameMap[gid]
	if exists {
		return name, nil
	}

	group, err := user.LookupGroupId(strconv.Itoa(gid))
	if err != nil {
		return "", err
	}
	name = group.Name
	db.cacheGroup(gid, name)

	return name, nil
}

func (db *userGroupDB) cacheGroup(gid int, name string) {
	db.nameToUidMap[name] = gid
	db.uidToNameMap[gid] = name
}
