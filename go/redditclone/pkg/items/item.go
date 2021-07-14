package items

import "time"

type AuthorInfo struct {
	ID       uint32 `json:"id"`
	Username string `json:"username"`
}

type Comment struct {
	ID      string     `json:"id"`
	Author  AuthorInfo `json:"author"`
	Body    string     `json:"body"`
	Created time.Time  `json:"created"`
}

type Vote struct {
	User uint32 `json:"user"`
	Vote int32  `json:"vote"`
}

type Item struct {
	Author           AuthorInfo `json:"author"`
	Category         string     `json:"category"`
	Comments         []Comment  `json:"comments"`
	Created          time.Time  `json:"created"`
	ID               uint32     `json:"id"`
	Score            int        `json:"score"`
	Text             string     `json:"text"`
	Title            string     `json:"title"`
	Type             string     `json:"type"`
	UpvotePersentage int        `json:"upvotePercentage"`
	Url              string     `json:"url"`
	Views            uint32     `json:"views"`
	Votes            []Vote     `json:"votes"`
}
