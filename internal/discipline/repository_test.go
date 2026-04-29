package discipline

import (
	"testing"
	"tournament_manager/internal/testutil"
	"tournament_manager/internal/tournament"
)

func TestAdd_Success(t *testing.T) {
	disciplineRepo, tournamentRepo := setupRepo(t)

	tournament, err := tournamentRepo.Add(tournament.Tournament{Name: "t", Location: "l"})

	if err != nil {
		t.Fatal("Unexpected tournament creation error")
	}

	d, err := disciplineRepo.Add(Discipline{Name: "MS", NoOfTeamPlayers: 2, TournamentID: tournament.ID})

	if err != nil {
		t.Fatalf("Expected success , got %v", err)
	}

	if d.Name != "MS" || d.NoOfTeamPlayers != 2 {
		t.Fatalf("Expected Add to be successful")
	}
}

func TestAdd_Failure(t *testing.T) {
	disciplineRepo, tournamentRepo := setupRepo(t)

	tournament, err := tournamentRepo.Add(tournament.Tournament{Name: "t", Location: "l"})

	if err != nil {
		t.Fatal("Unexpected tournament creation error")
	}

	_, err = disciplineRepo.Add(Discipline{TournamentID: tournament.ID})

	if err == nil {
		t.Fatal("Expected failure, got Success")
	}
}

func TestRepoGetAll_Success(t *testing.T) {
	disciplineRepo, tournamentRepo := setupRepo(t)

	tournament, err := tournamentRepo.Add(tournament.Tournament{Name: "t", Location: "l"})

	if err != nil {
		t.Fatal("Unexpected tournament creation error")
	}

	_, err = disciplineRepo.Add(Discipline{Name: "MS", NoOfTeamPlayers: 1, TournamentID: tournament.ID})
	_, err = disciplineRepo.Add(Discipline{Name: "WS", NoOfTeamPlayers: 1, TournamentID: tournament.ID})

	disciplines, err := disciplineRepo.GetAll(tournament.ID)
	if err != nil {
		t.Fatalf("Expected success, got %v", err)
	}

	if len(disciplines) != 2 {
		t.Fatalf("expected 2 entries returned, got %d", len(disciplines))
	}
}

func TestRepoShow_Success(t *testing.T) {
	disciplineRepo, tournamentRepo := setupRepo(t)

	tournament, err := tournamentRepo.Add(tournament.Tournament{Name: "t", Location: "l"})

	if err != nil {
		t.Fatal("Unexpected tournament creation error")
	}

	discipline, err := disciplineRepo.Add(Discipline{Name: "MS", NoOfTeamPlayers: 1, TournamentID: tournament.ID})

	_, err = disciplineRepo.Show(tournament.ID, discipline.ID)

	if err != nil {
		t.Fatalf("Expected success, got %v", err)
	}
}

func TestRepoShow_Failure(t *testing.T) {
	disciplineRepo, tournamentRepo := setupRepo(t)

	tournament, err := tournamentRepo.Add(tournament.Tournament{Name: "t", Location: "l"})

	if err != nil {
		t.Fatal("Unexpected tournament creation error")
	}

	_, err = disciplineRepo.Show(tournament.ID, "abc")

	if err == nil {
		t.Fatalf("Expected failure, got success")
	}
}

func TestRepoUpdate_Success(t *testing.T) {
	disciplineRepo, tournamentRepo := setupRepo(t)

	tournament, err := tournamentRepo.Add(tournament.Tournament{Name: "t", Location: "l"})

	if err != nil {
		t.Fatal("Unexpected tournament creation error")
	}

	discipline, err := disciplineRepo.Add(Discipline{Name: "MS", NoOfTeamPlayers: 1, TournamentID: tournament.ID})

	newDiscipline := Discipline{Name: "WD", NoOfTeamPlayers: 2, TournamentID: tournament.ID, ID: discipline.ID}

	_, err = disciplineRepo.Update(newDiscipline)

	if err != nil {
		t.Fatalf("Expected success, got %v", err)
	}

	disciplineUpdated, _ := disciplineRepo.Show(tournament.ID, discipline.ID)

	if disciplineUpdated.Name != "WD" || disciplineUpdated.NoOfTeamPlayers != 2 {
		t.Fatalf("Expected update to be successful")
	}
}

func TestRepoUpdate_Failure(t *testing.T) {
	disciplineRepo, tournamentRepo := setupRepo(t)

	tournament, err := tournamentRepo.Add(tournament.Tournament{Name: "t", Location: "l"})

	if err != nil {
		t.Fatal("Unexpected tournament creation error")
	}

	_, err = disciplineRepo.Add(Discipline{Name: "MS", NoOfTeamPlayers: 1, TournamentID: tournament.ID})

	newDiscipline := Discipline{Name: "WS", TournamentID: tournament.ID, ID: "abc"}

	_, err = disciplineRepo.Update(newDiscipline)

	if err == nil {
		t.Fatalf("Expected failure, got success")
	}
}

func TestRepoDelete_Success(t *testing.T) {
	disciplineRepo, tournamentRepo := setupRepo(t)

	tournament, err := tournamentRepo.Add(tournament.Tournament{Name: "t", Location: "l"})

	if err != nil {
		t.Fatal("Unexpected tournament creation error")
	}

	discipline, err := disciplineRepo.Add(Discipline{Name: "MS", TournamentID: tournament.ID})

	err = disciplineRepo.Delete(discipline.TournamentID, discipline.ID)

	if err != nil {
		t.Fatalf("Expected success, got %v", err)
	}

	_, err = disciplineRepo.Show(tournament.ID, discipline.ID)

	if err == nil {
		t.Fatalf("Expected deleted to be successful")
	}
}

func TestRepoDelete_Failure(t *testing.T) {
	disciplineRepo, tournamentRepo := setupRepo(t)

	tournament, err := tournamentRepo.Add(tournament.Tournament{Name: "t", Location: "l"})

	if err != nil {
		t.Fatal("Unexpected tournament creation error")
	}

	discipline, err := disciplineRepo.Add(Discipline{Name: "MS", TournamentID: tournament.ID})

	err = disciplineRepo.Delete(discipline.TournamentID, "abc")

	if err == nil {
		t.Fatalf("Expected failure, got success")
	}
}

func setupRepo(t *testing.T) (*Repository, *tournament.Repository) {
	dbConn := testutil.SetupTestRepo(t)
	disciplineRepo := NewRepository(dbConn)
	tournamentRepo := tournament.NewRepository(dbConn)

	return disciplineRepo, tournamentRepo
}
