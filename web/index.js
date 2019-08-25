const express = require('express')
var request = require('request')

const app = express()
const port = 4000

app.get('/item', (req, res) => {
  console.log('hola baka')
  res.send('Welcome to Make REST API Calls In Express!')
})
  

app.listen(port, function () {
  console.log('Example app listening on port 4000!');
});


