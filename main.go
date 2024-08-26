package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	goaway "github.com/TwiN/go-away"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/vmihailenco/msgpack"
)

type Vector3 struct {
	X float64 `msgpack:"x"`
	Y float64 `msgpack:"y"`
	Z float64 `msgpack:"z"`
}

// /////////
// todo: fiddle with broadcast channel pointers, ensure async correct, implement sim, impelemtn racetrack
var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  4096 * 2,
		WriteBufferSize: 4096 * 2,

		CheckOrigin: func(r *http.Request) bool {

			return true
		},
	}

	morphMap    = genRefMap("base.json")
	morphMax    = 100
	morphBuffer = []*Kaizomorph{}

	users            = make(map[string]int)
	clients          = make(map[*websocket.Conn]Player)
	broadcastChannel = make(chan Player, 1000)
	leaderboard      = make(map[string]interface{})

	morphChannel   = make(chan MorphMessage, 1000)
	keyChannel     = make(chan Message)
	spawnChannel   = make(chan MorphMessage)
	clientsMutex   sync.Mutex
	broadcastMutex sync.Mutex
	keysMutex      sync.Mutex
	spawnMutex     sync.Mutex
	morphMutex     sync.Mutex
)

const PLAYERY = 108

type Message struct {
	Type       string  `json:"type"`
	Pos        Vector3 `json:"pos,omitempty"`
	RotY       float64 `json:"roty,omitempty"`
	Key        Msg     `json:"key,omitempty"`
	Beam       bool    `json:"beam,omitempty"`
	Positional bool    `json:"positional,omitempty"`
}
type MorphMessage struct {
	Type  string     `json:"type"`
	Pos   Vector3    `json:"pos,omitempty"`
	Key   Msg        `json:"key,omitempty"`
	Morph Kaizomorph `json:"Kaizomorph,omitempty"`
	UUID  string     `json:"UUID"`
}
type Destroy struct {
	Type string `json:"type"`
	UUID string `json:"UUID"`
}
type Join struct {
	Type        string         `json:"type"`
	MorphBuffer []MorphMessage `json:"morphbuffer"`
	PlayerList  []Player       `json:"players"`
}

type Msg struct {
	Key    string `msgpack:"key"`
	Action string `msgpack:"action"`
	Type   string `msgpack:"type"`
}

type Player struct {
	Type        string    `msgpack:"type"`
	Pos         *Vector3  `msgpack:"Pos"`
	RotY        float64   `msgpack:"RotY"`
	Vel         float64   `msgpack:"vel"`
	ID          string    `msgpack:"id"`
	Timestamp   time.Time `msgpack:"timestamp"`
	SizeX       float64   `msgpack:"sizeX"`
	SizeZ       float64   `msgpack:"sizeZ"`
	Colliding   bool      `msgpack:"colliding"`
	Cooldown    float64   `msgpack:"Cooldown"`
	Rares       int       `msgpack:"rares"`
	TotalPoints int       `msgpack:"total"`
	Name        string    `msgpack:"name"`
}

// playerMsg, just a lightweight version of player

type Collision struct {
	ObjectID int
	Distance float64
}

