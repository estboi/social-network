package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"social-network/core"
	database "social-network/database/methods"
	"social-network/server"
	"social-network/server/websocket"
	"social-network/sessions"

	routing "github.com/Zewasik/go-router"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
		return
	}
	core := core.NewCore(db)
	manager := websocket.NewManager(*core)
	httpadpt := server.NewHttpAdapter(core, manager)

	SetupRoutes(manager, httpadpt)

	fmt.Printf("Starting server at port 8080\n")
	fmt.Printf("http://localhost:8080/\n")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func SetupRoutes(manager *websocket.Manager, httpadpt *server.HttpAdapter) {
	r := routing.NewRouterBuilder().
		SetAllowOrigin("http://localhost:3000").
		SetAllowMethods([]string{"POST", "GET", "OPTIONS", "PUT", "DELETE"}).
		SetAllowHeaders([]string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token"}).
		SetCredantials(true)

	// WS SERVING
	http.HandleFunc("/ws", manager.ServeWS)

	r.NewRoute("get", `/api/auth`, server.CheckIfAuth)
	// AUTH ROUTERS
	r.NewRoute("post", `/api/register`, httpadpt.RegisterHandlerPost)
	r.NewRoute("post", `/api/login`, httpadpt.LoginHandlerPost)
	r.NewRoute("post", `/api/logout`, sessions.CloseSession)

	// IMAGE ROUTERS
	r.NewRoute("get", `/api/image/user/(?P<id>\d+)`, httpadpt.ImageUserHandler, middleware())
	r.NewRoute("get", `/api/image/post/(?P<id>\d+)`, httpadpt.ImagePostHandler, middleware())
	r.NewRoute("get", `/api/image/group/(?P<id>\d+)`, httpadpt.ImageGroupHandler, middleware())
	// IMAGE CREATION ROUTERS
	r.NewRoute("post", `/api/image/user`, httpadpt.ImageUserCreateHandler)
	r.NewRoute("post", `/api/image/post/(?P<id>\d+)`, httpadpt.ImagePostCreateHandler, middleware())
	r.NewRoute("post", `/api/image/group/(?P<id>\d+)`, httpadpt.ImageGroupCreateHandler, middleware())

	// NAVBAR ROUTERS
	r.NewRoute("get", `/api/navbar`, httpadpt.NavbarHandlerGet, middleware())
	r.NewRoute("get", `/api/notification`, httpadpt.NotificationsHandlerGet, middleware())

	// POSTS ROUTERS
	r.NewRoute("get", `/api/posts/(?P<id>\d+)`, httpadpt.PostsHandlerGet, middleware())
	r.NewRoute("get", `/api/posts/home`, httpadpt.PostsHandlerGetHome, middleware())
	r.NewRoute("get", `/api/posts/group/(?P<id>\d+)`, httpadpt.PostsHandlerGetGroup, middleware())
	r.NewRoute("get", `/api/posts/users/(?P<id>\d+)`, httpadpt.PostsHandlerGetUser, middleware())
	// POSTS CREATION ROUTERS
	r.NewRoute("post", `/api/posts/create/(?P<id>\d+)`, httpadpt.PostsHandlerPostPost, middleware())
	// POSTS RELATED ACTIONS ROUTERS
	r.NewRoute("post", `/api/posts/vote/(?P<id>\d+)`, httpadpt.PostActionVoteHandler, middleware())

	// COMMENST ROUTERS
	r.NewRoute("get", `/api/posts/comments/(?P<id>\d+)`, httpadpt.CommentsGetHandler, middleware())
	// COMMENTS CREATION ROUTERS
	r.NewRoute("post", `/api/posts/comments/create/(?P<id>\d+)`, httpadpt.CommentsPostHandler, middleware())
	// COMMENTS IMAGE ROUTERS
	r.NewRoute("get", `/api/image/comment/(?P<id>\d+)`, httpadpt.CommentsImageGetHandler, middleware())
	r.NewRoute("post", `/api/image/comment/(?P<id>\d+)`, httpadpt.CommentsImagePostHandler, middleware())
	// COMMENTS RELATED ACTIONS ROUTERS
	r.NewRoute("post", `/api/posts/comment/vote/(?P<id>\d+)`, httpadpt.CommentsLikeHandler, middleware())

	// USERS ROUTERS
	r.NewRoute("get", `/api/users/all`, httpadpt.UsersHandlerGetAll, middleware())
	r.NewRoute("get", `/api/users/followers`, httpadpt.UsersHandlerGetFollowers, middleware())
	r.NewRoute("get", `/api/users/followed`, httpadpt.UsersHandlerGetFollowed, middleware())
	r.NewRoute("post", `/api/users/modify`, httpadpt.UsersHandlerModify, middleware())
	r.NewRoute("get", `/api/users/profile/(?P<id>\d+)`, httpadpt.UsersHandlerGetProfile, middleware())
	// USERS RELATED ACTIONS ROUTERS
	r.NewRoute("post", `/api/users/subscribe/(?P<id>\d+)`, httpadpt.ActionHandlerPostSubscribeOnUser, middleware())
	r.NewRoute("post", `/api/users/unsubscribe/(?P<id>\d+)`, httpadpt.ActionHandlerPostUnsubscribeOnUser, middleware())

	// GROUPS ROUTERS
	r.NewRoute("get", `/api/groups/all`, httpadpt.GroupsHandlerGetAll, middleware())
	r.NewRoute("get", `/api/groups/connected`, httpadpt.GroupsHandlerGetConnected, middleware())
	r.NewRoute("get", `/api/groups/created`, httpadpt.GroupsHandlerGetCreated, middleware())
	r.NewRoute("get", `/api/groups/profile/(?P<id>\d+)`, httpadpt.GroupsHandlerGetProfile, middleware())
	r.NewRoute("get", `/api/groups/notmembers/(?P<id>\d+)`, httpadpt.GroupsHandlerGetNotMembers, middleware())
	r.NewRoute("get", `/api/groups/requested/(?P<id>\d+)`, httpadpt.GroupsHandlerGetRequested, middleware())
	// GROUPS CREATION ROUTERS
	r.NewRoute("post", `/api/groups/create`, httpadpt.GroupsHandlerPost, middleware())
	// GROUPS RELATED ACTIONS ROUTERS
	r.NewRoute("post", `/api/groups/request/(?P<id>\d+)`, httpadpt.GroupsRequestHandler, middleware())
	r.NewRoute("post", `/api/groups/invite/(?P<id>\d+)`, httpadpt.GroupsInviteHandler, middleware())
	r.NewRoute("post", `/api/groups/accept/(?P<id>\d+)`, httpadpt.GroupsAcceptHandler, middleware())
	r.NewRoute("post", `/api/groups/deny/(?P<id>\d+)`, httpadpt.GroupsDenyHandler, middleware())
	r.NewRoute("post", `/api/groups/inviteAccept/(?P<id>\d+)`, httpadpt.GroupsAcceptInviteHandler, middleware())
	r.NewRoute("post", `/api/groups/inviteDeny/(?P<id>\d+)`, httpadpt.GroupsInviteDenyHandler, middleware())

	// EVENTS ROUTERS
	r.NewRoute("get", `/api/events/all`, httpadpt.EventHandlerGetAll, middleware())
	r.NewRoute("get", `/api/events/group/(?P<id>\d+)`, httpadpt.EventHandlerGetGroup, middleware())
	// EVENTS CREATION ROUTERS
	r.NewRoute("post", `/api/events/create/(?P<id>\d+)`, httpadpt.EventHandlerPost, middleware())
	// GROUPS RELATED ACTIONS ROUTERS
	r.NewRoute("post", `/api/events/attend/(?P<id>\d+)`, httpadpt.ActionHandlerPostJoinEvent, middleware())
	r.NewRoute("post", `/api/events/deny/(?P<id>\d+)`, httpadpt.ActionHandlerPostDenyEvent, middleware())

	// CHATS ROUTERS
	r.NewRoute("get", `/api/chats/all`, httpadpt.ChatsHandlerGetAll, middleware())
	r.NewRoute("get", `/api/chats/group/(?P<id>\d+)`, httpadpt.ChatsHandlerGetGroup, middleware())
	r.NewRoute("get", `/api/chats/user/(?P<id>\d+)`, httpadpt.ChatsHandlerGetOneChat, middleware())

	http.HandleFunc("/", r.ServeWithCORS())
}

// Checks for session validity. If user not logged should return error redirect.
func middleware() routing.Middleware {
	return func(main http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id, err := sessions.Validate(r)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), routing.ContextKey("user"), id)
			r = r.WithContext(ctx)

			main.ServeHTTP(w, r)
		})
	}
}
