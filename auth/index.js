var http = require("http"),
    url = require("url"),
    port = process.env.PORT || '3000';

var auth = require('./auth.json');

function GetQuery(name, source) {
    name = name.replace(/[\[\]]/g, "\\$&");
    var regex = new RegExp("[?&]" + name + "(=([^&#]*)|&|#|$)"),
        results = regex.exec("?" + source);
    if (!results) return null;
    if (!results[2]) return '';
    return decodeURIComponent(results[2].replace(/\+/g, " "));
}

function normalizePort(val) {
    var port = parseInt(val, 10);
    if (isNaN(port)) {
        return val;
    }
    if (port >= 0) {
        return port;
    }
    return false;
}

http.createServer(function (req, res) {
    var route = req.url.split('?')[0].split('/').pop();

    if (req.method == "POST" && route == "auth") {
        var body = "";
        req.on('readable', (chunk) => {
            var data = req.read();
            body += data ? data : "";
        });
        req.on('end', () => {
            var param = {
                "app": GetQuery("app", body),
                "name": GetQuery("name", body)
            };

            if (param.app in auth.application) {
                if (auth.application[param.app].indexOf(param.name) > -1) {
                    console.log('Allow publishing : ' + param.app + "/" + param.name);
                    res.writeHead(204, { "Content-Type": "text/plain" });
                    res.write("204 No Content\n");
                } else {
                    console.log('Block publishing : ' + param.app + "/" + param.name);
                    res.writeHead(403, { "Content-Type": "text/plain" });
                    res.write("403 Forbidden\n");
                }
            } else {
                console.log('Block publishing : ' + param.app + "/" + param.name);
                res.writeHead(403, { "Content-Type": "text/plain" });
                res.write("403 Forbidden\n");
            }
            res.end();
        });
    } else {
        res.writeHead(400, { "Content-Type": "text/plain" });
        res.write("400 Bad Request\n");
        res.end();
    }
}).listen(normalizePort(port));

console.log("VRChat-RTMP Gateway Authenticator is running.");