func wserve(w http.ResponseWriter, r *http.Request, username string) {

	// Perform username validation before upgrading to WebSocket
	//verified, err := validUsername(string(username))
	//if err != nil {
	//	w.WriteHeader(http.StatusBadRequest)
	//	fmt.Fprint(w, err)
	//	return
	//}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	verified := true
	if verified {

		clientsMutex.Lock()
		ID := len(clients)

		users[username] = ID

		player := Player{
			ID:          fmt.Sprintf("%d", ID),
			Pos:         &Vector3{X: 0, Y: 0, Z: 0},
			RotY:        0,
			Vel:         0,
			Timestamp:   time.Now(),
			Type:        "player",
			Name:        fmt.Sprintf("Player%d", ID),
			SizeX:       32,
			SizeZ:       32,
			Colliding:   false,
			TotalPoints: 0,
		}
		clients[conn] = player // = chan

		playerList := []Player{}
		morphList := []MorphMessage{}
		for _, p := range clients {
			playerList = append(playerList, p)
		}
		for _, m := range morphBuffer {
			msg := MorphMessage{}
			msg.Morph = *m
			msg.UUID = m.UUID
			msg.Pos = *m.Pos
			morphList = append(morphList, msg)
		}
		joinMsg := Join{
			Type:        "join",
			MorphBuffer: morphList,
			PlayerList:  playerList,
		}
		conn.WriteJSON(joinMsg)

		clientsMutex.Unlock()
		morphPacket := MorphMessage{Type: "beam"}
		destroyPacket := Destroy{Type: "destroy"}

		// options to unfuck this:
		// - one large message / game state to everyone
		// - send client side messages back to client (no need for locking mutex)
		// - refalsify beam
		for {

			var msg Message

			_, message, _ := conn.ReadMessage()

			err := msgpack.Unmarshal(message, &msg)
			if err != nil {
				fmt.Println("NO MESSAGE")
			}

			if msg.Positional {
				player.Pos = &msg.Pos
				player.RotY = msg.RotY
				player.Timestamp = time.Now()
			}

			for i := 0; i < len(morphBuffer); i++ {
				if handleCollision(&player, morphBuffer[i], i) {
					player.Colliding = true
					fmt.Println("COLLIDING: ", i)
					if msg.Beam {
						// send beam = true client side for all players but not player
						/*
						   						for _, p := range clients {
						                               if p.ID!= player.ID {
						                                   p.Colliding = true
						                               }
						                           }
						*/
						fmt.Println("BEAMED")
						morphBuffer[i].Pos = morphBuffer[i].Pos.Lerp(&Vector3{X: morphBuffer[i].Pos.X, Y: 108, Z: morphBuffer[i].Pos.Z}, 0.05)
						morphPacket.Pos = *morphBuffer[i].Pos
						morphPacket.UUID = morphBuffer[i].UUID
						if morphBuffer[i].Pos.Y > PLAYERY-12 {

							// morphMap.remove(UUID)

							//morphBuffer[i] = nil
							//2. if morphBuffer[i] == nil,

							// send morphbuffer update to sync
							player.TotalPoints++
							fmt.Println("POINT GET: ", player.TotalPoints)
							leaderboard[player.Name] = player.TotalPoints
							destroyPacket.UUID = morphBuffer[i].UUID

							// AVAILABLE SPOT CHANNEL <- INDEX

							clientsMutex.Lock()

							for client := range clients {
								morpherror := client.WriteJSON(leaderboard)
								if morpherror != nil {
									fmt.Println(morpherror)
								}
								morpherror2 := client.WriteJSON(destroyPacket)
								if morpherror2 != nil {
									fmt.Println(morpherror)
								}

							}

							clientsMutex.Unlock()
							morphBuffer = RemoveIndex(morphBuffer, i)
							//break
						}
						// send destroy msg
						morphChannel <- morphPacket
					}
				} else {
					player.Colliding = false
					if morphBuffer[i].Pos.Y > 26 && msg.Beam {
						morphBuffer[i].Pos = morphBuffer[i].Pos.Lerp(&Vector3{X: morphBuffer[i].Pos.X, Y: 24, Z: morphBuffer[i].Pos.Z}, 0.06)
						morphPacket.Pos = *morphBuffer[i].Pos
						morphPacket.UUID = morphBuffer[i].UUID
						morphChannel <- morphPacket
					}

				}

				if morphBuffer[i].Pos.Y > 26 && !msg.Beam {
					morphBuffer[i].Pos = morphBuffer[i].Pos.Lerp(&Vector3{X: morphBuffer[i].Pos.X, Y: 24, Z: morphBuffer[i].Pos.Z}, 0.06)
					morphPacket.Pos = *morphBuffer[i].Pos
					morphPacket.UUID = morphBuffer[i].UUID
					morphChannel <- morphPacket
				}

			}
			/*
				keysMutex.Lock()
				if msg.Type == "key" || msg.Type == "beam" {
					select {
					case keyChannel <- msg:
					default:
						fmt.Println("keyChannel is full, skipping")
					}
				}
				keysMutex.Unlock()
			*/
			//broadcastMutex.Lock()
			//select {
			//case broadcastChannel <- player:
			//default:
			//	fmt.Println("broadcastChannel is unused")
			//}

			// when player stops moving, this becomes congested and prevents beam msgs from being sent.
			// or, there's no msg to be sent
			broadcastChannel <- player

			//broadcastMutex.Unlock()

		}

	} else {
		fmt.Println("Invalid username:", username)
		conn.WriteJSON("Inavlid username") // wouldnt it be nice to have something that can writejson, println, and log errors all in one?
		conn.Close()
	}

}

func validUsername(u string) (bool, error) {

	if goaway.IsProfane(u) {

		return false, Errlog{p1: "please pick politely"}

	}

	if len(u) < 3 || len(u) > 32 {

		return false, Errlog{p1: "3-32 characters please"}

	}
	//_, ok := users[u]
	//if ok {
	//	return false, Errlog{p1: "User exists."}
	//}

	return true, nil
}

