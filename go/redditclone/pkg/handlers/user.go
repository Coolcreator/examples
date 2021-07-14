package handlers

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"

	"myapp/pkg/session"
	"myapp/pkg/user"
)

type UserHandler struct {
	Tmpl     *template.Template
	Logger   *zap.SugaredLogger
	UserRepo *user.UserRepo
	Sessions *session.SessionsManager
}

type SecretInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *UserHandler) Index(w http.ResponseWriter, r *http.Request) {
	err := h.Tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Failed to execute template", http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userData := &SecretInfo{}
	err = json.Unmarshal(b, userData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u, err := h.UserRepo.Signup(userData.Username, userData.Password)
	if err == user.ErrUserExists {
		errInfo := []map[string]string{
			{
				"location": "body",
				"param":    "username",
				"value":    userData.Username,
				"msg":      "already exists",
			},
		}
		resp, err := json.Marshal(map[string]interface{}{
			"errors": errInfo,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(resp)
		return
	}

	sess, token, err := h.Sessions.Create(w, u.ID, u.Login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	tokenString, err := token.SignedString([]byte("ACCESS_KEY"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(map[string]interface{}{
		"token": tokenString,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(201)
	w.Write(resp)

	h.Logger.Infof("Created session for %v", sess.UserID)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userData := &SecretInfo{}
	err = json.Unmarshal(b, userData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u, err := h.UserRepo.Authorize(userData.Username, userData.Password)
	if err == user.ErrNoUser { // объединить ошибки
		resp, err := json.Marshal(map[string]interface{}{
			"message": err.Error(),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(resp)
		return
	}
	if err == user.ErrBadPass {
		resp, err := json.Marshal(map[string]interface{}{
			"message": err.Error(),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(resp)
		return
	}

	sess, token, err := h.Sessions.Create(w, u.ID, u.Login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tokenString, err := token.SignedString([]byte("secret token"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(map[string]interface{}{
		"token": tokenString,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(200)
	w.Write(resp)

	h.Logger.Infof("created session for %v", sess.UserID)
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	err := h.Sessions.DestroyCurrent(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
