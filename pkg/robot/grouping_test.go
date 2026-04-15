package robot

import (
	"sort"
	"testing"
)

func TestGroupStoreAddAndMembers(t *testing.T) {
	gs := NewGroupStore()
	gs.Add("alpha", "r1")
	gs.Add("alpha", "r2")

	members := gs.Members("alpha")
	if len(members) != 2 {
		t.Fatalf("expected 2 members, got %d", len(members))
	}
}

func TestGroupStoreMembersUnknownGroup(t *testing.T) {
	gs := NewGroupStore()
	if gs.Members("nope") != nil {
		t.Fatal("expected nil for unknown group")
	}
}

func TestGroupStoreRemove(t *testing.T) {
	gs := NewGroupStore()
	gs.Add("beta", "r1")
	gs.Add("beta", "r2")
	gs.Remove("beta", "r1")

	members := gs.Members("beta")
	if len(members) != 1 || members[0] != "r2" {
		t.Fatalf("expected [r2], got %v", members)
	}
}

func TestGroupStoreRemoveEmptiesGroup(t *testing.T) {
	gs := NewGroupStore()
	gs.Add("gamma", "r1")
	gs.Remove("gamma", "r1")

	if gs.Members("gamma") != nil {
		t.Fatal("expected group to be removed when empty")
	}
}

func TestGroupStoreGroups(t *testing.T) {
	gs := NewGroupStore()
	gs.Add("a", "r1")
	gs.Add("b", "r1")
	gs.Add("c", "r2")

	groups := gs.Groups("r1")
	sort.Strings(groups)
	if len(groups) != 2 || groups[0] != "a" || groups[1] != "b" {
		t.Fatalf("expected [a b], got %v", groups)
	}
}

func TestGroupStoreDeleteGroup(t *testing.T) {
	gs := NewGroupStore()
	gs.Add("delta", "r1")
	gs.DeleteGroup("delta")

	if gs.Members("delta") != nil {
		t.Fatal("expected group to be deleted")
	}
}

func TestGroupStoreListGroups(t *testing.T) {
	gs := NewGroupStore()
	gs.Add("x", "r1")
	gs.Add("y", "r2")

	list := gs.ListGroups()
	if len(list) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(list))
	}
}

func TestGroupStoreString(t *testing.T) {
	gs := NewGroupStore()
	gs.Add("z", "r1")
	s := gs.String()
	if s == "" {
		t.Fatal("expected non-empty string")
	}
}
