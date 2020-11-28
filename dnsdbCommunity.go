package main

import (
	"encoding/json"
	"fmt"
	//"net/http"
	//"io/ioutil"
	//"strings"
)

func getDNSDBCommunity (domain string) ([]byte, error){
    query := fmt.Sprintf("*.%s",domain)
    var url []string
    //Did not use rtype ANY because number of results > number of result limit(256) in Community Edition  
    url = append(url ,fmt.Sprintf("https://api.dnsdb.info/dnsdb/v2/lookup/rrset/name/%s/CNAME",query))
    url = append(url, fmt.Sprintf("https://api.dnsdb.info/dnsdb/v2/lookup/rrset/name/%s/A",query))
    url = append(url, fmt.Sprintf("https://api.dnsdb.info/dnsdb/v2/lookup/rrset/name/%s/AAAA",query))
    url = append(url, fmt.Sprintf("https://api.dnsdb.info/dnsdb/v2/lookup/rrset/name/%s/MX",query))
    url = append(url, fmt.Sprintf("https://api.dnsdb.info/dnsdb/v2/lookup/rrset/name/%s/HINFO",query))
    url = append(url, fmt.Sprintf("https://api.dnsdb.info/dnsdb/v2/lookup/rrset/name/%s/NS",query))

    body := ""
    for _, urlElement := range url{
        tempVar,err := reqDNSDB(urlElement)
        if err != nil{
            return []byte{},err
        }
        body = body+tempVar
    }

    return formatJSON(body), nil
    }

//=================================
// ====>FUNCTION FROM dnsdb.go<====
//=================================
//// Function to send Request to DNSDB
//func reqDNSDB (url string) (string,error) {
//    apikey := "" 
//    if apikey == "" {
//        return "",nil
//    }
//
//    req,err := http.NewRequest("GET", url, nil)
//    if err != nil{
//        return "", err
//    }
//
//    req.Header.Set("X-API-Key",apikey)
//    req.Header.Set("Accept", "*/*")
//
//    resp, err := http.DefaultClient.Do(req)
//    if err != nil{
//        return "", err
//    }
//    defer resp.Body.Close()
//
//    body, err := ioutil.ReadAll(resp.Body)
//    //body, err := ioutil.ReadFile("data")
//    if err != nil{
//        return "", err
//    }
//    bodySlice :=strings.Split(string(body), "\n")
//    if len(bodySlice)!=2{
//        bodystr := strings.Join(bodySlice[1:len(bodySlice)-2],"")
//        return bodystr, nil
//    }
//    bodystr :=""
//    return bodystr, nil
//}
//
//
//=================================
// ====>FUNCTION FROM dnsdb.go<====
//=================================
////Had to do it because couldn't Unmarshal it 
//func formatJSON (bodystr string) []byte{
//    occurence := strings.Count(bodystr,"]}}")
//    bodystr = strings.Replace(bodystr,"]}}","]}},\n",occurence-1)
//    body := []byte(bodystr)
//    body = append([]byte("["),body...)
//    body = append(body, []byte("]")...)
//    return body
//}

func fetchDNSDBCommunity (domain string) ([]string, error){
    body, _ := getDNSDBCommunity(domain)
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
	    if objelement.Obj.Rrtype == "CNAME" && !domainRepeatCheck[clean(objelement.Obj.Rdata[0])]{
		    domainRepeatCheck[clean(objelement.Obj.Rdata[0])]= true

		    tempvar := []byte(clean(objelement.Obj.Rdata[0])) // Removing trailing '.' from subdomains Eg: "www.tesla.com."
		    tempvar = tempvar[:len(tempvar)-2]  // Removing trailing '.' from subdomains Eg: "www.tesla.com."
		    domains = append(domains,string(tempvar))
	    }
            if domainRepeatCheck[clean(objelement.Obj.Rrname)]{
                continue
            }
            domainRepeatCheck[clean(objelement.Obj.Rrname)]= true

            tempvar := []byte(clean(objelement.Obj.Rrname)) // Removing trailing '.' from subdomains Eg: "www.tesla.com."
            tempvar = tempvar[:len(tempvar)-2]  // Removing trailing '.' from subdomains Eg: "www.tesla.com."
            domains = append(domains,string(tempvar))
        }
    return domains, nil
}
