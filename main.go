package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
			Street   string `json:"street"`
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

}

func loginto() {

	_, err := http.PostForm("http://localhost:8080/logout", nil)
	time.Sleep(1 * time.Second)
	_, err = http.PostForm("http://localhost:8080/checkAuthLog", url.Values{
		"usernameL": {"fukaraca"},
		"passwordL": {"Password"},
	})
	fmt.Println(err)
}

func loremipsumgenerator() {
	resp, err := http.Get("https://baconipsum.com/api/?type=all-meat&paras=2&start-with-lorem=1")
	if err != nil {
		println(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Println(string(body)[2 : len(body)-2])
}

func usergenerator() {
	resp, err := http.Get("https://randomuser.me/api/1.3/")
	if err != nil {
		println(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	randomy := MyJsonName{}
	json.Unmarshal(body, &randomy)
	fmt.Println(string(body) /*[2 : len(body)-2]*/)
	fmt.Println(randomy.Results[0].Name.First)
}
