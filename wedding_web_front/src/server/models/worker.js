const mongoose = require('mongoose');

var Worker = mongoose.model('Worker', {
  name : { type: String },
  age  : { type: Number },
})

module.exports = {Worker};
