package tournament

import (
	"testing"
	"tournament_manager/internal/testutil"
)

const (
	testTournamentID    = "11111111-1111-4111-8111-111111111111"
	missingTournamentID = "22222222-2222-4222-8222-222222222222"
)

func TestRepository_Add(t *testing.T) {
	repo := initialiseRepo(t)
	tournament := Tournament{
		Name:     "Test Cup",
		Location: "SG",
	}

	created, err := repo.Add(tournament)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if created.ID == "" {
		t.Fatalf("Expected ID")
	}
}

func TestRepository_GetAll(t *testing.T) {
	repo := initialiseRepo(t)

	_, _ = repo.Add(Tournament{Name: "A", Location: "SG"})
	_, _ = repo.Add(Tournament{Name: "B", Location: "MY"})

	list, err := repo.GetAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(list) != 2 {
		t.Fatalf("expected 2 tournaments, got %d", len(list))
	}
}

func TestRepository_Show_Success(t *testing.T) {
	repo := initialiseRepo(t)

	added, _ := repo.Add(Tournament{Name: "ShowTest", Location: "SG"})

	got, err := repo.Show(added.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.Name != "ShowTest" {
		t.Fatalf("expected ShowTest, got %s", got.Name)
	}
}

func TestRepository_Show_NotFound(t *testing.T) {
	repo := initialiseRepo(t)

	_, err := repo.Show(missingTournamentID)
	if err == nil {
		t.Fatal("expected error")
	}

	if err.Error() != "tournament not found" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRepository_Update_Success(t *testing.T) {
	repo := initialiseRepo(t)

	added, _ := repo.Add(Tournament{Name: "Old", Location: "SG"})

	err := repo.Update(Tournament{
		ID:       added.ID,
		Name:     "New",
		Location: "MY",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updated, _ := repo.Show(added.ID)

	if updated.Name != "New" {
		t.Fatalf("expected New, got %s", updated.Name)
	}
}

func TestRepository_Update_NotFound(t *testing.T) {
	repo := initialiseRepo(t)

	err := repo.Update(Tournament{ID: missingTournamentID, Name: "X", Location: "Y"})
	if err == nil {
		t.Fatal("expected error")
	}

	if err.Error() != "tournament not found" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRepository_Delete_Success(t *testing.T) {
	repo := initialiseRepo(t)

	added, _ := repo.Add(Tournament{Name: "DeleteMe", Location: "SG"})

	err := repo.Delete(added.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = repo.Show(added.ID)
	if err == nil {
		t.Fatal("expected not found after delete")
	}
}

func TestRepository_Delete_NotFound(t *testing.T) {
	repo := initialiseRepo(t)

	err := repo.Delete(missingTournamentID)
	if err == nil {
		t.Fatal("expected error")
	}

	if err.Error() != "tournament not found" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func initialiseRepo(t *testing.T) *Repository {
	dbConn := testutil.SetupTestRepo(t)

	repo := NewRepository(dbConn)
	return repo
}
