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
    Recurring       *OtpApiStartMessageRecurring    `json:"recurring,omitempty"`
}

type SimpleResponse struct {
    ResponseCode    int     `json:"r"`
    TransactionNum  int     `json:"t"`
    Event           string  `json:"e"`
    Merchant        string  `json:"m"`
    OrderId         string  `json:"o"`
}

type SimpleQuery struct {
    Merchant        string      `json:"merchant"`
    OrderRefs       []string    `json:"orderRefs,omitempty"`
    TransactionIds  []string    `json:"transactionIds,omitempty"`
    Salt            string      `json:"salt"`
    SdkVersion      string      `json:"sdkVersion"`
}

type SimpleQueryTranstactions struct {
     Salt           string      `json:"salt"`
     Merchant       string      `json:"merchant"`
     OrderRef       string      `json:"orderRef"`
     Total          int         `json:"total"`
     TransactionId  int         `json:"transactionId"`
     Status         string      `json:"status"`
     ResultCode     string      `json:"resultCode"`
     RemainingTotal int         `json:"remainingTotal"`
     PaymentDate    string      `json:"paymentDate"`
     FinishDate     string      `json:"finishDate"`
     Method         string      `json:"method"`
}

type SimpleQueryResponse struct {
    Salt            string                      `json:"salt"`
    Merchant        string                      `json:"merchant"`
    Transactions    []SimpleQueryTranstactions  `json:"transactions"`
    TotalCount      int                         `json:"totalCount"`
}

