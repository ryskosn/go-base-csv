package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "./statik"
	"github.com/gocarina/gocsv"
	"github.com/rakyll/statik/fs"
	"golang.org/x/oauth2"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

const (
	credentialsFile = "/creds.json"
	tokenFile       = "data/token.json"
	callbackURL     = "http://127.0.0.1:5903/callback"
	defaultBaseURL  = "https://api.thebase.in/"
)

var creds struct {
	ClientID     string `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
}

var statikFS http.FileSystem

func init() {
	sFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}
	statikFS = sFS

	credsf, err := statikFS.Open(credentialsFile)
	if err != nil {
		log.Fatal(err)
	}

	bytes, err := ioutil.ReadAll(credsf)
	if err != nil {
		log.Fatalf("Failed to read credentials file: %v", err)
	}
	json.Unmarshal(bytes, &creds)
}

// Unescape json string of error message sent from BASE server
// use like this:
// if err != nil {
// 	fmt.Println(UnescapeResponse(err.Error()))
// 	log.Fatalf("Token exchange failed: %v", err)
// }
func unescapeResponse(s string) string {
	type errorMsg struct {
		Error       string `json:"error"`
		Description string `json:"error_description"`
	}

	lines := strings.Split(s, "\n")
	resp := strings.Split(lines[1], " ")
	u, _ := strconv.Unquote("`" + resp[1] + "`")
	jsonStr := strings.TrimRight(u, `"`)

	var msg errorMsg
	if err := json.Unmarshal([]byte(jsonStr), &msg); err != nil {
		log.Fatalf("json.Unmarshal error: %v", err)
	}
	return fmt.Sprintf("%#v", msg)
}

// writeToken saves token to local file.
func writeToken(token *oauth2.Token) error {
	bytes, err := json.MarshalIndent(token, "", "  ")
	err = ioutil.WriteFile(tokenFile, bytes, 0777)
	fmt.Printf("token is saved to %v\n", tokenFile)
	return err
}

// // readToken reads token from local file.
// func readToken(filename string) *oauth2.Token {
// 	f, err := os.Open(filename)
// 	if err != nil {
// 		log.Fatalf("open token file: %v", err)
// 	}
// 	defer f.Close()
// 	var bytes []byte
// 	bytes, err = ioutil.ReadAll(f)
// 	if err != nil {
// 		log.Fatalf("read token file: %v", err)
// 	}
// 	var token *oauth2.Token
// 	if err = json.Unmarshal(bytes, &token); err != nil {
// 		log.Fatal(err)
// 	}
// 	return token
// }

// // readTokenBinData reads from binary data.
// // using go-bindata.
// func readTokenBinData() *oauth2.Token {
// 	bytes, err := Asset("data/token.json")
// 	if err != nil {
// 		log.Fatalf("open token from bindata failed: %v", err)
// 	}
// 	var token *oauth2.Token
// 	if err = json.Unmarshal(bytes, &token); err != nil {
// 		log.Fatal(err)
// 	}
// 	return token
// }

// exists checks file.
func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func decodeBody(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	return dec.Decode(out)
}

// Order から OrderToCsv へフィールドを代入
func convertOrderDetail(o Order) OrderToCsv {
	oc := new(OrderToCsv)
	oc.UniqueKey = o.Order.UniqueKey

	// 日時
	const layout = "2006/01/02 15:04:05"
	t := time.Unix(int64(o.Order.Ordered), 0)
	oc.Ordered = t.Format(layout)

	// 氏名
	s := []string{o.Order.LastName, o.Order.FirstName}
	oc.FullName = strings.Join(s, "")
	oc.LastName = o.Order.LastName
	oc.FirstName = o.Order.FirstName
	oc.Total = o.Order.Total
	oc.MailAddress = o.Order.MailAddress
	oc.Tel = o.Order.Tel

	// 備考
	oc.Remark = o.Order.Remark

	// 割引きクーポン
	oc.Discount = o.Order.OrderDiscount.Discount
	oc.Note = o.Order.OrderDiscount.Note
	re := regexp.MustCompile(`クーポン「.+?」\((.+?)\)`)
	if re.MatchString(oc.Note) {
		result := re.FindAllStringSubmatch(oc.Note, -1)
		oc.Coupon = result[0][1]
	}

	// 商品
	oc.ItemID = o.Order.OrderItems[0].ItemID
	oc.Title = o.Order.OrderItems[0].Title
	oc.Price = o.Order.OrderItems[0].Price

	// 住所
	oc.Country = o.Order.Country
	oc.ZipCode = o.Order.ZipCode
	oc.Prefecture = o.Order.Prefecture

	// oc.Address = o.Order.Address
	// oc.Address2 = o.Order.Address2
	return *oc
}

