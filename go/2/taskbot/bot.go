package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

var (
	BotToken   = "1442002945:AAHbN0D4uwG3gMbQTEfix3AUEVsYo6oCn0I"
	WebhookURL = "https://telegram-task-manager.herokuapp.com/"
)

type Task struct {
	ID       int
	Name     string
	Creator  int64
	Executor int64
}

type TaskList struct {
	tasks map[int]*Task
	mu    sync.Mutex
	nexID int
}

func NewTaskList() *TaskList {
	return &TaskList{
		tasks: map[int]*Task{},
		mu:    sync.Mutex{},
	}
}

func (tl *TaskList) GetTasks() ([]Task, error) {
	tl.mu.Lock()
	defer tl.mu.Unlock()

	if len(tl.tasks) == 0 {
		return []Task{}, fmt.Errorf("Нет задач")
	}

	keys := make([]int, 0, len(tl.tasks))
	for k := range tl.tasks {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	result := make([]Task, 0, len(tl.tasks))
	for _, k := range keys {
		result = append(result, *tl.tasks[k])
	}

	return result, nil
}

func (tl *TaskList) CreateTask(taskName string, taskCreator int64) (*Task, error) {
	tl.mu.Lock()
	defer tl.mu.Unlock()

	if len(taskName) == 0 {
		return &Task{}, fmt.Errorf("Невозможно создать задачу без имени")
	}

	for id, task := range tl.tasks {
		if task.Name == taskName {
			return &Task{}, fmt.Errorf("Такая задача существует в списке задач с id=%d", id)
		}
	}

	tl.nexID++
	newTask := Task{
		ID:       tl.nexID,
		Name:     taskName,
		Creator:  taskCreator,
		Executor: 0,
	}
	tl.tasks[newTask.ID] = &newTask

	return tl.tasks[tl.nexID], nil
}

func (tl *TaskList) AssignTask(taskID int, taskExecutor int64) (Task, error) {
	tl.mu.Lock()
	defer tl.mu.Unlock()

	if tl.tasks[taskID].Executor == taskExecutor {
		return Task{}, fmt.Errorf("Задача на вас")
	}
	oldTask := *tl.tasks[taskID]

	tl.tasks[taskID].Executor = taskExecutor
	return oldTask, nil
}

func (tl *TaskList) UnassignTask(taskID int, taskExecutor int64) (Task, error) {
	tl.mu.Lock()
	defer tl.mu.Unlock()

	if tl.tasks[taskID].Executor != taskExecutor {
		return Task{}, fmt.Errorf("Задача не на вас")
	}

	tl.tasks[taskID].Executor = 0
	return *tl.tasks[tl.nexID], nil
}

func (tl *TaskList) ResolveTask(taskID int, taskExecutor int64) (Task, error) {
	tl.mu.Lock()
	defer tl.mu.Unlock()

	oldTask := *tl.tasks[taskID]

	delete(tl.tasks, taskID)
	return oldTask, nil
}

func startTaskBot(ctx context.Context) error {
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Fatalf("NewBotAPI failed: %s", err)
	}

	bot.Debug = true
	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(WebhookURL))
	if err != nil {
		log.Fatalf("SetWebhook failed: %s", err)
	}

	updates := bot.ListenForWebhook("/")

	http.HandleFunc("/state", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("All is working"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	go func() {
		log.Fatalln("HTTP err:", http.ListenAndServe(":"+port, nil))
	}()

	go ExecuteCommand(bot, updates)
	fmt.Println("Start listen :" + port)
	return nil
}

func handleCommand(command string) (string, string) {
	arguments := strings.Split(command[1:], " ")
	if strings.Contains(arguments[0], "assign_") || strings.Contains(arguments[0], "unassign_") || strings.Contains(arguments[0], "resolve_") {
		arguments = strings.Split(arguments[0], "_")
	}
	command = arguments[0]
	parameters := strings.Join(arguments[1:], " ")
	return command, parameters
}

func ExecuteCommand(bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel) {

	tasks := NewTaskList()
	users := make(map[int64]string)
	for update := range updates {
		if _, ok := users[update.Message.Chat.ID]; !ok {
			users[update.Message.Chat.ID] = "@" + update.Message.Chat.UserName
		}

		if strings.HasPrefix(update.Message.Text, "/") {
			command, parameters := handleCommand(update.Message.Text)
			errorMessage := ""
			switch command {
			case "tasks", "my", "owner":
				at, err := tasks.GetTasks()
				if err != nil {
					errorMessage = err.Error()
					break
				}

				answer := ""
				if command == "tasks" {
					for _, task := range at {
						options := ""
						if task.Executor == 0 {
							options = "/assign_" + strconv.Itoa(task.ID)
						} else if task.Executor == update.Message.Chat.ID {
							options = "assignee: я\n" + "/unassign_" + strconv.Itoa(task.ID) + " /resolve_" + strconv.Itoa(task.ID)
						} else {
							options = "assignee: " + users[task.Executor]
						}
						answer += strconv.Itoa(task.ID) + ". " + task.Name + " by " + users[task.Creator] + "\n" + options + "\n\n"
					}
				} else if command == "my" {
					for _, task := range at {
						if task.Executor == update.Message.Chat.ID {
							answer += strconv.Itoa(task.ID) + ". " + task.Name + " by " + users[task.Creator] + "\n" +
								"/unassign_" + strconv.Itoa(task.ID) +
								" /resolve_" + strconv.Itoa(task.ID) + "\n\n"

						}
					}

					if answer == "" {
						bot.Send(tgbotapi.NewMessage(
							update.Message.Chat.ID,
							"На вас не назначено ни одной задачи",
						))
						break
					}

				} else if command == "owner" {
					for _, task := range at {
						if task.Creator == update.Message.Chat.ID {
							answer += strconv.Itoa(task.ID) + ". " + task.Name + " by " + users[task.Creator] + "\n"
							if task.Creator != task.Executor {
								answer += "/assign_" + strconv.Itoa(task.ID) + "\n\n"
							}
						}
					}

					if answer == "" {
						bot.Send(tgbotapi.NewMessage(
							update.Message.Chat.ID,
							"Вы не создавали задач",
						))
						break
					}
				}

				bot.Send(tgbotapi.NewMessage(
					update.Message.Chat.ID,
					answer[:len(answer)-2],
				))

			case "new":
				ct, err := tasks.CreateTask(parameters, update.Message.Chat.ID)
				if err != nil {
					errorMessage = err.Error()
					break
				}

				bot.Send(tgbotapi.NewMessage(
					update.Message.Chat.ID,
					"Задача "+"\""+ct.Name+"\" создана, id="+strconv.Itoa(ct.ID),
				))

			case "assign", "unassign", "resolve":
				t := Task{}
				n, err := strconv.Atoi(parameters)
				if err != nil {
					errorMessage = "Некорректный ввод id задачи для команды " + command
					break
				}

				if _, ok := tasks.tasks[n]; !ok {
					errorMessage = "Нет задачи с id=" + strconv.Itoa(n)
					break
				}

				if command == "assign" {
					t, err = tasks.AssignTask(n, update.Message.Chat.ID)
				} else if command == "unassign" {
					t, err = tasks.UnassignTask(n, update.Message.Chat.ID)
				} else if command == "resolve" {
					t, err = tasks.ResolveTask(n, update.Message.Chat.ID)
				}

				if err != nil {
					errorMessage = err.Error()
					break
				}

				if command == "assign" {
					bot.Send(tgbotapi.NewMessage(
						update.Message.Chat.ID,
						"Задача \""+t.Name+"\" назначена на вас",
					))

					if t.Executor == 0 {
						for id := range users {
							if id == t.Creator && id != update.Message.Chat.ID {
								bot.Send(tgbotapi.NewMessage(
									id,
									"Задача \""+t.Name+"\" назначена на "+users[update.Message.Chat.ID],
								))
							}
						}
					} else {
						bot.Send(tgbotapi.NewMessage(
							t.Executor,
							"Задача \""+t.Name+"\" назначена на "+users[update.Message.Chat.ID],
						))
					}
				} else if command == "unassign" {
					bot.Send(tgbotapi.NewMessage(
						update.Message.Chat.ID,
						"Принято",
					))

					for id := range users {
						if id == t.Creator && id != update.Message.Chat.ID {
							bot.Send(tgbotapi.NewMessage(
								id,
								"Задача \""+t.Name+"\" осталась без исполнителя",
							))
							break
						}
					}
				} else if command == "resolve" {
					bot.Send(tgbotapi.NewMessage(
						update.Message.Chat.ID,
						"Задача \""+t.Name+"\" выполнена",
					))

					for id := range users {
						if id == t.Creator && t.Creator != update.Message.Chat.ID {
							bot.Send(tgbotapi.NewMessage(
								id,
								"Задача \""+t.Name+"\" выполнена "+users[update.Message.Chat.ID],
							))
							break
						}
					}
				}
			case "start":
				bot.Send(tgbotapi.NewMessage(
					update.Message.Chat.ID,
					"Привет! Я - простой менеджер по управлению задачами\n"+
						"Если хочешь узнать доступные тебе команды, то набери\n/help",
				))
			case "help":
				bot.Send(tgbotapi.NewMessage(
					update.Message.Chat.ID,
					"Доступные команды:\n\n"+
						"/tasks - выводит все активные задачи\n"+
						"/new XXX YYY ZZZ - создаёт новую задачу\n"+
						"/assign_$ID - делает пользователя исполнителем задачи\n"+
						"/unassign_$ID - снимает задачу с текущего исполнителя\n"+
						"/resolve_$ID - выполняет задачу, удаляет её из списка\n"+
						"/my - показывает задачи, которые назначены на меня\n"+
						"/owner - показывает задачи которые были созданы мной",
				))
			default:
				bot.Send(tgbotapi.NewMessage(
					update.Message.Chat.ID,
					"Некорректная команда. Список всех доступных команд лежит в /help",
				))
			}

			if errorMessage != "" {
				bot.Send(tgbotapi.NewMessage(
					update.Message.Chat.ID,
					errorMessage,
				))
				errorMessage = ""
			}

		} else {
			bot.Send(tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"Не команда",
			))
		}
	}
}

func main() {
	err := startTaskBot(context.Background())
	if err != nil {
		panic(err)
	}
}
