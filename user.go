package rdirsync

import (
	"os/user"
	"strconv"
	"sync"
)

type userGroupDB struct {
	muUser       sync.Mutex
	nameToUidMap map[string]uint32
	uidToNameMap map[uint32]string

	muGroup      sync.Mutex
	nameToGidMap map[string]uint32
	gidToNameMap map[uint32]string
}

func newUserGroupDB() *userGroupDB {
	return &userGroupDB{
		nameToUidMap: make(map[string]uint32),
		uidToNameMap: make(map[uint32]string),
		nameToGidMap: make(map[string]uint32),
		gidToNameMap: make(map[uint32]string),
	}
}

func (db *userGroupDB) LookupUser(name string) (uint32, error) {
	db.muUser.Lock()
	defer db.muUser.Unlock()

	uid, exists := db.nameToUidMap[name]
	if exists {
		return uid, nil
	}

	user_, err := user.Lookup(name)
	if err != nil {
		return 0, err
	}
	uid, err = parseUint32(user_.Uid)
	if err != nil {
		return 0, err
	}
	db.cacheUser(uid, name)

	return uid, nil
}

func (db *userGroupDB) LookupUid(uid uint32) (string, error) {
	db.muUser.Lock()
	defer db.muUser.Unlock()

	name, exists := db.uidToNameMap[uid]
	if exists {
		return name, nil
	}

	user_, err := user.LookupId(formatUint32(uid))
	if err != nil {
		return "", err
	}
	name = user_.Username
	db.cacheUser(uid, name)

	return name, nil
}

func (db *userGroupDB) cacheUser(uid uint32, name string) {
	db.nameToUidMap[name] = uid
	db.uidToNameMap[uid] = name
}

func (db *userGroupDB) LookupGroup(name string) (uint32, error) {
	db.muGroup.Lock()
	defer db.muGroup.Unlock()

	gid, exists := db.nameToGidMap[name]
	if exists {
		return gid, nil
	}

	group, err := user.LookupGroup(name)
	if err != nil {
		return 0, err
	}
	gid, err = parseUint32(group.Gid)
	if err != nil {
		return 0, err
	}
	db.cacheGroup(gid, name)

	return gid, nil
}

func (db *userGroupDB) LookupGid(gid uint32) (string, error) {
	db.muGroup.Lock()
	defer db.muGroup.Unlock()

	name, exists := db.gidToNameMap[gid]
	if exists {
		return name, nil
	}

	group, err := user.LookupGroupId(formatUint32(gid))
	if err != nil {
		return "", err
	}
	name = group.Name
	db.cacheGroup(gid, name)

	return name, nil
}

func (db *userGroupDB) cacheGroup(gid uint32, name string) {
	db.nameToUidMap[name] = gid
	db.uidToNameMap[gid] = name
}

func formatUint32(i uint32) string {
	return strconv.FormatUint(uint64(i), 10)
}

func parseUint32(s string) (uint32, error) {
	i, err := strconv.ParseUint(s, 10, 32)
	return uint32(i), err
}
