const fs = require("fs-extra")

const benchmarkHelpers = {}

// this helper functions is 
benchmarkHelpers.saveStatsToFile = (file, txHash, name, gas, duration) => {
    fs.ensureFile(file, err => {
      if (err)
        console.log(err)
  
      data = fs.readFileSync(file)
      data = data.length === 0 ? {} : JSON.parse(data)
      data[name] = {gas, txHash, duration}
      
      fs.writeJsonSync(file, data)
    })     
}

module.exports = benchmarkHelpers