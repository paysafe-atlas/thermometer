function onInit() {
  const WIFI_NAME = "PaysafeGuest";
  const WIFI_OPTIONS = { password: "Courag3ous" };
  const MEASURE_INTERVAL = 5000;

  const wifi = require("Wifi");
  var http = require("http");
  const sensor = require("DS18B20").connect(new OneWire(NodeMCU.D4));

  wifi.connect(WIFI_NAME, WIFI_OPTIONS, function (err) {
    if (err) {
      console.log("Connection error: " + err);
      return;
    }
    console.log("Connected!");

    startMeasure(function (temp) {
      console.log("Temp is " + temp + "°C");
      sendTemp(temp);
    }, MEASURE_INTERVAL);
  });

  function startMeasure(callback, interval) {
    setInterval(function () {
      sensor.getTemp(callback);
    }, interval);
  }

  function sendTemp(temp) {
    var content = "{temperature:"+temp+"}";
    var options = {
      host: '10.130.11.0:8080', // host ip (host name works as well)
      port: 8080,            // (optional) port, defaults to 80
      path: '/temperature/log',           // path sent to server
      method: 'POST',
      headers: { "Content-type" : "application/json",
               "Content-length": content.length} 
    };
    
    var req = require("http").request(options, function(res) {
      console.log('res',res);
      res.on('data', function(data) {
        console.log("HTTP> "+data);
      });
      res.on('close', function(data) {
        console.log("Connection closed");
      });
    });
    
    req.on('error',function(err){
      console.log(err);
    });
    
    req.end(content);  
    console.log("Request sent");  
  }

  function sendTemperature(temp) {
    var tempDto = {
      temperature: temp
    };
    const payload = JSON.stringify(temp);
    HTTP_OPTIONS.headers = {"Content-type": "application/json",
                            "Content-Length": payload.length};
    var req = http.request(HTTP_OPTIONS, function(res) {
      console.log('STATUS: ' + res.statusCode);
      console.log('HEADERS: ' + JSON.stringify(res.headers));
      res.setEncoding('utf8');
      res.on('data', function (chunk) {
        console.log('BODY: ' + chunk);
      });
    }).end(payload);
  }
}

save();  