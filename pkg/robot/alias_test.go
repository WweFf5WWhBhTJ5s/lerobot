package robot

import "testing"

func TestAliasStoreSetAndResolve(t *testing.T) {
	s := NewAliasStore()
	if err := s.Set("alpha", "robot-1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	id, ok := s.Resolve("alpha")
	if !ok || id != "robot-1" {
		t.Fatalf("expected robot-1, got %q ok=%v", id, ok)
	}
}

func TestAliasStoreResolveUnknown(t *testing.T) {
	s := NewAliasStore()
	_, ok := s.Resolve("ghost")
	if ok {
		t.Fatal("expected not found")
	}
}

func TestAliasStoreConflict(t *testing.T) {
	s := NewAliasStore()
	s.Set("alpha", "robot-1")
	err := s.Set("alpha", "robot-2")
	if err == nil {
		t.Fatal("expected conflict error")
	}
}

func TestAliasStoreReassignSameRobot(t *testing.T) {
	s := NewAliasStore()
	s.Set("alpha", "robot-1")
	if err := s.Set("alpha", "robot-1"); err != nil {
		t.Fatalf("reassigning same robot should not error: %v", err)
	}
}

func TestAliasStoreRobotGetsNewAlias(t *testing.T) {
	s := NewAliasStore()
	s.Set("alpha", "robot-1")
	s.Set("beta", "robot-1") // replaces old alias
	_, ok := s.Resolve("alpha")
	if ok {
		t.Fatal("old alias should be removed")
	}
	id, ok := s.Resolve("beta")
	if !ok || id != "robot-1" {
		t.Fatalf("expected robot-1 under beta, got %q", id)
	}
}

func TestAliasStoreAliasFor(t *testing.T) {
	s := NewAliasStore()
	s.Set("alpha", "robot-1")
	alias, ok := s.AliasFor("robot-1")
	if !ok || alias != "alpha" {
		t.Fatalf("expected alpha, got %q ok=%v", alias, ok)
	}
}

func TestAliasStoreDelete(t *testing.T) {
	s := NewAliasStore()
	s.Set("alpha", "robot-1")
	s.Delete("alpha")
	_, ok := s.Resolve("alpha")
	if ok {
		t.Fatal("alias should be deleted")
	}
	_, ok = s.AliasFor("robot-1")
	if ok {
		t.Fatal("reverse mapping should be removed")
	}
}

func TestAliasStoreAll(t *testing.T) {
	s := NewAliasStore()
	s.Set("alpha", "robot-1")
	s.Set("beta", "robot-2")
	all := s.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
	if all["alpha"] != "robot-1" || all["beta"] != "robot-2" {
		t.Fatalf("unexpected entries: %v", all)
	}
}
