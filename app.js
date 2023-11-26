const express = require('express')

const app = express()
const port = 80
const version = 1

app.get('/', (req, res) => {
  res.send(version.toString())
})

app.listen(port, () => {
  console.log(`Example app listening on port ${port}`)
})
