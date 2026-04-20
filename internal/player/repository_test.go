package player

import (
	"testing"
	"tournament_manager/internal/testutil"
	tournament "tournament_manager/internal/tournament"
)

const (
	testPlayerID    = "11111111-1111-4111-8111-111111111111"
	missingPlayerID = "22222222-2222-4222-8222-222222222222"
	FirstName       = "foo"
	LastName        = "bar"
	Gender          = "male"
)

func TestRepository_Add(t *testing.T) {
	repo, tournamentRepo := initialiseRepo(t)

	tour, err := tournamentRepo.Add(tournament.Tournament{
		Name:     "Tournament",
		Location: "Berlin",
	})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	p := Player{
		FirstName:    FirstName,
		LastName:     LastName,
		Gender:       Gender,
		TournamentID: tour.ID,
	}

	created, err := repo.Add(p)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if created.ID == "" {
		t.Fatalf("Expected ID")
	}

	if created.FirstName != FirstName || created.LastName != LastName || created.Gender != Gender {
		t.Fatalf("Data are not stored properly")
	}
}

func TestRepository_GetAll(t *testing.T) {
	repo, tournamentRepo := initialiseRepo(t)
	tour := createTournament(t, tournamentRepo)

	_, _ = repo.Add(Player{
		FirstName:    "A",
		LastName:     "Player",
		Gender:       Gender,
		TournamentID: tour.ID,
	})
	_, _ = repo.Add(Player{
		FirstName:    "B",
		LastName:     "Player",
		Gender:       Gender,
		TournamentID: tour.ID,
	})

	list, err := repo.GetAll(tour.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(list) != 2 {
		t.Fatalf("expected 2 players, got %d", len(list))
	}
}

func TestRepository_Show_Success(t *testing.T) {
	repo, tournamentRepo := initialiseRepo(t)
	tour := createTournament(t, tournamentRepo)
	added := createPlayer(t, repo, tour.ID)

	got, err := repo.Show(added.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.FirstName != FirstName {
		t.Fatalf("expected %s, got %s", FirstName, got.FirstName)
	}

	if got.TournamentID != tour.ID {
		t.Fatalf("expected tournament ID %s, got %s", tour.ID, got.TournamentID)
	}
}

func TestRepository_Show_NotFound(t *testing.T) {
	repo, _ := initialiseRepo(t)

	_, err := repo.Show(missingPlayerID)
	if err == nil {
		t.Fatal("expected error")
	}

	if err.Error() != "player not found" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRepository_Update_Success(t *testing.T) {
	repo, tournamentRepo := initialiseRepo(t)
	tour := createTournament(t, tournamentRepo)
	added := createPlayer(t, repo, tour.ID)

	err := repo.Update(Player{
		ID:           added.ID,
		FirstName:    "baz",
		LastName:     "qux",
		Gender:       "female",
		TournamentID: tour.ID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updated, err := repo.Show(added.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if updated.FirstName != "baz" || updated.LastName != "qux" || updated.Gender != "female" {
		t.Fatalf("player was not updated properly: %+v", updated)
	}
}

func TestRepository_Update_NotFound(t *testing.T) {
	repo, _ := initialiseRepo(t)

	err := repo.Update(Player{
		ID:        missingPlayerID,
		FirstName: "baz",
		LastName:  "qux",
		Gender:    "female",
	})
	if err == nil {
		t.Fatal("expected error")
	}

	if err.Error() != "player not found" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRepository_Delete_Success(t *testing.T) {
	repo, tournamentRepo := initialiseRepo(t)
	tour := createTournament(t, tournamentRepo)
	added := createPlayer(t, repo, tour.ID)

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
	repo, _ := initialiseRepo(t)

	err := repo.Delete(missingPlayerID)
	if err == nil {
		t.Fatal("expected error")
	}

	if err.Error() != "player not found" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func initialiseRepo(t *testing.T) (*Repository, *tournament.Repository) {
	dbConn := testutil.SetupTestRepo(t)

	repo := NewRepository(dbConn)
	tournamentRepo := tournament.NewRepository(dbConn)
	return repo, tournamentRepo
}

func createTournament(t *testing.T, repo *tournament.Repository) tournament.Tournament {
	t.Helper()

	tour, err := repo.Add(tournament.Tournament{
		Name:     "Tournament",
		Location: "Berlin",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	return tour
}

func createPlayer(t *testing.T, repo *Repository, tournamentID string) Player {
	t.Helper()

	p, err := repo.Add(Player{
		FirstName:    FirstName,
		LastName:     LastName,
		Gender:       Gender,
		TournamentID: tournamentID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	return p
}
