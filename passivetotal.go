package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)


func fetchPassiveTotal (domain string) ([]string,error){
    url := fmt.Sprintf("https://api.passivetotal.org/v2/enrichment/subdomains?query=%s",domain)
    username := ""
    secret := ""
    if username=="" || secret == ""{
        return []string{}, nil
    }
    wrapper := struct {
        PrimaryDomain string `json:"primaryDomain"`
        QueryValue    string `json:"queryValue"`
        Success        bool `json:"success"`
        Subdomains []string `json:"subdomains"`
    }{}

    //Web Request and Response Block
    req, err := http.NewRequest("GET", url, nil)
    if err != nil{
        return []string{}, nil
    }

    req.SetBasicAuth(username, secret)

    resp, err := http.DefaultClient.Do(req)
    if err != nil{
       return []string{}, err
   }

   defer resp.Body.Close()

   dec := json.NewDecoder(resp.Body)
   err = dec.Decode(&wrapper)
   if err != nil{
       return []string{}, err
   }

   var domains []string
   for _,subdomain := range wrapper.Subdomains{
       domains = append(domains, fmt.Sprintf("%v.%v",subdomain,wrapper.PrimaryDomain))
   }

   return domains,nil
}
