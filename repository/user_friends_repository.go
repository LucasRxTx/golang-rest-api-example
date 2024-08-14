package repository

import (
	"database/sql"
	"fmt"
	"rest-api/domain"
	"strings"
)

type IUserFriendsRepo interface {
	Update(tx *sql.Tx, userId string, friends []string) error
	Delete(tx *sql.Tx, userId string, friendId string) error
	GetAllForUserId(tx *sql.Tx, id string) ([]domain.UserFriend, error)
}

type UserFriendsRepository struct{}

func (repo *UserFriendsRepository) GetAllForUserId(tx *sql.Tx, id string) ([]domain.UserFriend, error) {
	rows, err := tx.Query(`
		SELECT
			f.id AS friend_id,
			f.name AS friend_name,
			s.score AS friend_score
		FROM
			game_user AS u
		JOIN game_friends AS jt
		ON
			u.id = jt.user_id
		JOIN game_user AS f
		ON
			jt.friend_id = f.id
		JOIN game_state AS s ON f.id = s.user_id
		WHERE u.id = $1;
		`,
		id,
	)

	if err != nil {
		return []domain.UserFriend{}, err
	}

	var friendsList []domain.UserFriend
	var friend domain.UserFriend
	for rows.Next() {
		err = rows.Scan(&friend.Id, &friend.Name, &friend.Highscore)
		if err != nil {
			return []domain.UserFriend{}, err
		}
		friendsList = append(friendsList, friend)
	}

	return friendsList, nil
}

func (repo *UserFriendsRepository) Delete(tx *sql.Tx, userId string, friendId string) error {
	query := `DELETE FROM game_friends WHERE user_id = $1 AND friend_id = $2;`
	fmt.Println("user_id", userId, "friend_id", friendId)
	_, err := tx.Exec(query, userId, friendId)

	return err
}

func (repo *UserFriendsRepository) Update(tx *sql.Tx, id string, friends []string) error {
	currentFriends, err := repo.GetAllForUserId(tx, id)
	fmt.Printf("Friends: %v\n", currentFriends)
	if err != nil {
		return err
	}

	friendsMap := map[string]int{}
	for _, newFriend := range friends {
		friendsMap[newFriend]++
	}

	for _, currentFriend := range currentFriends {
		friendsMap[currentFriend.Id.String()]--
	}

	var friendsToInsert []string
	var friendsToDelete []string
	for key, value := range friendsMap {
		if key == id { // cannot friend yourself
			continue
		}
		if value < 0 {
			friendsToDelete = append(friendsToDelete, key)
		} else if value > 0 {
			friendsToInsert = append(friendsToInsert, key)
		} else {
			continue
		}
	}

	if len(friendsToInsert) > 0 {
		query := `INSERT INTO game_friends (user_id, friend_id) VALUES `
		var inserts []string
		params := []any{id}
		for i, friendId := range friendsToInsert {
			// i + 2 since 0 is not valid in sql and $1 is the users id, all other will be friend id
			inserts = append(inserts, (fmt.Sprintf("($1, $%d)", i+2)))
			params = append(params, friendId)
		}

		query += strings.Join(inserts, ", ") + ";"
		_, err = tx.Exec(query, params...)
		if err != nil {
			return err
		}
	}

	if len(friendsToDelete) > 0 {
		for _, friendId := range friendsToDelete {
			err = repo.Delete(tx, id, friendId)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
