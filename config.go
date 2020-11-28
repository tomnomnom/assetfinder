package main

import (
    "gopkg.in/yaml.v2"
    "os"
)

//Config is config.yml's data structure
type Config struct {
     Flags struct{
         BufferOverrun bool `yaml:"BufferOverrun"`
         CertSpotter bool `yaml:"CertSpotter"`
         CrtSh bool `yaml:"CrtSh"`
         DNSDBCommunity bool `yaml:"DNSDBCommunity"`
         DNSDB bool `yaml:"DNSDB"`
         Facebook bool `yaml:"Facebook"`
         FindSubDomains bool `yaml:"FindSubDomains"`
         HackerTarget bool `yaml:"HackerTarget"`
         PassiveTotal bool `yaml:"PassiveTotal"`
         SubsOnly bool `yaml:"subs-only"`
         ThreatCrowd bool `yaml:"ThreatCrowd"`
         Urlscan bool `yaml:"Urlscan"`
         VirusTotal bool `yaml:"VirusTotal"`
         Wayback bool `yaml:"Wayback"`
     } `yaml:"flags"`

    Credentials struct{
        DNSDB struct{
            APIKey string `yaml:"api-key"`
        } `yaml:"dnsdb"`

        Facebook struct{
            APPID string `yaml:"app-id"`
            AppSecret string `yaml:"app-secret"`
        } `yaml:"facebook"`

        FindSubDomains struct{
            APIToken string `yaml:"api-token"`
        } `yaml:"spyse"`

        PassiveTotal struct {
            Username string `yaml:"username"`
            Secret string `yaml:"secret"`
        } `yaml:"passivetotal"`

        VirusTotal struct {
            APIKey string `yaml:"api-key"`
        } `yaml:"virustotal"`
    } `yaml:"credentials"`
}

func (cfg *Config)ymlparser() error{
    homeDir,err := os.UserHomeDir()
    if err!= nil{
        return err
    }

    f, err := os.Open(homeDir + "/.config/assetfinder/config.yml")
    if err != nil {
        return err
    }

    defer f.Close()

    decoder := yaml.NewDecoder(f)
    err = decoder.Decode(cfg)
    if err != nil {
        return err
    }
    return nil
}
