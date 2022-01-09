package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4"
	"golang.org/x/net/publicsuffix"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"sync"
	"time"
)

type MyJsonName struct {
	Info struct {
		Page    int64  `json:"page"`
		Results int64  `json:"results"`
		Seed    string `json:"seed"`
		Version string `json:"version"`
	} `json:"info"`

	Results []struct {
		Cell string `json:"cell"`
		Dob  struct {
			Age  int64  `json:"age"`
			Date string `json:"date"` //birthday
		} `json:"dob"`
		Email  string `json:"email"`  //email
		Gender string `json:"gender"` //gender
		ID     struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"id"`
		Location struct {
			City        string `json:"city"` //location
			Coordinates struct {
				Latitude  string `json:"latitude"`
				Longitude string `json:"longitude"`
			} `json:"coordinates"`
			Postcode string `json:"postcode"`
			State    string `json:"state"`
			Street   struct {
				Name   string      `json:"name"`
				Number interface{} `json:"number"` //for averting type mismatch err
			} `json:"street"`
			Country  string `json:"country"`
			Timezone struct {
				Description string `json:"description"`
				Offset      string `json:"offset"`
			} `json:"timezone"`
		} `json:"location"`
		Login struct {
			Md5      string `json:"md5"`
			Password string `json:"password"`
			Salt     string `json:"salt"`
			Sha1     string `json:"sha1"`
			Sha256   string `json:"sha256"`
			Username string `json:"username"` //username password email
			UUID     string `json:"uuid"`
		} `json:"login"`
		Name struct {
			First string `json:"first"` //name
			Last  string `json:"last"`  //lastname
			Title string `json:"title"`
		} `json:"name"`
		Nat     string `json:"nat"`
		Phone   string `json:"phone"` //mobile
		Picture struct {
			Large     string `json:"large"`
			Medium    string `json:"medium"` //avatar
			Thumbnail string `json:"thumbnail"`
		} `json:"picture"`
		Registered struct {
			Age  int64  `json:"age"`
			Date string `json:"date"`
		} `json:"registered"`
	} `json:"results"`
}

//YtEmbedJson is used for Youtube embed link requests from oembed api
type YtEmbedJson struct {
	Html string                 `json:"html"`
	X    map[string]interface{} `json:"-"`
}

var db_Host = "127.0.0.1"
var db_Port = "5432"
var db_Name = "dbForFaceClone"
var db_User = "postgres"
var db_Password = "123456"
var Ctx = context.Background()
var wg sync.WaitGroup

func main() {

	/*link := "https://www.youtube.com/watch?ssv=O3VPs9b_HZE"
	fmt.Println(GetYtEmbed(link))*/

	/*	wg.Add(1)
		userList := connectDBAndFetch()
		wg.Wait()
		for i := 0; i < 10; i++ {
			wg.Add(len(userList))
			for _, username := range userList {
				fmt.Println(username)
				if populatePost(username) {
					wg.Done()
				} else {
					log.Println("problem with populate post:", username)
					wg.Done()
				}
			}
			wg.Wait()
		}*/

}

//random user
//create account
//login
//edit profile
//update avatar
//logout

//populatePost function does populate Post with lorem ipsum generator
func populatePost(username string) bool {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Println("cookiejar creation error:", err)
		return false
	}

	c := http.Client{ //no need to redirect
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil
		},
		Jar: jar, //cookies will be stored here
	}
	if loginTo(c, username) {
		if makePost(c, username) {
			if logOut(c) {
				fmt.Println("successfully posted:", username)
				c.CloseIdleConnections()
				return true
			}
		}
	}
	return false
}

//populateUser func populates user
func populateUser() bool {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Println("cookiejar creation error:", err)
		return false
	}
	c := http.Client{ //no need to redirect
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil
		},
		Jar: jar, //cookies will be stored here
	}
	profile := userGenerator()

	//time.Sleep(1000 * time.Millisecond)
	if createAccount(profile.Results[0].Login.Username) {
		//time.Sleep(2000 * time.Millisecond)
		if loginTo(c, profile.Results[0].Login.Username) {
			//	time.Sleep(1000 * time.Millisecond)
			if editProfile(c, profile) {
				//	time.Sleep(1000 * time.Millisecond)
				if updateAvatar(c, profile.Results[0].Picture.Large) {
					//		time.Sleep(1000 * time.Millisecond)
					if logOut(c) {
						fmt.Println("succesfully created:", profile.Results[0].Login.Username)
						return true
					}
				}
			}
		}
	}

	return false
}

//logOut logs out
func logOut(c http.Client) bool {
	resp, err := c.Get("http://localhost:8080/logout")
	defer resp.Body.Close()
	if err != nil {
		log.Println("logout failed:", err)
		return false
	}
	return true
}

func makePost(c http.Client, username string) bool {
	resp, err := c.PostForm("http://localhost:8080/postIt", url.Values{
		"postmessage": {loremipsumGenerator()},
	})
	defer resp.Body.Close()
	if err != nil {
		log.Println("post make request failed:", err)
		return false
	}
	return true
}

//editProfile sends post request for edit profile div
func editProfile(c http.Client, profile MyJsonName) bool {
	mobile := profile.Results[0].Cell
	modile := []byte{}
	for _, v := range mobile {
		if v >= 48 && v <= 57 {
			modile = append(modile, byte(v))
		}
	}
	resp, err := c.PostForm("http://localhost:8080/updateprofile", url.Values{
		"firstname":    {profile.Results[0].Name.First},
		"lastname":     {profile.Results[0].Name.Last},
		"gender":       {profile.Results[0].Gender},
		"birthday":     {profile.Results[0].Dob.Date},
		"mobilenumber": {string(modile)},
		"country":      {profile.Results[0].Location.Country},
	})
	defer resp.Body.Close()
	if err != nil {
		log.Println("post updateprofile failed:", err)
		return false
	}
	return true
}

//createAccount uses seed as username, password and also email creatively
func createAccount(seed string) bool {
	username := seed
	password := seed
	email := fmt.Sprintf("%s@example.com", seed)
	resp, err := http.PostForm("http://localhost:8080/checkReg", url.Values{
		"usernameReg": {username},
		"passwordReg": {password},
		"emailReg":    {email},
	})
	defer resp.Body.Close()
	if err != nil {
		log.Println("create account failed while post:", err)
		return false
	}

	return true

}

//loginTo func simply logs in
func loginTo(c http.Client, seed string) bool {

	formData := url.Values{
		"usernameL": {seed},
		"passwordL": {seed},
	}
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/checkAuthLog", strings.NewReader(formData.Encode()))
	if err != nil {
		log.Println("new req error", err)
		return false
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") //header must be set in order to send form values in request body
	resp, err := c.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Println("client do error:", err)
		return false
	}
	return true

}

//loremipsumGenerator populates message bodies
func loremipsumGenerator() string {
	resp, err := http.Get("https://baconipsum.com/api/?type=all-meat&paras=2&start-with-lorem=1")
	if err != nil {
		println(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	return string(body)[2 : len(body)-2]
}

//updateAvatar func update profile picture in link url
func updateAvatar(c http.Client, link string) bool {
	resp, err := http.Get(link)
	if err != nil {
		log.Println("getlink error upavatar:", err)
		return false
	}

	body := &bytes.Buffer{}

	writer := multipart.NewWriter(body)
	fw, err := writer.CreateFormFile("change_pp", "avatar.jpg")
	if err != nil {
		log.Println("readform error upavatar:", err)
		return false
	}
	_, err = io.Copy(fw, resp.Body)
	if err != nil {
		log.Println("copy error upavatar:", err)
		return false
	}

	writer.Close() //kapatilmazsa unexpected eof hatasi veriyor
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/updatepp", body)
	if err != nil {
		log.Println("avatar update failed", err)
		return false
	}
	req.Header.Set("Content-Type", writer.FormDataContentType()) //header must be set in order to send form values in request body
	req.Header.Set("Boundary", writer.Boundary())

	resp, err = c.Do(req)
	defer resp.Body.Close()
	defer req.Body.Close()

	if err != nil {
		log.Println("avatar update failed", err)
		return false
	}

	return true
}

//userGenerator func returns a random user as MyJsonName type
func userGenerator() MyJsonName {
	resp, err := http.Get("https://randomuser.me/api/1.3/")
	if err != nil {
		println(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	randomy := MyJsonName{}
	err = json.Unmarshal(body, &randomy)
	if err != nil {
		log.Println("unmarshal error:", err)
	}
	return randomy
}

func GetYtEmbed(shortlink string) string {

	reqLink := fmt.Sprintf("https://www.youtube.com/oembed?url=%s&format=json", shortlink)
	resp, err := http.Get(reqLink)
	if err != nil {
		log.Println("Video link json couldn't be get", err)
		return "Video can't be loaded!"
	}
	defer resp.Body.Close()
	jsonVideo := YtEmbedJson{}
	out, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Println("Video link json couldn't be read", err)
		return "Video can't be loaded"
	}
	if string(out) == "Not Found" {
		log.Println("invalid youtube url")
		return "false"
	}

	err = json.Unmarshal(out, &jsonVideo)
	if err != nil {
		log.Println("Video link json unmarshal failed", err)
		return "Video can't be loaded"
	}
	return jsonVideo.Html
}

func connectDBAndFetch() []string {
	var databaseURL = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", db_Host, db_Port, db_User, db_Password, db_Name)
	ctxT, cancel := context.WithTimeout(Ctx, 5*time.Second)
	defer cancel()
	conn, err := pgx.Connect(ctxT, databaseURL)
	defer conn.Close(ctxT)
	if err != nil {
		log.Println("DB connection error:", err)
	}
	//check whether connection is ok or not
	err = conn.Ping(ctxT)
	if err != nil {
		log.Println("Ping to DB error:", err)
	}

	rows, err := conn.Query(ctxT, "SELECT username FROM user_creds")
	defer rows.Close()
	if err != nil {
		log.Println("query username list failed:", err)
		return nil
	}
	userList := []string{}
	for rows.Next() { //while there is a next
		tempUser := ""
		err := rows.Scan(&tempUser)
		if err != nil {
			log.Println("row scan failed:", err)
		}
		userList = append(userList, tempUser)
	}
	wg.Done()
	return userList
}
