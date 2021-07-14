package items

import (
	"errors"
	"sync"
)

var (
	ErrNoItem = errors.New("the item not found")
)

type ItemsRepo struct {
	lastID uint32
	data   []*Item
	mu     sync.Mutex
}

func NewRepo() *ItemsRepo {
	return &ItemsRepo{
		data: make([]*Item, 0, 10),
		mu:   sync.Mutex{},
	}
}

func (repo *ItemsRepo) GetAll() ([]*Item, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	return repo.data, nil
}

func (repo *ItemsRepo) Add(item *Item) (uint32, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.lastID++
	item.ID = repo.lastID
	repo.data = append(repo.data, item)

	return repo.lastID, nil
}

func (repo *ItemsRepo) GetByID(id uint32) (*Item, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	for _, item := range repo.data {
		if item.ID == id {
			return item, nil
		}
	}

	return nil, nil
}

func (repo *ItemsRepo) Update(comment Comment, id uint32) (bool, *Item, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	for _, item := range repo.data {
		if item.ID != id {
			continue
		}
		item.Comments = append(item.Comments, comment)
		return true, item, nil
	}

	return false, &Item{}, nil
}

func (repo *ItemsRepo) Delete(id uint32) (bool, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	i := -1
	for idx, item := range repo.data {
		if item.ID != id {
			continue
		}
		i = idx
	}
	if i < 0 {
		return false, nil
	}

	if i < len(repo.data)-1 {
		copy(repo.data[i:], repo.data[i+1:])
	}
	repo.data[len(repo.data)-1] = nil
	repo.data = repo.data[:len(repo.data)-1]

	return true, nil
}

func (repo *ItemsRepo) Remove(id uint32, comment string) (*Item, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	i := -1
	for idx, item := range repo.data {
		if item.ID != id {
			continue
		}
		i = idx
	}

	post := repo.data[i]
	for idx, comm := range post.Comments {
		if comm.ID != comment {
			continue
		}
		i = idx
	}

	if i < 0 {
		return &Item{}, nil
	}

	if i < len(post.Comments)-1 {
		copy(post.Comments[i:], post.Comments[i+1:])
	}
	post.Comments[len(post.Comments)-1] = Comment{}
	post.Comments = post.Comments[:len(post.Comments)-1]

	return post, nil
}

func (repo *ItemsRepo) GetUserPosts(username string) ([]*Item, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	result := []*Item{}
	for _, item := range repo.data {
		if item.Author.Username == username {
			result = append(result, item)
		}
	}

	return result, nil
}

func (repo *ItemsRepo) GetCategoryPosts(category string) ([]*Item, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	result := []*Item{}
	for _, item := range repo.data {
		if item.Category == category {
			result = append(result, item)
		}
	}

	return result, nil
}

func (repo *ItemsRepo) Inc(username uint32, id uint32) (*Item, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	item, _ := repo.GetByID(id)
	if item != nil {
		for v, vote := range item.Votes {
			if username == vote.User {
				if vote.Vote == -1 {
					item.Votes[v].Vote = 1
					item.Score += 2
				}
				return item, nil
			}
		}
		item.Votes = append(item.Votes, Vote{User: username, Vote: 1})
		item.Score += 1
		return item, nil

	}
	return nil, ErrNoItem
}

func (repo *ItemsRepo) Dec(username uint32, id uint32) (*Item, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	item, _ := repo.GetByID(id)
	if item != nil {
		for v, vote := range item.Votes {
			if username == vote.User {
				if vote.Vote == 1 {
					item.Votes[v].Vote = -1
					item.Score -= 2
				}
				return item, nil
			}
		}
		item.Votes = append(item.Votes, Vote{User: username, Vote: -1})
		item.Score -= 1
		return item, nil
	}
	return nil, ErrNoItem
}

func (repo *ItemsRepo) CancelVote(username uint32, id uint32) (*Item, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	item, _ := repo.GetByID(id)
	if item != nil {
		for i, vote := range item.Votes {
			if username == vote.User {
				if vote.Vote == -1 {
					item.Score += 1
				} else if vote.Vote == 1 {
					item.Score -= 1
				}
				//delete vote
				if i < len(item.Votes)-1 {
					copy(item.Votes[i:], item.Votes[i+1:])
				}
				item.Votes[len(item.Votes)-1] = Vote{}
				item.Votes = item.Votes[:len(repo.data)-1]
				return item, nil
			}
		}

	}
	return item, ErrNoItem
}
