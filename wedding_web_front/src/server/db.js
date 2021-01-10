const mongoose = require('mongoose');

mongoose.connect("mongodb://localhost:27017/workerdb", (err) => {
    if(err) return console.log(err);
        console.log("mongo connected");
});

module.exports = mongoose;
