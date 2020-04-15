# assetfinder

Find domains and subdomains potentially related to a given domain.


## Install

If you have Go installed and configured (i.e. with `$GOPATH/bin` in your `$PATH`):

```
go get -u github.com/tomnomnom/assetfinder
```

Otherwise [download a release for your platform](https://github.com/tomnomnom/assetfinder/releases).
To make it easier to execute you can put the binary in your `$PATH`.

## Usage

```
assetfinder [--subs-only] <domain>
```

## Sources

Please feel free to issue pull requests with new sources! :)

### Implemented
* crt.sh
* certspotter
* hackertarget
* threatcrowd
* wayback machine
* dns.bufferover.run
* facebook
    * Needs `FB_APP_ID` and `FB_APP_SECRET` environment variables set (https://developers.facebook.com/)
    * You need to be careful with your app's rate limits
* virustotal
    * Needs `VT_API_KEY` environment variable set (https://developers.virustotal.com/reference)
* findsubdomains
    * Needs `SPYSE_API_TOKEN` environment variable set (the free version always gives the first response page, and you also get "25 unlimited requests") — (https://spyse.com/apidocs)

### Sources to be implemented
* http://api.passivetotal.org/api/docs/
* https://community.riskiq.com/ (?)
* https://riddler.io/
* http://www.dnsdb.org/
* https://certdb.com/api-documentation

## TODO
* Flags to control which sources are used
    * Likely to be all on by default and a flag to disable
* Read domains from stdin
