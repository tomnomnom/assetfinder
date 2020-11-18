package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"io/ioutil"
	"strings"
)

func getDNSDB (domain string) ([]byte, error){
    query := fmt.Sprintf("*.%s",domain)
    url := fmt.Sprintf("https://api.dnsdb.info/dnsdb/v2/lookup/rrset/name/%s/ANY",query)
    body,err := reqDNSDB(url)
    if err != nil{
        return []byte{},nil
    }
    return formatJSON(body), nil
    }



// Function to send Request to DNSDB
func reqDNSDB (url string) (string,error) {
    apikey := cfg.Credentials.DNSDB.APIKey
    // Failing silently if no API Key
    if apikey == "" {
        return "",nil
    }

    req,err := http.NewRequest("GET", url, nil)
    if err != nil{
        return "", err
    }

    req.Header.Set("X-API-Key",apikey)
    req.Header.Set("Accept", "*/*")

    resp, err := http.DefaultClient.Do(req)
    if err != nil{
        return "", err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil{
        return "", err
    }
    bodySlice :=strings.Split(string(body), "\n")
    if len(bodySlice)!=2{
        bodystr := strings.Join(bodySlice[1:len(bodySlice)-2],"")
        return bodystr, nil
    }
    bodystr :=""
    return bodystr, nil
}

//Had to do it because couldn't Unmarshal it 
func formatJSON (bodystr string) []byte{
    occurence := strings.Count(bodystr,"]}}")
    bodystr = strings.Replace(bodystr,"]}}","]}},\n",occurence-1)
    body := []byte(bodystr)
    body = append([]byte("["),body...)
    body = append(body, []byte("]")...)
    return body
}

func fetchDNSDB (domain string) ([]string, error){
    body, _ := getDNSDB(domain)
    wrapper := []struct{
            Obj struct{
                Count int `json:"count"`
                TimeFirst int `json:"time_first"`
                TimeLast int `json:"time_last"`
                Rrname string `json:"rrname"`
                Rrtype string `json:"rrtype"`
                Bailiwick string `json:"bailiwick"`
                Rdata []string `json:"rdata"`
            }`json:"obj"`
    }{}
    domainRepeatCheck := make(map[string]bool)
    var domains []string
    err := json.Unmarshal(body, &wrapper)
    if err != nil {
        return []string{}, err
    }
    for _, objelement := range wrapper{
            if domainRepeatCheck[objelement.Obj.Rrname]{
                continue
            }
            domainRepeatCheck[objelement.Obj.Rrname]= true

            tempvar := []byte(objelement.Obj.Rrname) // Removing trailing '.' from subdomains Eg: "www.tesla.com."
            tempvar = tempvar[:len(tempvar)-2]  // Removing trailing '.' from subdomains Eg: "www.tesla.com."
            domains = append(domains,string(tempvar))        }
    return domains, nil
}
