package attendant

import (
	"slices"
	"testing"
	"tournament_manager/internal/player"
	"tournament_manager/internal/testutil"
	"tournament_manager/internal/tournament"
)

func TestRepository_Add_Success(t *testing.T) {
	repo, playerRepo, tournamentRepo := initialiseRepo(t)

	player, _ := playerRepo.Add(player.Player{FirstName: "foo", LastName: "bar", Gender: "male"})
	tournament, _ := tournamentRepo.Add(tournament.Tournament{Location: "Berlin", Name: "Coppi Cup"})

	attendant, err := repo.Add(Attendant{TournamentID: tournament.ID, PlayerID: player.ID})
	if err != nil {
		t.Fatal("Unexpected error")
	}

	if attendant.ID == "" {
		t.Fatal("Failed to create attendant")
	}
}

func TestRepository_Add_Failure(t *testing.T) {
	repo, playerRepo, _ := initialiseRepo(t)

	player, _ := playerRepo.Add(player.Player{FirstName: "foo", LastName: "bar", Gender: "male"})

	_, err := repo.Add(Attendant{TournamentID: "random_id", PlayerID: player.ID})
	if err == nil {
		t.Fatalf("unexpected error")
	}
}

func TestRepository_GetAll_Success(t *testing.T) {
	repo, playerRepo, tournamentRepo := initialiseRepo(t)

	player1, _ := playerRepo.Add(player.Player{FirstName: "Martin", LastName: "bar", Gender: "male"})
	player2, _ := playerRepo.Add(player.Player{FirstName: "Lisa", LastName: "bar", Gender: "female"})
	tournament, _ := tournamentRepo.Add(tournament.Tournament{Location: "Berlin", Name: "Coppi Cup"})

	repo.Add(Attendant{TournamentID: tournament.ID, PlayerID: player1.ID})
	repo.Add(Attendant{TournamentID: tournament.ID, PlayerID: player2.ID})

	attendants, err := repo.GetAll(tournament.ID)

	if err != nil {
		t.Fatalf("unexpected err %v", err)
	}

	if len(attendants) != 2 {
		t.Fatalf("expected 2 attendants, got %d", len(attendants))
	}

	firstNames := make([]string, 0, len(attendants))
	for _, attendant := range attendants {
		firstNames = append(firstNames, attendant.FirstName)
	}

	expectedFirstNames := []string{"Lisa", "Martin"}

	if !slices.Equal(firstNames, expectedFirstNames) {
		t.Fatalf("Expected %v, got %v", expectedFirstNames, firstNames)
	}
}

func TestRepository_Show_Success(t *testing.T) {
	repo, playerRepo, tournamentRepo := initialiseRepo(t)

	player, _ := playerRepo.Add(player.Player{FirstName: "Martin", LastName: "bar", Gender: "male"})
	tournament, _ := tournamentRepo.Add(tournament.Tournament{Location: "Berlin", Name: "Coppi Cup"})

	attendant, _ := repo.Add(Attendant{TournamentID: tournament.ID, PlayerID: player.ID})

	displayAttendant, err := repo.Show(tournament.ID, attendant.ID)

	if err != nil {
		t.Fatalf("unexpected err %v", err)
	}

	if displayAttendant.FirstName != "Martin" {
		t.Fatalf("Expected Martin, got %v", displayAttendant.FirstName)
	}
}

func TestRepository_Show_NotFound(t *testing.T) {
	repo, _, tournamentRepo := initialiseRepo(t)

	tournament, _ := tournamentRepo.Add(tournament.Tournament{Location: "Berlin", Name: "Coppi Cup"})

	_, err := repo.Show(tournament.ID, "random_id")

	if err == nil {
		t.Fatalf("Expected error to occur")
	}
}

func TestRepository_Delete_Success(t *testing.T) {
	repo, playerRepo, tournamentRepo := initialiseRepo(t)

	player, _ := playerRepo.Add(player.Player{FirstName: "Martin", LastName: "bar", Gender: "male"})
	tournament, _ := tournamentRepo.Add(tournament.Tournament{Location: "Berlin", Name: "Coppi Cup"})

	attendant, _ := repo.Add(Attendant{TournamentID: tournament.ID, PlayerID: player.ID})

	err := repo.Delete(tournament.ID, attendant.ID)

	if err != nil {
		t.Fatalf("unexpected err %v", err)
	}

	_, err = repo.Show(tournament.ID, attendant.ID)

	if err == nil {
		t.Fatal("expected attendant to be deleted")
	}
}

func TestRepository_Delete_NotFound(t *testing.T) {
	repo, _, tournamentRepo := initialiseRepo(t)

	tournament, _ := tournamentRepo.Add(tournament.Tournament{Location: "Berlin", Name: "Coppi Cup"})

	err := repo.Delete(tournament.ID, "random_id")

	if err == nil {
		t.Fatalf("Expected error to occur")
	}
}

func initialiseRepo(t *testing.T) (*Repository, *player.Repository, *tournament.Repository) {
	dbConn := testutil.SetupTestRepo(t)

	repo := NewRepository(dbConn)
	playerRepo := player.NewRepository(dbConn)
	tournamentRepo := tournament.NewRepository(dbConn)

	return repo, playerRepo, tournamentRepo
}