// NewConfig construct Config
func NewConfig() *oauth2.Config {
	if len(creds.ClientID) == 0 {
		log.Println("missing clientID")
		return nil
	}
	if len(creds.ClientSecret) == 0 {
		log.Println("missing clientSecret")
		return nil
	}
	if len(callbackURL) == 0 {
		log.Println("missing callbackURL")
		return nil
	}
	parsedURL, err := url.ParseRequestURI(callbackURL)
	if err != nil {
		fmt.Printf("Failed to parse url: %s", callbackURL)
		return nil
	}
	return &oauth2.Config{
		ClientID:     creds.ClientID,
		ClientSecret: creds.ClientSecret,
		RedirectURL:  parsedURL.String(),
		Scopes:       []string{"read_users_mail", "read_items", "read_orders"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://api.thebase.in/1/oauth/authorize",
			TokenURL: "https://api.thebase.in/1/oauth/token",
		},
	}
}

// Client is my http client.
type Client struct {
	*http.Client
}

// Orders struct
// start_ordered: yyyy-mm-dd
// end_ordered  : yyyy-mm-dd
// limit        : default 20, MAX 100
// offset       : default 0
func (c *Client) getOrders(limit int, offset int) Orders {
	values := url.Values{}
	values.Add("limit", fmt.Sprint(limit))
	values.Add("offset", fmt.Sprint(offset))
	req, err := http.NewRequest("GET", "https://api.thebase.in/1/orders", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.URL.RawQuery = values.Encode()
	resp, err := c.Client.Do(req)
	if err != nil {
		log.Fatalf("Failed to get orders: %v\n", err)
	}
	orders := new(Orders)
	if err := decodeBody(resp, orders); err != nil {
		log.Fatal(err)
	}
	return *orders
}

func (c *Client) getOrderDetail(key string) Order {
	url := fmt.Sprintf("https://api.thebase.in/1/orders/detail/%v", key)
	resp, err := c.Client.Get(url)
	if err != nil {
		log.Fatalf("Failed to get order detail:%v", err)
	}
	o := new(Order)
	if err := decodeBody(resp, o); err != nil {
		log.Fatalf("Order decode failed:%v", err)
	}
	return *o
}

// NewClient is constructor of Client.
func NewClient(config *oauth2.Config, token *oauth2.Token) *Client {
	var c Client
	c.Client = config.Client(oauth2.NoContext, token)
	return &c
}

// Token がない場合、認証コードの取得および、Token との Exchange を行う
func initToken() {
	config := NewConfig()

	// 認証コードを引数で受け取る
	// flag.Arg() は該当するものがない場合、空文字列を返す
	code := flag.Arg(0)

	// 引数なしの場合、ブラウザで認証する
	if code == "" {
		url := config.AuthCodeURL("state", oauth2.AccessTypeOffline)
		fmt.Println("ブラウザで以下の URL にアクセスし、認証してください。")
		fmt.Println(url)
		// open は Mac 依存かもしれない
		if err := exec.Command("open", url).Run(); err != nil {
			log.Fatal(err)
		}
		fmt.Println("$ ./go-base-csv --init=true <auth code>")
		return
	}
	// exchange code for token
	token, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Println(unescapeResponse(err.Error()))
		log.Fatalf("Failed to exchange token: %v\n", err)
	}
	if err := writeToken(token); err != nil {
		log.Fatalf("writeToken: %v\n", err)
	}
}

func main() {
	var (
		init   = flag.Bool("init", false, "initialize token, default false.")
		limit  = flag.Int("limit", 10, "how many orders from newest, default 10, max 100.")
		offset = flag.Int("offset", 0, "skip from newest, default 0.")
	)
	flag.Parse()

	if *init == true {
		initToken()
		return
	}

	f, err := statikFS.Open("/token.json")
	if err != nil {
		log.Fatalf("Failed to open token.json: %v\n", err)
	}

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Token seems broken, initialize it.")
		initToken()
	}

	var token *oauth2.Token
	if err = json.Unmarshal(bytes, &token); err != nil {
		log.Fatalf("Failed to Unmarshal token: %v\n", err)
	}

	config := NewConfig()
	client := NewClient(config, token)
	orders := client.getOrders(*limit, *offset)

	// Orders から n 件取り出す
	var n = *limit
	rows := make([]OrderToCsv, 0, n)
	for i := 0; i < n; i++ {
		key := orders.Orders[i].UniqueKey
		fmt.Printf("%v\n", key)
		o := client.getOrderDetail(key)
		oc := convertOrderDetail(o)
		rows = append(rows, oc)
	}

	// csv
	// https://qiita.com/nkumag/items/ef372ea35dcfbfa19310
	gocsv.SetCSVWriter(func(out io.Writer) *gocsv.SafeCSVWriter {
		writer := csv.NewWriter(transform.NewWriter(out, japanese.ShiftJIS.NewEncoder()))
		return gocsv.NewSafeCSVWriter(writer)
	})

	// 現在日時をファイル名にする
	now := time.Now()
	const layout = "20060102_150405"
	filename := now.Format(layout) + "_BASE.csv"
	filename = strings.TrimLeft(filename, "20")

	file, _ := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
	defer file.Close()
	if err := gocsv.MarshalFile(rows, file); err != nil {
		log.Fatalf("csv marshal: %v", err)
	}
}
