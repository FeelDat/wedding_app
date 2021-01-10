const express = require('express');
var { Worker } = require('../models/worker');
var ObjectId = require('mongoose').Types.ObjectId;

var router = express.Router();

router.get('/', (req, res) => {
  Worker.find((err, docs) => {
    if(err){
      console.log('error: ' + JSON.stringify(err));
    }
    res.send(docs);
  })
});

router.get('/:id', (req, res) => {
  if(!ObjectId.isValid(req.params.id)){
    return res.status(400).send(`No records with given id: ${req.params.id}`)
  }
  Worker.findById(req.params.id, (err, docs) => {
    if(err){
      console.log('error: ' + JSON.stringify(err));
    }
    res.send(docs);
  })
});

router.post('/', (req, res) => {
  console.log(req.body)
  let worker = new Worker({
    name: req.body.name,
    age: req.body.age
  })

  worker.save((err, docs) => {
    if(err){
      console.log('error: ' + JSON.stringify(err));
    }
    res.send(docs);
  })
});

router.put('/:id', (req, res) => {
    if(!ObjectId.isValid(req.params.id)){
        return res.status(400).send(`No records with given id: ${req.params.id}`)
    }
    let worker = {
        name: req.body.name,
        age: req.body.age
    }
    Worker.findByIdAndUpdate(req.params.id, { $set: worker }, { new: true }, (err, docs) => {
        if(err){
            console.log('error: ' + JSON.stringify(err));
        }
        res.send(docs);
    })
})

router.delete('/:id', (req, res) => {
    if(!ObjectId.isValid(req.params.id)){
        return res.status(400).send(`No records with given id: ${req.params.id}`)
    }
    Worker.findByIdAndRemove(req.params.id, (err, docs) => {
        if(err){
            console.log('error: ' + JSON.stringify(err));
        }
        res.send(docs);
    })
})

module.exports = router;
