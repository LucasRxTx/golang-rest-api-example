package dto

type UserFriendDto struct {
	Id        string `json:"id" binding:"required,uuid"`
	Name      string `json:"name" binding:"required"`
	Highscore int    `json:"highscore" binding:"required"`
}

type UserFriendsListDto struct {
	Freinds []UserFriendDto `json:"friends" binding:"required"`
}

func (fl *UserFriendsListDto) AddFriend(friend UserFriendDto) {
	fl.Freinds = append(fl.Freinds, friend)
}