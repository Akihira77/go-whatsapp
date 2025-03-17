package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"sort"
	"time"

	"github.com/Akihira77/go_whatsapp/src/components"
	"github.com/Akihira77/go_whatsapp/src/services"
	"github.com/Akihira77/go_whatsapp/src/types"
	"github.com/Akihira77/go_whatsapp/src/utils"
	"github.com/Akihira77/go_whatsapp/src/views"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PageHandler struct {
	userService *services.UserService
	chatService *services.ChatService
	v           *utils.MyValidator
}

func NewPageHandler(userService *services.UserService, chatService *services.ChatService) *PageHandler {
	return &PageHandler{
		userService: userService,
		chatService: chatService,
		v:           utils.NewMyValidator(),
	}
}

func (ph *PageHandler) RenderSignup(c *gin.Context) {
	views.Signup().Render(c, c.Writer)
}

func (ph *PageHandler) RenderSignin(c *gin.Context) {
	views.Signin().Render(c, c.Writer)
}

func (ph *PageHandler) RenderHome(c *gin.Context) {
	user, ok := c.MustGet("user").(*types.User)
	if !ok {
		slog.Error("Retrieving user's info")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed retrieving user's info"})
		return
	}

	userName := c.Query("username")
	groupName := c.Query("groupname")

	if userName != "" {
		userName = "%" + userName + "%"
	}

	if groupName != "" {
		groupName = "%" + groupName + "%"
	}

	if groupName == "" && userName == "" {
		userName = "%" + userName + "%"
		groupName = "%" + groupName + "%"
	}

	chatList, err := ph.chatService.SearchChat(c, user.ID, userName, groupName)
	if err != nil {
		slog.Error("Retrieving last messages",
			"error", err,
		)
	}

	if c.GetHeader("X-Page-Query") != "" {
		components.ChatList(chatList).Render(c, c.Writer)
		return
	} else if c.GetHeader("X-From-Group") != "" {
		components.HomeSidebar(user, chatList).Render(c, c.Writer)
		return
	}

	user.Status = types.ONLINE
	views.Home(chatList, nil).Render(c, c.Writer)
}

func (ph *PageHandler) RenderChatPage(c *gin.Context) {
	user, ok := c.MustGet("user").(*types.User)
	if !ok {
		slog.Error("Retrieving user's info")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed retrieving your user's info"})
		return
	}

	senderId := c.Query("userId")
	if senderId != "" {
		u, err := ph.userService.FindUserByID(c, senderId)
		if err != nil {
			slog.Error("Retrieving sender's info",
				"error", err,
			)
			return
		}

		msgs, err := ph.chatService.MarkMessagesAsRead(c, senderId, &user.ID, nil)
		if err != nil {
			slog.Error("Retrieving reads messages",
				"error", err,
			)
			return
		}

		components.ChatPage(user, u, msgs).Render(c, c.Writer)
	} else {
		groupId := c.Query("groupId")
		g, err := ph.userService.FindGroupByID(c, groupId)
		if err != nil {
			slog.Error("Retrieving group's info",
				"error", err,
			)
			return
		}

		msgs, err := ph.chatService.MarkMessagesAsRead(c, senderId, nil, &groupId)
		if err != nil {
			slog.Error("Retrieving reads messages",
				"error", err,
			)
			return
		}
		g.Messages = msgs

		sort.Slice(g.Member, func(i, j int) bool {
			memberI := utils.GetFullName(&g.Member[i].User)
			memberJ := utils.GetFullName(&g.Member[j].User)

			return memberI <= memberJ
		})
		components.GroupPage(user, g).Render(c, c.Writer)
	}
}

func (ph *PageHandler) RenderMyProfile(c *gin.Context) {
	c.Header("HX-Redirect", "/users/profile")

	views.MyProfile().Render(c, c.Writer)
}

func (ph *PageHandler) RenderEditProfile(c *gin.Context) {
	views.EditUser().Render(c, c.Writer)
}

func (ph *PageHandler) RenderChangePassword(c *gin.Context) {
	views.ChangePassword().Render(c, c.Writer)
}

func (ph *PageHandler) RenderMyContacts(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 500*time.Millisecond)
	defer cancel()

	user, ok := c.MustGet("user").(*types.User)
	if !ok {
		slog.Error("Failed retrieve user's data from context")
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Failed retrieving your user info"})
		return
	}

	search := c.Query("search")
	users, err := ph.userService.GetMyContacts(ctx, user.ID, search)
	if err != nil {
		slog.Error("Failed retrieve user's contacts",
			"error", err,
		)
	}

	views.MyContacts(user, users).Render(c, c.Writer)
}

