package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strconv"
)

type game_info struct {
	game_name string
	teams     []game_team
}
type game_team struct {
	team_name   string
	team_person []string
}

// 注册队伍
func ReTeamCompetition(tName string, user string, playerID string, gameId string) string {
	db, err := sql.Open("mysql", databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	//录入用户信息
	insertPlayer(db, user, playerID)
	//录入队伍名称
	insertTeam(db, tName)
	//登记队伍和比赛
	insertCompetitionTeam(db, gameId, getTeamID(db, tName))
	//登记比赛和人员
_:
	insertTeamPlayer(db, getTeamID(db, tName), getPlayerID(db, user))
	//登记比赛和人员

	return "登记成功"

}

func insertPlayer(db *sql.DB, playerName string, playerID string) string {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM Player WHERE id = ?", playerID).Scan(&count)
	if err != nil {
		return "注册失败"
	}
	if !(count > 0) {
		_, err := db.Exec("INSERT INTO Player (PlayerName,id) VALUES (?,?)", playerName, playerID)
		if err != nil {
			return "注册失败"
		}
		return "注册成功"
	}
	return "已注册"
}

func insertTeam(db *sql.DB, TeamName string) string {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM Team WHERE TeamName = ?", TeamName).Scan(&count)
	if err != nil {
		return "队伍注册失败"
	}
	if !(count > 0) {
		_, err := db.Exec("INSERT INTO  Team  (TeamName) VALUES (?)", TeamName)
		if err != nil {
			return "队伍注册失败"
		}
		return "队伍注册成功"
	}
	return "队伍已注册"
}

func insertCompetitionTeam(db *sql.DB, competitionID, teamID string) {
	_, err := db.Exec("INSERT INTO CompetitionTeam (CompetitionID, TeamID) VALUES (?, ?)", competitionID, teamID)
	fmt.Print(err)
}

func getTeamID(db *sql.DB, teamName string) string {
	var teamID int
	err := db.QueryRow("SELECT TeamID FROM Team WHERE TeamName = ?", teamName).Scan(&teamID)
	if err != nil {
		log.Fatal(err)
	}
	return strconv.Itoa(teamID)
}

func insertTeamPlayer(db *sql.DB, teamID, playerID string) string {
	var count int
	err2 := db.QueryRow("SELECT COUNT(*) FROM TeamPlayer WHERE TeamID = ? AND PlayerID=?", teamID, playerID).Scan(&count)
	if err2 != nil {
		return "注册失败"
	}
	if !(count > 0) {
		_, err := db.Exec("INSERT INTO TeamPlayer (TeamID, PlayerID) VALUES (?, ?)", teamID, playerID)
		if err != nil {
			log.Fatal(err)
		}
		return "注册成功"
	}
	return "已报名"
}

func getPlayerID(db *sql.DB, Name string) string {
	var playerId int
	err := db.QueryRow("SELECT PlayerId FROM Player WHERE PlayerName = ?", Name).Scan(&playerId)
	if err != nil {
		log.Fatal(err)
	}
	return strconv.Itoa(playerId)
}

func getComptitionTeam(CompetitionID string) string {
	db, err := sql.Open("mysql", databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var game game_info

	query1 := `
		SELECT CompetitionName
		FROM Competition
		WHERE CompetitionID = ?
	`
	Name, err := db.Query(query1, CompetitionID)
	if err != nil {
		log.Fatal("无法找到竞赛id")
	}
	for Name.Next() {
		Name.Scan(&game.game_name)
	}

	if err != nil {
		log.Fatal("无法找到竞赛id")
	}

	query := `
		SELECT Team.TeamName
		FROM CompetitionTeam
		JOIN Team ON CompetitionTeam.TeamID = Team.TeamID
		WHERE CompetitionTeam.CompetitionID = ?
	`

	rows, err := db.Query(query, CompetitionID)
	if err != nil {
		log.Fatal("无法找到竞赛id")
	}
	defer rows.Close()
	var team game_team
	for rows.Next() {
		var team_info game_team
		err := rows.Scan(&team.team_name)
		if err != nil {
			log.Fatal(err)
		}

		team_info = getTeamPlayer(db, team)
		game.teams = append(game.teams, team_info)

	}

	var condata = ""
	condata = condata + "比赛名称:" + game.game_name
	for _, team := range game.teams {
		condata = condata + "\n参赛队伍：" + team.team_name + "\n"
		condata = condata + "参赛人员："

		for _, user := range team.team_person {
			condata = condata + user + " "
		}
		condata = condata + "\n" + "----------------"
	}

	return condata
}

func getTeamPlayer(db *sql.DB, Team game_team) game_team {
	var TeamId int
	err := db.QueryRow("SELECT TeamID FROM Team WHERE TeamName = ?", Team.team_name).Scan(&TeamId)
	if err != nil {
		log.Fatal(err)
	}
	query := `
		SELECT Player.PlayerName
FROM TeamPlayer
JOIN Player ON TeamPlayer.PlayerID = Player.PlayerID
WHERE TeamPlayer.TeamID = ?;

	`
	rows, err := db.Query(query, TeamId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var person = ""

	for rows.Next() {
		err := rows.Scan(&person)
		Team.team_person = append(Team.team_person, person)
		if err != nil {
			log.Fatal(err)
		}
	}

	return Team
}
