const express = require('express');
var mysql = require('mysql');
var cors = require('cors');
const path = require('path');

const app = express();
app.use(express.json())
app.use(cors())

var con = mysql.createConnection({
    host: "mysql-spgroup24.alwaysdata.net",
    user: "spgroup24",
    password: "spgroup-24",
    database: "spgroup24_status"
  });

con.connect((err) => {
    if (err) throw err;
    console.log("DB Connected!");
});

app.get('/', (req, res) => {
    res.sendFile(path.join(__dirname, '/index.html'));
});

app.post('/lights/:room', (req, res) => {
    var sql = "INSERT INTO statuses (type, room, json) VALUES (1, " + req.params['room'] + ", '" + mysql_real_escape_string(JSON.stringify(req.body)) + "')";
    con.query(sql, function (err, result) {
        if (err) res.status(500).send("nok");
        res.status(201).send("ok");
        console.log("1 record inserted");
    });
});

app.get('/lights/:room', (req, res) => {
    var sql = "SELECT * FROM statuses WHERE type = 1 AND room = " + req.params['room'] + " ORDER BY datetime DESC LIMIT 1"
    con.query(sql, function (err, result) {
        if (err) res.status(500).send("nok");
        res.json(result)
    });
});

app.post('/aircon/:room', (req, res) => {
    var sql = "INSERT INTO statuses (type, room, json) VALUES (2, " + req.params['room'] + ", '" +  mysql_real_escape_string(JSON.stringify(req.body)) + "')";
    con.query(sql, function (err, result) {
        if (err) throw err;
        //res.status(500).send("nok");
        else {
            res.status(201).send("ok");
            console.log("1 record inserted");
        }
    });
});

app.get('/aircon/:room', (req, res) => {
    var sql = "SELECT * FROM statuses WHERE type = 2 AND room = " + req.params['room'] + " ORDER BY datetime DESC LIMIT 1"
    con.query(sql, function (err, result) {
        if (err) res.status(500).send("nok");
        res.json(result)
    });
});

app.listen(8100, () => console.log(`Listening on ${8100}`));

function mysql_real_escape_string (str) {
    return str.replace(/[\0\x08\x09\x1a\n\r"'\\\%]/g, function (char) {
        switch (char) {
            case "\0":
                return "\\0";
            case "\x08":
                return "\\b";
            case "\x09":
                return "\\t";
            case "\x1a":
                return "\\z";
            case "\n":
                return "\\n";
            case "\r":
                return "\\r";
            case "\"":
            case "'":
            case "\\":
            case "%":
                return "\\"+char; // prepends a backslash to backslash, percent,
                                  // and double/single quotes
            default:
                return char;
        }
    });
}