func (ph *PageHandler) RenderUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 500*time.Millisecond)
	defer cancel()

	user, ok := c.MustGet("user").(*types.User)
	if !ok {
		slog.Error("Failed retrieve user's data from context")
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Failed retrieving your user info"})
		return
	}

	var query types.UserQuerySearch
	if err := c.ShouldBindQuery(&query); err != nil {
		slog.Error("Failed extract query",
			"error", err,
		)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed retrieving query search"})
		return
	}

	if query.Size == 0 {
		query.Size = 10
	}
	slog.Info("Paginate",
		"query", query,
	)

	users, err := ph.userService.GetUsers(ctx, user, &query)
	if err != nil {
		slog.Error("Failed retrieve user's contacts",
			"error", err,
		)
	}

	if c.GetHeader("X-Page-Query") != "" {
		components.UserList(users, &query, false).Render(c, c.Writer)
		return
	}

	views.Users(users, &query).Render(c, c.Writer)
}

func (ph *PageHandler) RenderMakeGroup(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 500*time.Millisecond)
	defer cancel()

	user, ok := c.MustGet("user").(*types.User)
	if !ok {
		slog.Error("Failed retrieve user's data from context")
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Failed retrieving your user info"})
		return
	}

	var query types.UserQuerySearch
	if err := c.ShouldBindQuery(&query); err != nil {
		slog.Error("Failed extract query",
			"error", err,
		)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed retrieving query search"})
		return
	}

	if query.Size == 0 {
		query.Size = 10
	}

	users, err := ph.userService.GetUsers(ctx, user, &query)
	if err != nil {
		slog.Error("Failed retrieve user's contacts",
			"error", err,
		)
	}

	if c.GetHeader("X-Page-Query") != "" {
		components.UserList(users, &query, true).Render(c, c.Writer)
		return
	}

	views.MakeGroup(users, &query).Render(c, c.Writer)
}

func (ph *PageHandler) RenderGroupInfo(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 500*time.Millisecond)
	defer cancel()

	groupId := c.Param("groupId")
	g, err := ph.userService.FindGroupByID(ctx, groupId)
	if err != nil {
		slog.Error("Failed retrieve group info",
			"groupId", groupId,
			"err", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Group not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed retrieve group info"})
		return
	}

	components.GroupInfo(g).Render(c, c.Writer)
}

func (ph *PageHandler) RenderNamingGroup(c *gin.Context) {
	components.NamingGroup().Render(c, c.Writer)
}

func (ph *PageHandler) CreateGroup(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 500*time.Millisecond)
	defer cancel()

	user, ok := c.MustGet("user").(*types.User)
	if !ok {
		slog.Error("Failed retrieve user's data from context")
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Failed retrieving your user info"})
		return
	}

	var data types.CreateGroup
	if err := c.ShouldBindJSON(&data); err != nil {
		slog.Error("Failed extract request payload",
			"err", err,
		)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Request payload is invalid"})
		return
	}

	errs := ph.v.Validate(&data)
	if errs != nil || len(errs) > 0 {
		slog.Error("Request payload is invalid",
			"error", errs,
		)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Create group failed"})
		return
	}

	var member []string
	err := json.Unmarshal([]byte(data.Member), &member)
	if err != nil {
		slog.Error("Failed extract member of group",
			"error", err,
		)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed extract member of group"})
		return
	}

	imageData, err := base64.StdEncoding.DecodeString(data.GroupProfile)
	if data.GroupProfile != "" && err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image data"})
		return
	}

	data.Creator = user

	member = append(member, user.ID)
	group, err := ph.userService.CreateGroup(ctx, data, imageData, member)
	if err != nil {
		slog.Error("Failed creating group",
			"error", err,
		)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed saving group information"})
		return
	}

	group.Messages, _ = ph.chatService.GetMessagesInsideGroup(ctx, group.ID)
	chatList, err := ph.chatService.SearchChat(c, user.ID, "%%", "%%")
	if err != nil {
		slog.Error("Retrieving last messages",
			"error", err,
		)
	}

	// c.Header("HX-Redirect", "/")
	views.Home(chatList, group).Render(c, c.Writer)
}

func (ph *PageHandler) ExitGroup(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 500*time.Millisecond)
	defer cancel()

	user, ok := c.MustGet("user").(*types.User)
	if !ok {
		slog.Error("Failed retrieve user's data from context")
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Failed retrieving your user info"})
		return
	}

	groupId := c.Param("groupId")
	group, err := ph.userService.FindGroupByID(ctx, groupId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Group does not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed searching group information"})
		return
	}

	ph.userService.ExitGroup(ctx, user.ID, group)
}
