package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/net/publicsuffix"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
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
			State    string `json:"state"` //country
			Street   struct {
				Name   string `json:"name"`
				Number string `json:"number"`
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

func main() {
	/*_, err := http.Get("http://localhost:8080/logout")
	time.Sleep(200 * time.Millisecond)

	//create a client
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Fatal("cookiejar creation error:", err)
	}
	c := http.Client{ //no need to redirect
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil
		},
		Jar: jar, //cookies will be stored here
	}
	formData := url.Values{
		"usernameL": {"doodi44"},
		"passwordL": {"doodi44"},
	}
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/checkAuthLog", strings.NewReader(formData.Encode()))
	fmt.Println("newrequest error:", err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") //header must be set in order to send form values in request body
	resp, err := c.Do(req)
	fmt.Println("client do error:", err)
	fmt.Println("request jari:", c.Jar.Cookies(req.URL))
	for _, v := range c.Jar.Cookies(req.URL) {
		fmt.Println("cookieler:", v.Name, v.Value)
	}
	fmt.Println()
	defer resp.Body.Close()
	denemeedit(c, "deneme3")
	for _, v := range c.Jar.Cookies(req.URL) {
		fmt.Println(v.Name, v.Value)
	}*/
	resp, _ := http.Get("http://localhost:8080/logout")
	defer resp.Body.Close()

	fmt.Println(populateUser())

}

func denemeedit(c http.Client, str string) {

	formData1 := url.Values{
		"firstname":    {str},
		"lastname":     {str},
		"gender":       {"male"},
		"birthday":     {"2021-12-04 23:28:15.000000 +00:00"},
		"mobilenumber": {str},
		"country":      {str},
	}
	_, err := c.PostForm("http://localhost:8080/updateprofile", formData1)
	//req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/updateprofile", strings.NewReader(formData1.Encode()))
	fmt.Println(err)
	//resp, err := c.Do(req)

}

func logOut(c http.Client) bool {
	resp, err := c.Get("http://localhost:8080/logout")
	defer resp.Body.Close()
	if err != nil {
		log.Println("logout failed:", err)
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

//createAccount
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
func loremipsumGenerator() {
	resp, err := http.Get("https://baconipsum.com/api/?type=all-meat&paras=2&start-with-lorem=1")
	if err != nil {
		println(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Println(string(body)[2 : len(body)-2])
}

func updateAvatar(c http.Client, link string) bool {
	fmt.Println(link)
	resp, err := http.Get(link)
	if err != nil {
		log.Println("getlink error upavatar:", err)
		return false
	}
	/*formData := url.Values{
		"change_pp": {link},
	}*/
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
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/updatepp", body) //strings.NewReader(formData.Encode())
	if err != nil {
		log.Println("avatar update failed", err)
		return false
	}
	req.Header.Set("Content-Type", writer.FormDataContentType()) //header must be set in order to send form values in request body
	req.Header.Set("boundary", writer.Boundary())
	//err = writer.Close() //bu yapılmazsa unexpected error hatası veriyor
	//fmt.Println(err)
	resp, err = c.Do(req)
	defer resp.Body.Close()
	defer req.Body.Close()
	/*_, err := c.PostForm("http://localhost:8080/updatepp", url.Values{
		"change_pp": {link},
	})*/
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

//random user
//create account
//login
//edit profile
//update avatar
//logout

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

	time.Sleep(1000 * time.Millisecond)
	if createAccount(profile.Results[0].Login.Username) {
		time.Sleep(2000 * time.Millisecond)
		if loginTo(c, profile.Results[0].Login.Username) {
			time.Sleep(1000 * time.Millisecond)
			if editProfile(c, profile) {
				time.Sleep(1000 * time.Millisecond)
				if updateAvatar(c, profile.Results[0].Picture.Large) {
					time.Sleep(1000 * time.Millisecond)
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
