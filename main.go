package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

type Joke struct {
	ID    int    `json:"id" binding:"required"`
	Likes int    `json:"likes"`
	Joke  string `json:"joke" binding:"required"`
}

type LOGIN struct {
	ID       string `json:"id" binding:"required"`
	PASSWORD string `json:"password" binding:"required"`
}

type JWTtoken struct {
	TOKEN string `json:"token" binding:"required"`
}

/** we'll create a list of jokes */
var jokes = []Joke{
	Joke{1, 0, "Did you hear about the restaurant on the moon? Great food, no atmosphere."},
	Joke{2, 0, "What do you call a fake noodle? An Impasta."},
	Joke{3, 0, "How many apples grow on a tree? All of them."},
	Joke{4, 0, "Want to hear a joke about paper? Nevermind it's tearable."},
	Joke{5, 0, "I just watched a program about beavers. It was the best dam program I've ever seen."},
	Joke{6, 0, "Why did the coffee file a police report? It got mugged."},
	Joke{7, 0, "How does a penguin build it's house? Igloos it together."},
	Joke{8, 0, "Dad, did you get a haircut? No I got them all cut."},
	Joke{9, 0, "What do you call a Mexican who has lost his car? Carlos."},
	Joke{10, 0, "Dad, can you put my shoes on? No, I don't think they'll fit me."},
	Joke{11, 0, "Why did the scarecrow win an award? Because he was outstanding in his field."},
	Joke{12, 0, "Why don't skeletons ever go trick or treating? Because they have no body to go with."},
}

// Predefined ID and password. (THIS IS THE PROTOTYPE.)
var authID = os.Getenv("USER_ID")
var authPassword = os.Getenv("USER_PASSWORD")

func main() {

	// Set the router as the default one shipped with Gin
	router := gin.Default()

	// Serve the frontend
	router.Use(static.Serve("/", static.LocalFile("./views", true)))

	api := router.Group("/api")
	{
		api.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
		api.GET("/jokes", authMiddleware(), JokeHandler)
		api.POST("/jokes/like/:jokeID", authMiddleware(), LikeJoke)
		api.POST("/auth", auth)
	}
	// Start the app
	router.Run(":3000")
}

func auth(c *gin.Context) {

	var login LOGIN
	c.BindJSON(&login)

	// To-Do: Check ID and password is right or not.
	if authID != login.ID {
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusBadRequest, "")
		return
	}
	if authPassword != login.PASSWORD {
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusBadRequest, "")
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":      login.ID,
		"timestamp": int32(time.Now().Unix()),
	})

	// Sign and get the complete encoded token as a string using the secret.
	tokenString, err := token.SignedString([]byte(login.PASSWORD))

	if err != nil {
		log.Fatal("Faile to generate signed string.")
	}

	fmt.Println("user : " + login.ID + " / " + "pw : " + login.PASSWORD + "/" + "token string:" + tokenString)

	var jwtToken JWTtoken
	jwtToken = JWTtoken{tokenString}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, jwtToken)
}

// authMiddleware intercepts the requests, and check for a valid jwt token
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get the JWT token sting in Authorization header.
		var header = c.Request.Header.Get("Authorization")
		fmt.Println("auth : " + header)
		var tokenString = header[7:]
		
		// Parse and validate JWT token.
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte(authPassword), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			fmt.Println(claims["user"], claims["timestamp"])
		} else {
			fmt.Println(err)
			c.Abort()
			c.Writer.WriteHeader(http.StatusUnauthorized)
			c.Writer.Write([]byte("Unauthorized"))
			return
		}
	}
}

// JokeHandler returns a list of jokes available (in memory)
func JokeHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, jokes)
}

// LikeJoke ...
func LikeJoke(c *gin.Context) {
	// Check joke ID is valid
	if jokeid, err := strconv.Atoi(c.Param("jokeID")); err == nil {
		// find joke and increment likes
		for i := 0; i < len(jokes); i++ {
			if jokes[i].ID == jokeid {
				jokes[i].Likes = jokes[i].Likes + 1
			}
		}
		c.JSON(http.StatusOK, &jokes)
	} else {
		// the jokes ID is invalid
		c.AbortWithStatus(http.StatusNotFound)
	}
}

// getJokesByID returns a single joke
func getJokesByID(id int) (*Joke, error) {
	for _, joke := range jokes {
		if joke.ID == id {
			return &joke, nil
		}
	}
	return nil, errors.New("Joke not found")
}
