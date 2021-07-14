package handlers

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/xid"
	"go.uber.org/zap"

	"myapp/pkg/items"
	"myapp/pkg/session"
)

type ItemsHandler struct {
	Tmpl      *template.Template
	ItemsRepo *items.ItemsRepo
	Logger    *zap.SugaredLogger
}

func (h *ItemsHandler) List(w http.ResponseWriter, r *http.Request) {
	elems, err := h.ItemsRepo.GetAll()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	resp, err := json.Marshal(elems)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(200)
	w.Write(resp)
}

func (h *ItemsHandler) AddPost(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	item := new(items.Item)
	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	item.Author.ID = sess.UserID
	item.Author.Username = sess.Login
	item.Comments = []items.Comment{}
	item.Created = time.Now()
	item.Score = 1
	item.UpvotePersentage = 100
	item.Votes = []items.Vote{}
	item.Votes = append(item.Votes, items.Vote{User: sess.UserID, Vote: 1})

	err = json.Unmarshal(b, item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	lastID, err := h.ItemsRepo.Add(item)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(200)
	w.Write(resp)

	h.Logger.Infof("Insert with id %v", lastID)
}

func (h *ItemsHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Bad ID", http.StatusBadGateway)
		return
	}

	item, err := h.ItemsRepo.GetByID(uint32(id))
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if item == nil {
		http.Error(w, "Invalid item request", http.StatusNotFound)
		return
	}

	item.Views += 1
	resp, err := json.Marshal(item)
	w.WriteHeader(200)
	w.Write(resp)

}

func (h *ItemsHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Bad ID", http.StatusBadGateway)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newComment := items.Comment{}
	commentBody := map[string]string{}
	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = json.Unmarshal(b, &commentBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newComment.Author.Username = sess.Login
	newComment.Author.ID = sess.UserID
	newComment.Body = commentBody["comment"]
	newComment.Created = time.Now()
	newComment.ID = xid.New().String()

	ok, item, err := h.ItemsRepo.Update(newComment, uint32(id))
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if !ok {
		http.Error(w, "Invalid item request", http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(200)
	w.Write(resp)
}

func (h *ItemsHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Bad ID", http.StatusBadGateway)
		return
	}

	ok, err := h.ItemsRepo.Delete(uint32(id))
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if !ok {
		http.Error(w, "Invalid item request", http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(map[string]bool{
		"success": ok,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(200)
	w.Write(resp)
}

func (h *ItemsHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Bad ID", http.StatusBadGateway)
		return
	}

	comment, ok := vars["comment"]
	if !ok {
		http.Error(w, "Bad comment ID", http.StatusBadGateway)
		return
	}

	post, err := h.ItemsRepo.Remove(uint32(id), comment)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(200)
	w.Write(resp)
}

func (h *ItemsHandler) UserList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username, ok := vars["username"]
	if !ok {
		http.Error(w, "Bad username", http.StatusBadGateway)
		return
	}

	posts, err := h.ItemsRepo.GetUserPosts(username)
	if err != nil {
		http.Error(w, "Database error", http.StatusBadGateway)
		return
	}

	resp, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(200)
	w.Write(resp)
}

func (h *ItemsHandler) CategoryList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	category, ok := vars["category"]
	if !ok {
		http.Error(w, "Bad category name", http.StatusBadGateway)
		return
	}

	posts, err := h.ItemsRepo.GetCategoryPosts(category)
	if err != nil {
		http.Error(w, "Database error", http.StatusBadGateway)
		return
	}

	resp, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(200)
	w.Write(resp)
}

func (h *ItemsHandler) Upvote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Bad ID", http.StatusBadGateway)
		return
	}

	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	post, err := h.ItemsRepo.Inc(sess.UserID, uint32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	post.UpvotePersentage = h.VotePercent(post)

	resp, err := json.Marshal(post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(200)
	w.Write(resp)
}

func (h *ItemsHandler) Downvote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Bad ID", http.StatusBadGateway)
		return
	}

	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	post, err := h.ItemsRepo.Dec(sess.UserID, uint32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	post.UpvotePersentage = h.VotePercent(post)

	resp, err := json.Marshal(post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(200)
	w.Write(resp)
}

func (h *ItemsHandler) Unvote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Bad ID", http.StatusBadGateway)
		return
	}

	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	post, err := h.ItemsRepo.CancelVote(sess.UserID, uint32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	post.UpvotePersentage = h.VotePercent(post)

	resp, err := json.Marshal(post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(200)
	w.Write(resp)
}

func (repo *ItemsHandler) VotePercent(item *items.Item) int {
	if len(item.Votes) == 0 || item.Score < 0 {
		return 0
	}
	return int(item.Score) / len(item.Votes) * 100
}
