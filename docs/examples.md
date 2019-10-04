## Exemplos

### curl
```curl
curl https://brcep-qnlohrjtbl.now.sh/78048000/json
```

### Javascript
```javascript
var request = require('request');
var options = {
    url: 'https://brcep-qnlohrjtbl.now.sh/78048000/json',
    }
};
function callback(error, response, body) {
    if (!error && response.statusCode == 200) {
        var info = JSON.parse(body);
        console.log(info);
    }
}
request(options, callback);
```

### Python
```python
import urllib.request
import json

url = "https://brcep-qnlohrjtbl.now.sh/78048000/json"
result = urllib.request.urlopen(url)
data = result.read()
encoding = result.info().get_content_charset('utf-8')
print(json.loads(data.decode(encoding)))
```

### Golang
```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "net/url"
)

type brCep struct {
	Cep         string `json:"cep"`
	Endereco    string `json:"endereco"`
	Bairro      string `json:"bairro"`
	Complemento string `json:"complemento"`
	Cidade      string `json:"cidade"`
	Uf          string `json:"uf"`
	Latitude    string `json:"latitude"`
	Longitude   string `json:"longitude"`
	Ddd         string `json:"ddd"`
	Unidade     string `json:"unidade"`
	Ibge        string `json:"ibge"`
}

func main() {
    cep := "78048000"
    cepSeguro := url.QueryEscape(cep)

    url := fmt.Sprintf("https://brcep-dlfeappmhe.now.sh/%s/json", cepSeguro)

    req, err := http.NewRequest("GET", url, nil)

    client := &http.Client{}

    resp, err := client.Do(req)
    if err != nil {
        log.Fatal("Do: ", err)
        return
    }

    defer resp.Body.Close()
    var resultado brCep

    if err := json.NewDecoder(resp.Body).Decode(&resultado); err != nil {
        log.Println(err)
    }

    fmt.Printf("%+v\n", resultado)
}
```

### Ruby
```ruby
require "net/http"
require "uri"
require 'json'

url = "https://brcep-qnlohrjtbl.now.sh/78048000/json"
uri = URI.parse(url)

http = Net::HTTP.new(uri.host, uri.port)
http.use_ssl = true if url =~ /^https/

request = Net::HTTP::Post.new(uri.request_uri, 'Content-Type' => 'application/json')
response = http.request(request)
if response.code == "200"
    result = JSON.parse(response.body)
    puts(result)
end
```

### PHP
```php
<?php
    header ( "Content-Type: application/json;charset=utf-8" );
    $url = 'https://brcep-qnlohrjtbl.now.sh/78048000/json';
    $json = file_get_contents($url);
    echo $json;
?>
```

### Java
```java
public class BrCepGetExample {
    private static final String URL = "https://brcep-qnlohrjtbl.now.sh/78048000/json";

    public static void main(String[] args) {
        BrCepGetExample http = new BrCepGetExample();

        try {
            System.out.println("Response from brcep endpoint: " + http.sendGet());
        } catch (Exception e) {
            System.out.println("Something went wrong: " + e.getMessage());
        }
    }

    private String sendGet() throws Exception {
        URL obj = new URL(URL);

        HttpURLConnection con = (HttpURLConnection) obj.openConnection();
        con.setRequestMethod("GET");

        if (con.getResponseCode() != 200) {
            throw new IllegalStateException("Unexpected status code: " + con.getResponseCode());
        }

        BufferedReader in = new BufferedReader(new InputStreamReader(con.getInputStream()));

        String inputLine;

        StringBuilder response = new StringBuilder();
        while ((inputLine = in.readLine()) != null) {
            response.append(inputLine);
        }

        in.close();

        return response.toString();
    }
}
```

### C-Sharp
```c#
using System;
using System.Windows.Forms;
using System.Net;
using System.IO;

namespace WindowsFormsApp1
{
    public partial class Form1 : Form
    {
        public string UserAgent = @"Mozilla/5.0 (Windows; Windows NT 6.1) AppleWebKit/534.23 (KHTML, like Gecko) Chrome/11.0.686.3 Safari/534.23";

        public string HttpGet(string url)
        {
            HttpWebRequest request = (HttpWebRequest)WebRequest.Create(url);
            request.UserAgent = UserAgent;
            request.KeepAlive = false;
            request.Method = "GET";
            HttpWebResponse response = (HttpWebResponse)request.GetResponse();
            StreamReader sr = new StreamReader(response.GetResponseStream());
            return sr.ReadToEnd();
        }

        public Form1()
        {
            InitializeComponent();
        }

        private void button2_Click(object sender, EventArgs e)
        {
            richTextBox2.Text = HttpGet(textBox2.Text);
        }
    }
}
```

### Delphi
```pascal
uses idHTTP;

procedure TForm1.ButtonCEPClick(Sender: TObject);
var
    HTTP: TIdHTTP;
    IDSSLHandler : TIdSSLIOHandlerSocketOpenSSL;    
    Response: TStringStream;
    URL: String;
begin
    URL := 'https://brcep-qnlohrjtbl.now.sh/78048000/json';
    MemoReturn.Lines.Clear;
    try
        HTTP := TIdHTTP.Create;
        IDSSLHandler := TIdSSLIOHandlerSocketOpenSSL.Create;	
        HTTP.IOHandler := IDSSLHandler;	
        Response := TStringStream.Create('');
        HTTP.Get(URL, Response);
        if HTTP.ResponseCode = 200 then
           MemoReturn.Text := Utf8ToAnsi( Response.DataString );
    finally
        Response.Free;
        HTTP.Destroy;
    end;
end;
```