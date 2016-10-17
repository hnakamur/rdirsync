package rdirsync

import "testing"

func TestLookupUser(t *testing.T) {
	db := newUserGroupDB()

	uid, err := db.LookupUser("root")
	if err != nil {
		t.Fatalf("failed to lookup user \"root\"; %s", err)
	}
	if uid != 0 {
		t.Errorf("unexpected uid, got %d, want %d", uid, 0)
	}

	uid, err = db.LookupUser("root")
	if err != nil {
		t.Fatalf("failed to lookup user \"root\"; %s", err)
	}
	if uid != 0 {
		t.Errorf("unexpected uid, got %d, want %d", uid, 0)
	}
}

func TestLookupUid(t *testing.T) {
	db := newUserGroupDB()

	name, err := db.LookupUid(0)
	if err != nil {
		t.Fatalf("failed to lookup uid 0; %s", err)
	}
	if name != "root" {
		t.Errorf("unexpected uid, got %d, want %d", name, "root")
	}

	name, err = db.LookupUid(0)
	if err != nil {
		t.Fatalf("failed to lookup uid 0; %s", err)
	}
	if name != "root" {
		t.Errorf("unexpected uid, got %d, want %d", name, "root")
	}
}

func TestLookupGroup(t *testing.T) {
	db := newUserGroupDB()

	gid, err := db.LookupGroup("root")
	if err != nil {
		t.Fatalf("failed to lookup group \"root\"; %s", err)
	}
	if gid != 0 {
		t.Errorf("unexpected gid, got %d, want %d", gid, 0)
	}

	gid, err = db.LookupGroup("root")
	if err != nil {
		t.Fatalf("failed to lookup group \"root\"; %s", err)
	}
	if gid != 0 {
		t.Errorf("unexpected gid, got %d, want %d", gid, 0)
	}
}

func TestLookupGid(t *testing.T) {
	db := newUserGroupDB()

	name, err := db.LookupGid(0)
	if err != nil {
		t.Fatalf("failed to lookup uid 0; %s", err)
	}
	if name != "root" {
		t.Errorf("unexpected uid, got %d, want %d", name, "root")
	}

	name, err = db.LookupGid(0)
	if err != nil {
		t.Fatalf("failed to lookup uid 0; %s", err)
	}
	if name != "root" {
		t.Errorf("unexpected uid, got %d, want %d", name, "root")
	}
}
