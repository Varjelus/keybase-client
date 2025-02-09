var promise = require('bluebird')
var fs = promise.promisifyAll(require('fs'))
var path = require('path')
var root = 'json'

var protocols = {}

fs.readdirAsync(root).map(load).each(analyze).then(write)

function load (file) {
  return fs.readFileAsync(path.join(root, file)).then(JSON.parse)
}

function analyze (json) {
  var prot = protocols[json.protocol]
  if (!prot) {
    prot = protocols[json.protocol] = {}
  }

  addEnums(prot, json)
}

function addEnums (prot, json) {
  json.types.filter(function (t) {
    return t.type === 'enum'
  }).forEach(function (t) {
    var en = prot[t.name] = {}

    t.symbols.forEach(function (s) {
      var parts = s.split('_')
      var val = parseInt(parts.pop(), 10)
      var name = fixCase(parts.join('_'))
      en[name] = val
    })
  })
}

function write () {
  var stream = fs.createWriteStream('js/keybase_v1.js')
  stream.once('open', function(fd) {
    Object.keys(protocols).forEach(function (p) {
      stream.write('export const ' + p + ' = ' + JSON.stringify(protocols[p], null, 2).replace(/\"/g, '\'') + '\n\n')
    })

    stream.write('export default {\n' + Object.keys(protocols).map(function(a) {return '  ' + a}).join(',\n') + '\n}')
    stream.end()
  })
}

function fixCase (s) {
  function replace (s) {
    return s[1][0].toUpperCase() + s[1].substr(1)
  }
  return s.toLowerCase().replace(/(\_\w)/g, replace)
}