type SimpleIpn struct {
    Salt            string  `json:"salt"`
    OrderRef        string  `json:"orderRef"`
    Method          string  `json:"method,omitempty"`
    Merchant        string  `json:"merchant"`
    FinishDate      string  `json:"finishDate"`
    PaymentDate     string  `json:"paymentDate"`
    TransactionId   int     `json:"transactionId"`
    Status          string  `json:"status"`
    ReceiveDate     string  `json:"receiveDate,omitempty"`
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

const simpleSdkVersion = "SimplePayV2.1_Payment_PHP_SDK_2.0.7_190701:dd236896400d7463677a82a47f53e36e"

func translateStatus(donation *Donation, otpStatus string) {
    switch otpStatus {
    case "SUCCESS", "FINISHED":
        donation.Status = "SUCCESSFUL"
        donation.Successful = true
    case "FAIL", "NOTAUTHORIZED":
        donation.Status = "FAILURE"
        donation.Successful = false
    case "TIMEOUT":
        donation.Status = "TIMEOUT"
        donation.Successful = false
    case "CANCEL", "CANCELLED":
        donation.Status = "CANCELLED"
        donation.Successful = false
    case "INIT", "INPAYMENT", "INFRAUD", "AUTHORIZED", "REVERSED":
        donation.Status = "INPROGRESS"
        donation.Successful = false
    default:
        donation.Status = "UNKNOWN"
        donation.Successful = false
    }
}

func DonationGetPublicStatus(donation Donation, dict dictionary.Dictionary) string {
    switch donation.Status {
    case "SUCCESSFUL":
        return dict.Donate.StatusSuccess
    case "FAILURE":
        return dict.Donate.StatusFailure
    case "TIMEOUT":
        return dict.Donate.StatusTimeout
    case "CANCELLED":
        return dict.Donate.StatusCancelled
    case "INPROGRESS":
        return dict.Donate.StatusInProgress
    default:
        return "UNKNOWN"
    }
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

func otpHttpRequest(url string, content any) ([]byte, error) {
    body, _ := json.Marshal(content)

    req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
    req.Header.Set("merchantKey", config.Config.Donation.SecretKey)
    req.Header.Set("Signature", otpGenerateSignature(body))
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if nil != err {
        return []byte{}, err
    }
    defer resp.Body.Close()

    respBody, err := io.ReadAll(resp.Body)
    if nil != err {
        return []byte{}, err
    }

    if !signatureMatch(respBody, resp.Header.Get("Signature")) {
        return []byte{}, errors.New("Signature mismatch!")
    }

    return respBody, nil
}

func CheckTransactionProgress(donation Donation, callback func(Donation)) {
    url := config.Config.Donation.SimplePayURL + "/payment/v2/query"
    log.Printf("URL> %s\n", url)

    requestTimes := (config.Config.Donation.TimeoutMinutes * 60) / config.Config.Donation.PollSeconds
    for range(requestTimes) {
        time.Sleep(time.Duration(config.Config.Donation.PollSeconds) * time.Second)

        query := SimpleQuery{
            Merchant:   config.Config.Donation.Merchant,
            OrderRefs:  []string{donation.Id},
            Salt:       generateSalt(32),
            SdkVersion: simpleSdkVersion,
        }

        respBody, err := otpHttpRequest(url, query)
        if nil != err {
            log.Println("SimplePay Query ERROR: ", err)
            continue
        }

        ret := SimpleQueryResponse{}
        err = json.Unmarshal(respBody, &ret)

        if ret.TotalCount > 0 {
            // FIXME: Statuses should be mapped to result values...
            err := donation.SelectLock()
            if nil != err {
                log.Println("Failed to select donation...")
                continue
            }
            translateStatus(&donation, ret.Transactions[0].Status)
            donation.UpdateUnlock()

            if donation.Successful || "FAILURE" == donation.Status {
                callback(donation)
                return
            }
        }
    }
    log.Println("checking TIMED OUT")
}

func DonationIpn(w http.ResponseWriter, signature string, body []byte, callback func(Donation)) ([]byte, error) {
    if !signatureMatch(body, signature) {
        return []byte{}, errors.New("Signature mismatch..")
    }

    ipn := SimpleIpn{}
    err := json.Unmarshal(body, &ipn)
    if nil != err {
        return []byte{}, errors.New("Error parsing body")
    }

    log.Println(ipn)
    ipn.ReceiveDate = otpGetTimeFormat(time.Now())

    donation := Donation{Id: ipn.OrderRef}
    donation.SelectLock()
    translateStatus(&donation, ipn.Status)
    donation.UpdateUnlock()

    callback(donation)

    resp, err := json.Marshal(ipn)
    w.Header().Add("Signature", otpGenerateSignature(resp))
    return resp, err
}

func RedirectToOtpApi(dict dictionary.Dictionary, donation *Donation) (OtpReturnPublic, error) {
    url := config.Config.Donation.SimplePayURL + "/payment/v2/start"
    log.Printf("URL> %s\n", url)
    log.Printf("mer> %s\n", config.Config.Donation.Merchant)

    donation.Status = "INPROGRESS" // TODO: Create a struct for these..
    err := donation.Add()
    if nil != err {
        return OtpReturnPublic{}, err
    }

    simple_start := OtpApiStartMessage{
        Salt:           generateSalt(32),
        Merchant:       config.Config.Donation.Merchant,
        // FIXME: ::
        OrderRef:       donation.Id,
        Currency:       config.Config.Donation.Currency,
        CustomerEmail:  donation.Email,
        Language:       strings.ToUpper(dict.Meta.CountryCode),
        SdkVersion:     simpleSdkVersion,
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

    respBody, err := otpHttpRequest(url, simple_start)
    if nil != err {
        log.Println("SimplePay ERROR:")
        log.Println(err)
        return OtpReturnPublic{PaymentUrl: "/donate/" + donation.Id}, err
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

func ProgressOtpReply(r string, s string) (Donation, error) {
    payload, err := base64.StdEncoding.DecodeString(r)

    if nil != err {
        log.Println("base64 decode error.")
        log.Println(err)
        return Donation{}, err
    }

    if !signatureMatch(payload, s) {
        log.Println("Signature mismatch..")
        return Donation{}, errors.New("Payload ignature mismatch..")
    }

    simple_resp := SimpleResponse{}
    json.Unmarshal(payload, &simple_resp)

    donation := Donation{Id: simple_resp.OrderId}
    donation.SelectLock()

    // translateStatus(&donation, simple_resp.Event)
    // if donation.Successful {
    //     donation.Occurences = []time.Time{time.Now()}
    //     donation.RecurringActive = donation.Recurring
    // }
    donation.UpdateUnlock()

    return donation, nil
}

func addExistingDonationsToNewUser(user dbase.User) {
    don := Donation{}
    dons := don.List()

    for _, d := range(dons) {
        if d.Email == user.Email {
            d.UserId = user.Id.Hex()
            d.Update()
        }
    }
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

func (donation *Donation) UpdateUnlock() error {
    ddon := donation.UnMap()
    return ddon.UpdateUnlock()
}

func (donation *Donation) SelectLock() error {
    ddon := dbase.Donation{}
    oid, _ := primitive.ObjectIDFromHex(donation.Id)
    err := ddon.SelectLock(oid)
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
