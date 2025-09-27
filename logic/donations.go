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
    Name        string  `json:"name"`
    Company     string  `json:"company"`
    Country     string  `json:"country"`
    State       string  `json:"state"`
    City        string  `json:"city"`
    Zip         string  `json:"zip"`
    Address     string  `json:"address"`
    Address2    string  `json:"address2"`
    Phone       string  `json:"phone"`
}

type OtpApiStartMessageRecurring struct {
    Times       int     `json:"times"`
    Until       string  `json:"until"`
    MaxAmount   int     `json:"maxAmount"`
}

type OtpApiStartMessage struct {
    Salt            string                          `json:"salt"`   // 32 char
    Merchant        string                          `json:"merchant"`
    OrderRef        string                          `json:"orderRef"`
    Currency        string                          `json:"currency"`
    CustomerEmail   string                          `json:"customerEmail"`
    Language        string                          `json:"language"`
    SdkVersion      string                          `json:"sdkVersion"`
    Methods         []string                        `json:"methods"`
    Total           string                          `json:"total"`
    Timeout         string                          `json:"timeout"`
    Url             string                          `json:"url"`    // Redirect back URL
    Invoice         *OtpApiStartMessageInvoice      `json:"invoice,omitempty"`
    Recurring       *OtpApiStartMessageRecurring    `json:"recurring,omitempty"`
}

type SimpleResponse struct {
    ResponseCode    int     `json:"r"`
    TransactionNum  int     `json:"t"`
    Event           string  `json:"e"`
    Merchant        string  `json:"m"`
    OrderId         string  `json:"o"`
}

type OtpJsonResponse struct {
    ErrorCodes      []int
    Merchant        string
    Salt            string
    OrderRef        string
    TransactionId   int
    Currency        string
    Timeout         string
    Total           float64
    PaymentUrl      string
    Tokens          []string
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

func otpGetTimeFormat(t time.Time) string {
    return t.Local().Format("2006-01-02T15:04:05-07:00")
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

func RedirectToOtpApi(dict dictionary.Dictionary, donation Donation) (OtpReturnPublic, error) {
    url := config.Config.Donation.SimplePayURL + "/payment/v2/start"
    log.Printf("URL> %s\n", url)
    log.Printf("mer> %s\n", config.Config.Donation.Merchant)

    donation.Status = "Initiated" // TODO: Create a struct for these..
    err := donation.Add()
    if nil != err {
        return OtpReturnPublic{}, err
    }

    simple_start := OtpApiStartMessage{
        Salt:           generateSalt(32),
        Merchant:       config.Config.Donation.Merchant,
        // FIXME: ::
        OrderRef:       donation.Id,
        Currency:       "HUF",
        CustomerEmail:  donation.Email,
        Language:       strings.ToUpper(dict.Meta.CountryCode),
        SdkVersion:     "SimplePayV2.1_Payment_PHP_SDK_2.0.7_190701:dd236896400d7463677a82a47f53e36e",
        Methods:        []string{"CARD"},
        Total:          strconv.Itoa(int(donation.Amount)),
        Timeout:        otpGetTimeFormat(time.Now().Add(5 * time.Minute)),
        Url:            config.Config.Donation.SimplePayReturnURL,
    }

    if donation.Recurring {
        simple_start.Recurring = &OtpApiStartMessageRecurring{
            Times: 12,
            Until: otpGetTimeFormat(time.Now().Add(8760 * time.Hour)), // multiple up to 1 year
            MaxAmount: int(donation.Amount),
        }
        log.Println("recurring...")
    }

    if donation.InvoiceNeeded {
        iv := donation.Invoice
        simple_start.Invoice = &OtpApiStartMessageInvoice{
            Name:       iv.Name,
            Company:    iv.Company,
            Country:    iv.Country,
            State:      iv.State,
            City:       iv.City,
            Zip:        iv.Zip,
            Address:    iv.Address,
            Address2:   iv.Address2,
            Phone:      iv.Phone,
        }
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
        return OtpReturnPublic{PaymentUrl: "/donate/" + donation.Id}, errors.New("Signature mismatch!")
    }

    retStuff := OtpJsonResponse{}
    err = json.Unmarshal(respBody, &retStuff)
    if nil != err || 0 != len(retStuff.ErrorCodes) {
        retStuff.PaymentUrl = "/donate/" + donation.Id
    }
    log.Println(retStuff)

    donation.Tokens = retStuff.Tokens
    donation.TransactionId = retStuff.TransactionId
    donation.Update()

    return OtpReturnPublic{PaymentUrl: retStuff.PaymentUrl}, err
}

func ProgressOtpReply(r string, s string) (string, bool, error) {
    payload, err := base64.StdEncoding.DecodeString(r)

    if nil != err {
        log.Println("base64 decode error.")
        log.Println(err)
        return "", false, err
    }

    if !signatureMatch(payload, s) {
        log.Println("Signature mismatch..")
        return "", false, errors.New("Payload ignature mismatch..")
    }

    simple_resp := SimpleResponse{}
    json.Unmarshal(payload, &simple_resp)
    log.Println(simple_resp)

    donation := Donation{Id: simple_resp.OrderId}
    donation.Select()

    donation.Status = simple_resp.Event
    donation.Successful = "SUCCESS" == simple_resp.Event
    log.Println("successful: ", donation.Successful)
    if donation.Successful {
        log.Println(donation.Occurences)
        donation.Occurences = []time.Time{time.Now()}
        log.Println("LastOccurance: ", donation.Occurences)
        donation.RecurringActive = donation.Recurring
    }
    donation.Update()

    // TODO: Send Email

    return donation.Id, donation.Successful, nil
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
    donation.Id = ddon.Id.Hex()

    if donation.InvoiceNeeded {
        iv := donation.Invoice.UnMap()
        iv.Id = primitive.NewObjectID()
        iv.Add()
        ddon.Invoice = iv.Id
        donation.Invoice.Id = iv.Id.Hex()
    }

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