func handleCollision(player *Player, obj *Kaizomorph, id int) bool {

	// only needs to be calculated once, do this later
	objMinX, objMinZ := obj.Pos.X-obj.SizeX/2, obj.Pos.Z-obj.SizeZ/2
	objMaxX, objMaxZ := obj.Pos.X+obj.SizeX/2, obj.Pos.Z+obj.SizeZ/2

	playerMinX, playerMinZ := player.Pos.X-player.SizeX/2, player.Pos.Z-player.SizeZ/2
	playerMaxX, playerMaxZ := player.Pos.X+player.SizeX/2, player.Pos.Z+player.SizeZ/2

	overlap := (playerMaxX >= objMinX && playerMinX <= objMaxX) && (playerMaxZ >= objMinZ && playerMinZ <= objMaxZ)
	if overlap {
		return true
	}
	return false

}

func spawner() {

	// fill up a queue and only send morphs when a player connects (might already work like that)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if len(morphBuffer) < morphMax {
				morphPacket := MorphMessage{}
				morphPos := randomPointSquare(1800, 1800, 900, 900)
				morph := genMorph(morphMap, plainsWeights, rand.Intn(50))
				morph.Pos = morphPos

				morphPacket.Type = "morph"
				morphPacket.Pos = *morphPos
				morphPacket.Morph = morph
				morphPacket.UUID = morph.UUID
				morphBuffer = append(morphBuffer, &morph)

				spawnMutex.Lock()
				spawnChannel <- morphPacket
				spawnMutex.Unlock()
			}

		}
	}
}

func broadcast() {
	// broadcast random morphs every so often
	// once-per-day goroutines

	// leaderboard
	// []usernames
	// morphs:120 rares: 0

	for {
		//player := <-broadcastChannel // this only updates the clients as the player moves

		select {
		case morphPacket := <-morphChannel:
			//fmt.Println("MORPHPOS: ", morphPacket.Pos)

			clientsMutex.Lock()

			for client := range clients {
				morpherror := client.WriteJSON(morphPacket)
				if morpherror != nil {
					fmt.Println(morpherror)
				}
			}

			clientsMutex.Unlock()

		case morphPacket := <-spawnChannel:

			clientsMutex.Lock()

			if len(morphBuffer) < morphMax {
				for client := range clients {
					morpherror := client.WriteJSON(morphPacket)
					if morpherror != nil {
						fmt.Println(morpherror)
					}
				}
			}

			clientsMutex.Unlock()

		case player := <-broadcastChannel:

			//msg, err := msgpack.Marshal(player)
			//if err != nil {
			//	panic(err)
			//

			clientsMutex.Lock()

			for client := range clients {
				err := client.WriteJSON(player)
				if err != nil {
					fmt.Println(err)
					client.Close()
					delete(clients, client)
				}

			}
			clientsMutex.Unlock()

		}
	}
}

func main() {
	//gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.Static("/client", "./client")

	router.LoadHTMLFiles("client/game.html")

	//config := cors.DefaultConfig()
	//config.AllowOrigins = []string{"http://127.0.0.1", "http://127.0.0.1:8081/", "http://10.0.0.4:8081"}
	//router.Use(cors.New(config))

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "game.html", gin.H{"message": "hi"})
	})

	// options:
	// - handle username with tokens and sessions GOOD

	//router.POST("/ws", func(c *gin.Context) {
	//	username := c.PostForm("username")
	//	wserve(c.Writer, c.Request, username)
	//})

	router.POST("/login", func(c *gin.Context) {
		username := c.PostForm("data")
		c.Set("userName", username)

		fmt.Print("User logged in: ", username)
		c.Redirect(http.StatusFound, "/ws")

	})

	router.GET("/ws", func(c *gin.Context) {

		username, _ := c.Get("userName")
		fmt.Print("USRNAME: ", username)

		// Set the Content-Disposition header to force download
		//c.Header("Content-Disposition", "attachment; filename=items.json")

		// Set the Content-Type header to application/json
		//c.Header("Content-Type", "application/json")

		wserve(c.Writer, c.Request, "")
	})

	go broadcast()
	go spawner()
	//go simulate()

	router.Use(gin.Recovery())

	fmt.Println("WebSocket server running")
	router.Run(":8080")
}

func simulate() {
	//delta := 1.0 / 60.0
	//accel := 24.0
	//decel := 0.95
	//player.RotY = 0

	//vel := NewVector3(0, 0, 0)
	///Pos := NewVector3(0, 0, 0)

	for {

		select {
		case player := <-broadcastChannel:
			// Check if player is colliding
			if player.Colliding {
				// Handle colliding state
				// ...
				select {
				case msg := <-keyChannel:
					if msg.Type == "beam" {
						fmt.Print("BEAMED")
						// lerp object up to y level 64 as long as this is true, else lerp down to y level 20
						// use object properties to determine lerp time
						// lerp(morphBuffer[collision.ObjectID].beamTime)

						// send id, morphpos
					} else {

					}

				}
			}
		}
	}
}

// fun little database idea: golang map + fsync
