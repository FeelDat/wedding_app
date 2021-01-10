const express = require('express');
const cors = require('cors');
const bodyParser = require('body-parser');

var { mongoose } = require('./db');
var workerController = require('./controllers/workerController');

const app = express();

app.use(bodyParser.json());
app.use(cors({ origin: 'http://localhost:4200' }));

app.listen(3000, () => console.log('server started!'));

app.use('/worker', workerController );
