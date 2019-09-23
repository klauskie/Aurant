const express = require('express')
var request = require('request')
const engines = require('consolidate');

const app = express()
const port = 4000

app.engine('ejs', engines.ejs);
app.set('views', './views');
app.set('views engine', 'ejs');
app.use(express.static(__dirname + '/public'));


const RESTAURANT_ID = 2;

app.get('/', function(req, res){
  res.render('home.ejs')
})

app.get('/item', (req, res) => {
  console.log('/item')
  request('http://localhost:8080/item/restaurant/' + RESTAURANT_ID, function (error, response, body) {
    if (error != null) {
      res.send("error")
    }
    if (response.statusCode != 200) {
      res.send("not 200")
    }
    var info = JSON.parse(body);
    res.render('home.ejs', {info:info})
  });
})
  

app.listen(port, function () {
  console.log('Example app listening on port 4000!');
});


// npx nodemon index.js