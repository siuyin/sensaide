const express = require('express');
var mysql = require('mysql');
const app = express();
app.use(express.json())

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

app.post('/lights/:room', (req, res) => {
    var sql = "INSERT INTO statuses (type, room, json) VALUES (1, " + req.params['room'] + ", '" + JSON.stringify(req.body) + "')";
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
    var sql = "INSERT INTO statuses (type, room, json) VALUES (2, " + req.params['room'] + ", '" + JSON.stringify(req.body) + "')";
    con.query(sql, function (err, result) {
        if (err) res.status(500).send("nok");
        res.status(201).send("ok");
        console.log("1 record inserted");
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
