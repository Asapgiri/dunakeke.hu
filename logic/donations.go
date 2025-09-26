package logic

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"dunakeke/config"
	"dunakeke/dbase"
	"dunakeke/dictionary"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// =====================================================================================================================
// private api logic

type OtpApiStartMessageInvoice struct {
    Name        string
    Company     string
    Country     string
    State       string
    City        string
    Zip         string
    Address     string
    Address2    string
    Phone       string
}

type OtpApiStartMessage struct {
    Salt            string                      `json:"salt"`   // 32 char
    Merchant        string                      `json:"merchant"`
    OrderRef        string                      `json:"orderRef"`
    Currency        string                      `json:"currency"`
    CustomerEmail   string                      `json:"customerEmail"`
    Language        string                      `json:"language"`
    SdkVersion      string                      `json:"sdkVersion"`
    Methods         []string                    `json:"methods"`
    Total           string                      `json:"total"`
    Timeout         string                      `json:"timeout"`
    Url             string                      `json:"url"`    // Redirect back URL
    // Invoice         OtpApiStartMessageInvoice   `json:"invoice"`
}

type SimpleResponse struct {
    ResponseCode    int     `json:"r"`
    TransactionNum  int     `json:"t"`
    Event           string  `json:"e"`
    Merchant        string  `json:"m"`
    OrderId         string  `json:"o"`
}

type OtpJsonResponse struct {
    Merchant        string
    Salt            string
    OrderRef        string
    TransactionId   int
    Currency        string
    Timeout         string
    Total           float64
    PaymentUrl      string
}

type OtpReturnPublic struct {
    PaymentUrl      string
}

type MerchantHasher struct {
    Body        string
    Merchant    string
    Hash        string
}

func otpGenerateSignature(body []byte) string {
    mac := hmac.New(sha512.New384, []byte(config.Config.Donation.SecretKey))
    mac.Write(body)
    hmacSum := mac.Sum(nil)
    return base64.StdEncoding.EncodeToString(hmacSum)
}

func otpGetTimeout() string {
    return time.Now().Add(5 * time.Minute).Local().Format("2006-01-02T15:04:05-07:00")
}

func generateSalt(length int) string {
    bytes := make([]byte, length)
    rand.Read(bytes)
    return base64.StdEncoding.EncodeToString(bytes)[:length]
}

func signatureMatch(payload []byte, signature string) bool {
    mac := hmac.New(sha512.New384, []byte(config.Config.Donation.SecretKey))
    mac.Write(payload)
    expectedMac := mac.Sum(nil)

    decodedSig, err := base64.StdEncoding.DecodeString(signature)
    if nil != err {
        return false
    }

    return hmac.Equal(decodedSig, expectedMac)
}

func RedirectToOtpApi(dict dictionary.Dictionary, amount float64, email string) (OtpReturnPublic, error) {
    url := config.Config.Donation.SimplePayURL + "/payment/v2/start"
    log.Printf("URL> %s\n", url)
    log.Printf("mer> %s\n", config.Config.Donation.Merchant)

    timeout := otpGetTimeout()

    simple_start := OtpApiStartMessage{
        Salt:           generateSalt(32),
        Merchant:       config.Config.Donation.Merchant,
        // FIXME: ::
        OrderRef:       timeout,
        Currency:       "HUF",
        CustomerEmail:  email,
        Language:       strings.ToUpper(dict.Meta.CountryCode),
        SdkVersion:     "SimplePayV2.1_Payment_PHP_SDK_2.0.7_190701:dd236896400d7463677a82a47f53e36e",
        Methods:        []string{"CARD"},
        Total:          strconv.Itoa(int(amount)),
        Timeout:        timeout,
        Url:            config.Config.Donation.SimplePayReturnURL,
        // TODO: Add invoice ...
    }

    body, _ := json.Marshal(simple_start)
    log.Println(string(body))

    req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
    req.Header.Set("merchantKey", config.Config.Donation.SecretKey)
    req.Header.Set("Signature", otpGenerateSignature(body))
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if nil != err {
        log.Println("SimplePay ERROR:")
        log.Println(err)
        return OtpReturnPublic{}, err
    }
    defer resp.Body.Close()

    log.Println("response status: ", resp.Status)
    log.Println("response headers: ", resp.Header)

    respBody, err := io.ReadAll(resp.Body)
    if nil != err {
        return OtpReturnPublic{}, err
    }
    log.Println("response body: ", string(respBody))

    if !signatureMatch(respBody, resp.Header.Get("Signature")) {
        return OtpReturnPublic{}, errors.New("Signature mismatch!")
    }

    retStuff := OtpJsonResponse{}
    err = json.Unmarshal(respBody, &retStuff)

    return OtpReturnPublic{PaymentUrl: retStuff.PaymentUrl}, err
}

func ProgressOtpReply(r string, s string) bool {
    decoded, err := base64.StdEncoding.DecodeString(r)
    simple_resp := SimpleResponse{}

    if !signatureMatch(decoded, s) {
        log.Println("Signature mismatch..")
        return false
    }

    json.Unmarshal(decoded, &simple_resp)

    log.Println(simple_resp)
    log.Println(err)

    return "SUCCESS" == simple_resp.Event
}


// =====================================================================================================================
// "public" page logic

func (donation *Donation) List() []Donation {
    ddon := dbase.Donation{}
    ddons, _ := ddon.List()

    donations := make([]Donation, len(ddons))
    for i, d := range(ddons) {
        donations[i].Map(d)
    }

    return donations
}

func (donation *Donation) Add() error {
    ddon := donation.UnMap()
    ddon.Id = primitive.NewObjectID()

    return ddon.Add()
}

func (donation *Donation) Update() error {
    ddon := donation.UnMap()
    return ddon.Update()
}

func (donation *Donation) Select() error {
    ddon := dbase.Donation{}
    oid, _ := primitive.ObjectIDFromHex(donation.Id)
    err := ddon.Select(oid)
    if nil != err {
        return err
    }

    donation.Map(ddon)
    return nil
}


func (do *DonationOption) List() []DonationOption {
    ddon := dbase.DonationOption{}
    ddons, err := ddon.List()
    if nil != err {
        log.Println(err)
        return []DonationOption{}
    }

    donations := make([]DonationOption, len(ddons))
    for i, d := range(ddons) {
        donations[i].Map(d)
    }

    sort.Slice(donations, func(i, j int) bool { return donations[i].Amount < donations[j].Amount })

    return donations
}

func (do *DonationOption) Add() error {
    ddon := do.UnMap()
    ddon.Id = primitive.NewObjectID()
    return ddon.Add()
}

func (do *DonationOption) Update() error {
    ddon := do.UnMap()
    return ddon.Update()
}

func (do *DonationOption) Select() error {
    ddon := dbase.DonationOption{}
    oid, _ := primitive.ObjectIDFromHex(do.Id)
    err := ddon.Select(oid)
    if nil != err {
        return err
    }

    do.Map(ddon)
    return nil
}
