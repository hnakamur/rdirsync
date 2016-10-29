package rdirsync

import (
	"os/user"
	"strconv"
	"sync"
)

type userGroupDB struct {
	muUser       sync.RWMutex
	nameToUidMap map[string]uint32
	uidToNameMap map[uint32]string

	muGroup      sync.RWMutex
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
	db.muUser.RLock()
	uid, exists := db.nameToUidMap[name]
	db.muUser.RUnlock()
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
	db.muUser.RLock()
	name, exists := db.uidToNameMap[uid]
	db.muUser.RUnlock()
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
	db.muUser.Lock()
	db.nameToUidMap[name] = uid
	db.uidToNameMap[uid] = name
	db.muUser.Unlock()
}

func (db *userGroupDB) LookupGroup(name string) (uint32, error) {
	db.muGroup.RLock()
	gid, exists := db.nameToGidMap[name]
	db.muGroup.RUnlock()
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
	db.muGroup.RLock()
	name, exists := db.gidToNameMap[gid]
	db.muGroup.RUnlock()
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
	db.muGroup.Lock()
	db.nameToUidMap[name] = gid
	db.uidToNameMap[gid] = name
	db.muGroup.Unlock()
}

func formatUint32(i uint32) string {
	return strconv.FormatUint(uint64(i), 10)
}

func parseUint32(s string) (uint32, error) {
	i, err := strconv.ParseUint(s, 10, 32)
	return uint32(i), err
}
